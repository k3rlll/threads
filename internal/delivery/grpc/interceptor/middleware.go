package interceptor

import (
	"context"
	"log/slog"
	ctxUtil "main/pkg/utils/context"
	"runtime/debug"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var publicMethods = map[string]struct{}{
	"/auth.v1.AuthService/Register": {},
	"/auth.v1.AuthService/Login":    {},
}

type JWTManager interface {
	VerifyAccessToken(tokenString string) (userID uuid.UUID, err error)
}

// AuthInterceptor is a gRPC middleware that intercepts incoming requests to perform authentication.
func AuthInterceptor(jwtManager JWTManager) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if _, ok := publicMethods[info.FullMethod]; ok {
			// Public method, proceed without authentication
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		values := md["authorization"]
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "missing authorization token")
		}

		accessToken := strings.TrimPrefix(values[0], "Bearer ")

		userID, err := jwtManager.VerifyAccessToken(accessToken)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		newCtx := ctxUtil.NewContext(ctx, userID.String())

		return handler(newCtx, req)
	}
}

// LoggingInterceptor is a gRPC middleware that intercepts errors returned by handlers and logs them appropriately.
func LoggingInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)

		if err == nil {
			logger.Info("gRPC Request",
				"method", info.FullMethod,
				"request", req,
				"response", resp,
			)
			return resp, nil
		}

		st, ok := status.FromError(err)

		if ok {

			logger.Warn("gRPC Client Error",
				"method", info.FullMethod,
				"code", st.Code(),
				"msg", st.Message(),
			)
			return resp, err
		}

		logger.Error("gRPC SYSTEM ERROR",
			"method", info.FullMethod,
			"err", err,
		)

		return nil, status.Error(codes.Internal, "internal server error")
	}
}

// RecoveryInterceptor is a gRPC middleware that recovers from panics in handlers and logs the panic details.
func RecoveryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {

				stackTrace := string(debug.Stack())

				logger.Error("PANIC RECOVERED",
					"method", info.FullMethod,
					"panic", r,
					"stack", stackTrace,
				)

				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}
