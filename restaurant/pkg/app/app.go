package app

import (
	"context"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/comp1x/final-task/restaurant/pkg/config"
	"github.com/comp1x/final-task/restaurant/pkg/repositories/menurepository"
	"github.com/comp1x/final-task/restaurant/pkg/repositories/orderrepository"
	"github.com/comp1x/final-task/restaurant/pkg/repositories/productrepository"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/restaurant"
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
	if err := env.Parse(&cfg.Restaurant); err != nil {
		log.Fatalf("failed to retrieve env variables, %v", err)
	}
	if err := env.Parse(&cfg.Customer); err != nil {
		log.Fatalf("failed to retrieve env variables, %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DB.PgHost, cfg.DB.PgUser, cfg.DB.PgPwd, cfg.DB.PgDBName, cfg.DB.PgPort,
	)

	ProductServiceServer, err := productrepository.New(dsn)
	if err != nil {
		log.Fatalf("ошибка при создании OfficeService: %v", err)
	}
	MenuServiceServer, err := menurepository.New(dsn)
	if err != nil {
		log.Fatalf("ошибка при создании MenuService: %v", err)
	}
	OrderServiceServer, err := orderrepository.New(dsn)
	if err != nil {
		log.Fatalf("ошибка при создании MenuService: %v", err)
	}

	if err != nil {
		log.Fatalf(
			"ошибка при создании UserService: %v", err)
	}

	restaurant.RegisterProductServiceServer(s, ProductServiceServer)
	restaurant.RegisterMenuServiceServer(s, MenuServiceServer)
	restaurant.RegisterOrderServiceServer(s, OrderServiceServer)

	l, err := net.Listen("tcp", cfg.Restaurant.GRPCAddr)
	if err != nil {
		log.Fatalf("failed to listen tcp %s, %v", cfg.Restaurant.GRPCAddr, err)
	}

	log.Printf("starting listening grpc server at %s", cfg.Restaurant.GRPCAddr)
	if err := s.Serve(l); err != nil {
		log.Fatalf("error services grpc server %v", err)
	}
}

func runHTTPServer(
	ctx context.Context, cfg config.Config, mux *runtime.ServeMux,
) {
	err := restaurant.RegisterMenuServiceHandlerFromEndpoint(
		ctx,
		mux,
		"0.0.0.0"+cfg.Restaurant.GRPCAddr,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)

	if err != nil {
		log.Fatal(err)
	}

	err = restaurant.RegisterOrderServiceHandlerFromEndpoint(
		ctx,
		mux,
		"0.0.0.0"+cfg.Restaurant.GRPCAddr,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)

	if err != nil {
		log.Fatal(err)
	}

	err = restaurant.RegisterProductServiceHandlerFromEndpoint(
		ctx,
		mux,
		"0.0.0.0"+cfg.Restaurant.GRPCAddr,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("starting listening http server at %s", cfg.Restaurant.HTTPAddr)
	if err := http.ListenAndServe(cfg.Restaurant.HTTPAddr, mux); err != nil {
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
