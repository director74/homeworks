package internalgrpc

import (
	"fmt"
	"net"

	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/app"
	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/cfg"
	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"google.golang.org/grpc"
)

type Server struct {
	srv     *grpc.Server
	host    string
	port    string
	logg    app.Logger
	storage app.Storage
}

func NewServer(logger app.Logger, storage app.Storage, grpcConf cfg.GRPCServerConf) *Server {
	return &Server{
		port:    grpcConf.Port,
		host:    grpcConf.Host,
		logg:    logger,
		storage: storage,
	}
}

func (s *Server) Start() error {
	lsn, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		return err
	}
	s.srv = grpc.NewServer(s.withServerUnaryInterceptor())
	pb.RegisterCalendarServer(s.srv, NewService(s.storage))
	s.logg.Info(fmt.Sprintf("starting grpc server on %s", lsn.Addr().String()))
	return s.srv.Serve(lsn)
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
	s.logg.Info("grpc server stopped")
}

func (s *Server) withServerUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(UnaryServerInterceptor(s.logg))
}
