package app

import (
	"context"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/comp1x/final-task/statistics/pkg/config"
	"github.com/comp1x/final-task/statistics/pkg/repositories/statisticsrepository"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/statistics"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg config.Config) error {
	s := grpc.NewServer()
	mux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())

	go runGRPCServer(cfg, s)

	go runHTTPServer(ctx, cfg, mux)

	gracefulShutDown(s, cancel)

	return nil
}

func runGRPCServer(cfg config.Config, s *grpc.Server) {
	if err := env.Parse(&cfg.DB); err != nil {
		log.Fatalf("failed to retrieve env variables, %v", err)
	}
	if err := env.Parse(&cfg.Statistics); err != nil {
		log.Fatalf("failed to retrieve env variables, %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.PgHost, cfg.PgUser, cfg.PgPwd, cfg.PgDBName, cfg.PgPort,
	)

	StatisticsServiceServer, err := statisticsrepository.New(dsn)
	if err != nil {
		log.Fatalf("ошибка при создании MenuService: %v", err)
	}

	statistics.RegisterStatisticsServiceServer(s, StatisticsServiceServer)

	l, err := net.Listen("tcp", cfg.Statistics.GRPCAddr)
	if err != nil {
		log.Fatalf("failed to listen tcp %s, %v", cfg.Statistics.GRPCAddr, err)
	}

	log.Printf("starting listening grpc server at %s", cfg.Statistics.GRPCAddr)
	if err := s.Serve(l); err != nil {
		log.Fatalf("error services grpc server %v", err)
	}
}

func runHTTPServer(
	ctx context.Context, cfg config.Config, mux *runtime.ServeMux,
) {
	err := statistics.RegisterStatisticsServiceHandlerFromEndpoint(
		ctx,
		mux,
		"0.0.0.0"+cfg.Statistics.GRPCAddr,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("starting listening http server at %s", cfg.Statistics.HTTPAddr)
	if err := http.ListenAndServe(cfg.Statistics.HTTPAddr, mux); err != nil {
		log.Fatalf("error service http server %v", err)
	}
}

func gracefulShutDown(s *grpc.Server, cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)

	sig := <-ch
	errorMessage := fmt.Sprintf("%s %v - %s", "Received shutdown signal:", sig, "Graceful shutdown done")
	log.Println(errorMessage)
	s.GracefulStop()
	cancel()
}
