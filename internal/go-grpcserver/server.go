package grpcserver

import (
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"

	"google.golang.org/grpc"
)

const (
	defaultGracefulShutdownTimeout = 10 * time.Second
)

type Server struct {
	grpcSrv *grpc.Server

	GracefulShutdownTimeout time.Duration
}

func New() *Server {
	return &Server{
		GracefulShutdownTimeout: defaultGracefulShutdownTimeout,
	}
}

func (s *Server) Create() (*grpc.Server, error) {

	chain := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
	}

	s.grpcSrv = grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(chain...),
		// grpc.StatsHandler()
		// TODO: this is supported by opencensus but unsupported by opentelemetry-go:
		// https://github.com/open-telemetry/opentelemetry-go/issues/764
	)

	return s.grpcSrv, nil
}

func (s *Server) Shutdown() bool {
	c := make(chan struct{})

	go func() {
		defer close(c)

		// Block until all pending RPCs are finished
		s.grpcSrv.GracefulStop()
	}()

	select {
	case <-time.After(s.GracefulShutdownTimeout):
		// Timeout
		s.grpcSrv.Stop()
		<-c
		return false

	case <-c:
		// Shutdown completed within the timeout
		return true
	}
}
