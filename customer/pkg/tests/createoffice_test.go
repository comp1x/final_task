package tests

import (
	"context"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/comp1x/final-task/customer/pkg/config"
	"github.com/comp1x/final-task/customer/pkg/repositories/officerepository"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/customer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
	"time"
)

func TestCreateOffice(t *testing.T) {
	// Server
	cfg := config.Config{}

	if err := env.Parse(&cfg); err != nil {
		t.Fatalf("env.Parse(%v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DB.PgHost, cfg.DB.PgUser, cfg.DB.PgPwd, cfg.DB.PgDBName, cfg.DB.PgPort,
	)

	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc, err := officerepository.New(dsn)

	if err != nil {
		t.Fatalf("officerepository.New %v", err)
	}

	customer.RegisterOfficeServiceServer(srv, svc)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("srv.Serve %v", err)
		}
	}()
	// Test
	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	t.Cleanup(func() {
		conn.Close()
	})
	if err != nil {
		t.Fatalf("grpc.DialContext %v", err)
	}
	client := customer.NewOfficeServiceClient(conn)

	_, err = client.CreateOffice(ctx, &customer.CreateOfficeRequest{
		Name:    "mediasoft",
		Address: "gagarina dom 3",
	})
	if err != nil {
		t.Fatalf("client.CreateOffice %v", err)
	}
}
