package app

import (
	"context"
	"fmt"
	"github.com/comp1x/final-task/customer/pkg/config"
	officerepository "github.com/comp1x/final-task/customer/pkg/repositories/officerepository"
	orderrepository "github.com/comp1x/final-task/customer/pkg/repositories/orderrepository"
	userrepository "github.com/comp1x/final-task/customer/pkg/repositories/userrepository"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	customer "gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/customer"
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
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DB.PgHost, cfg.DB.PgUser, cfg.DB.PgPwd, cfg.DB.PgDBName, cfg.DB.PgPort,
	)
	OfficeServiceServer, err := officerepository.New(dsn)
	if err != nil {
		log.Fatalf("ошибка при создании OfficeService: %v", err)
	}
	UserServiceServer, err := userrepository.New(dsn)
	OrderServiceServer, err := orderrepository.New(dsn)

	if err != nil {
		log.Fatalf("ошибка при создании UserService: %v", err)
	}

	customer.RegisterOfficeServiceServer(s, OfficeServiceServer)
	customer.RegisterOrderServiceServer(s, OrderServiceServer)
	customer.RegisterUserServiceServer(s, UserServiceServer)

	l, err := net.Listen("tcp", cfg.Customer.GRPCAddr)
	if err != nil {
		log.Fatalf("failed to listen tcp %s, %v", cfg.Customer.GRPCAddr, err)
	}

	log.Printf("starting listening grpc server at %s", cfg.Customer.GRPCAddr)
	if err := s.Serve(l); err != nil {
		log.Fatalf("error services grpc server %v", err)
	}
}

func runHTTPServer(
	ctx context.Context, cfg config.Config, mux *runtime.ServeMux,
) {
	err := customer.RegisterOfficeServiceHandlerFromEndpoint(
		ctx,
		mux,
		"0.0.0.0"+cfg.Customer.GRPCAddr,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)

	if err != nil {
		log.Fatal(err)
	}

	err = customer.RegisterOrderServiceHandlerFromEndpoint(
		ctx,
		mux,
		"0.0.0.0"+cfg.Customer.GRPCAddr,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)

	if err != nil {
		log.Fatal(err)
	}

	err = customer.RegisterUserServiceHandlerFromEndpoint(
		ctx,
		mux,
		"0.0.0.0"+cfg.Customer.GRPCAddr,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("starting listening http server at %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, mux); err != nil {
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
