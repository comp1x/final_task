package statisticsrepository

import (
	"context"
	"fmt"
	"github.com/comp1x/final-task/restaurant/pkg/models"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/statistics"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type StatisticsService struct {
	statistics.UnimplementedStatisticsServiceServer
	db *gorm.DB
}

func New(dbURL string) (*StatisticsService, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к базе данных: %w", err)
	}

	return &StatisticsService{
		db: db,
	}, nil
}

func (s *StatisticsService) GetAmountOfProfit(
	ctx context.Context, request *statistics.GetAmountOfProfitRequest,
) (*statistics.GetAmountOfProfitResponse, error) {
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var orders []models.Order
	if err := s.db.Table("orders").Where("created_at > ? AND created_at < ?", request.StartDate.AsTime(), request.EndDate.AsTime()).Find(&orders); err != nil {
		log.Printf("ошибка при получении заказов из базы данных: %v", err)
		return nil, fmt.Errorf("ошибка при получении списка заказов")
	}
	fmt.Println(orders)
	return nil, nil
}

func (s *StatisticsService) TopProducts(
	ctx context.Context, request *statistics.TopProductsRequest,
) (*statistics.TopProductsResponse, error) {
	return nil, nil
}
