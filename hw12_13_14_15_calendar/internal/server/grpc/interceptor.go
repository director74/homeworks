package internalgrpc

import (
	"fmt"
	"time"

	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/app"
	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/storage"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor(logg app.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// Calls the handler
		h, err := handler(ctx, req)

		rqDuration := time.Since(start).Seconds()

		ip, _ := getClientIP(ctx)
		useragent := getUserAgent(ctx)

		logg.Infof("%s [%s] %s %s %f %s",
			ip,
			start.Format(storage.LayoutLog),
			info.FullMethod,
			status.Code(err).String(),
			rqDuration,
			useragent,
		)

		return h, err
	}
}

func getClientIP(ctx context.Context) (string, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("couldn't parse client IP address")
	}
	return p.Addr.String(), nil
}

func getUserAgent(ctx context.Context) string {
	result := ""
	metadata, ok := metadata.FromIncomingContext(ctx)
	if ok {
		useragent := metadata.Get("user-agent")
		if len(useragent) > 0 {
			result = useragent[0]
		}
	}

	return result
}
