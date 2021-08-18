package controller

import (
	"storage/app/interface/controller/storage"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	hawk_grpc "github.com/tyr-tech-team/hawk/middleware/grpc"
	"google.golang.org/grpc"
)

//NewGrpcServer -
func NewGrpcServer(log *logrus.Entry, s storage.StorageServiceServer) *grpc.Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(
					grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor),
				),
				grpc_logrus.UnaryServerInterceptor(log),
				grpc_auth.UnaryServerInterceptor(hawk_grpc.TraceID),
				grpc_auth.UnaryServerInterceptor(hawk_grpc.GetOperator),
				grpc_recovery.UnaryServerInterceptor(),
				//interceptor.UnaryServerInterceptor(),
			),
		),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpc_ctxtags.StreamServerInterceptor(
					grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor),
				),
				grpc_logrus.StreamServerInterceptor(log),
				grpc_auth.StreamServerInterceptor(hawk_grpc.TraceID),
				grpc_auth.StreamServerInterceptor(hawk_grpc.GetOperator),
				grpc_recovery.StreamServerInterceptor(),
				//interceptor.StreamServerInterceptor(),
			),
		),
	)

	storage.RegisterStorageServiceServer(server, s)

	return server
}
