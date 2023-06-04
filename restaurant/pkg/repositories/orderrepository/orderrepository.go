package orderrepository

import (
	"context"
	"fmt"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/restaurant"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type OrderService struct {
	restaurant.UnimplementedOrderServiceServer

	db *gorm.DB
}

func New(dbURL string) (*OrderService, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к базе данных: %w", err)
	}

	return &OrderService{
		db: db,
	}, nil
}

func (s *OrderService) GetUpToDateOrderList(
	ctx context.Context, request *restaurant.GetUpToDateOrderListRequest,
) (*restaurant.GetUpToDateOrderListResponse, error) {
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &restaurant.GetUpToDateOrderListResponse{}, nil
}
