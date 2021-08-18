package cmd

import (
	"context"
	"log"
	"storage/app/domain/service"
	"storage/app/infra/config"
	"storage/app/infra/google"
	"storage/app/infra/logger"
	"storage/app/infra/mongo"
	"storage/app/infra/trace"
	"storage/app/interface/controller"
	"storage/app/interface/controller/storage"
	"storage/app/usecase"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tyr-tech-team/hawk/srv"
	"go.uber.org/fx"
	"google.golang.org/grpc"

	mongo_storage "storage/app/infra/mongo/storage"

	google_storage "storage/app/infra/google/storage"
)

var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&config.C.Info.RemoteHost, "consul", "c", "127.0.0.1:8500", "consul host")
	serverCmd.Flags().StringVarP(&config.C.Info.Port, "port", "p", "0", "service port")
	serverCmd.Flags().StringVarP(&config.C.Info.Mode, "mode", "", "development", "server mode")
	serverCmd.Flags().StringVarP(&config.C.Google.Credentials, "file", "f", "", "file path")
}

func start() {

	app := fx.New(
		fx.NopLogger,
		fx.Provide(
			context.Background,

			// config
			config.NewConsulClient,
			config.RemoteConfig,
			config.RegisterClient,

			// logger
			logger.NewLogger,

			// mongo
			mongo.NewDial,

			// google
			google.NewDial,

			// new mongo repository
			mongo_storage.NewRepository,

			// google repository
			google_storage.NewRepository,

			// service
			service.NewService,

			// usecase
			usecase.NewStorageUsecase,

			// grpc service
			storage.NewService,

			// grpc server
			controller.NewGrpcServer,
		),
		fx.Invoke(NewGRPCServer),
	)

	if err := app.Err(); err != nil {
		log.Fatal(err)
	}

	app.Run()
}

func NewGRPCServer(usecase usecase.StorageUsecase, lc fx.Lifecycle, server *grpc.Server, register srv.Register, log *logrus.Entry) error {
	s := srv.New(
		srv.SetName(config.C.Info.Name),
		srv.SetHost(config.C.Info.Host),
		srv.SetRegister(register),
		srv.SetGRPC(),
		srv.SetEnableTraefik(),
	)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			trace.NewTrace(config.C)
			go server.Serve(s.GetListener())
			s.Register()
			log.Info("start service on ", s.GetHost())
			return nil
		},
		OnStop: func(context.Context) error {
			s.Close()
			s.Deregister()
			return nil
		},
	})

	return nil
}
