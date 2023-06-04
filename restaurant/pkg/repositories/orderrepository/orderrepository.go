package orderrepository

import (
	"context"
	"fmt"
	modelsCustomer "github.com/comp1x/final-task/customer/pkg/models"
	"github.com/comp1x/final-task/restaurant/pkg/models"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/restaurant"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
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

	var orders []models.Order
	if err := s.db.Find(&orders).Error; err != nil {
		log.Printf("ошибка при получении заказов из базы данных: %v", err)
		return nil, fmt.Errorf("ошибка при получении списка заказов")
	}

	apiTotalOrders := make([]*restaurant.Order, 0)
	for _, order := range orders {
		apiOrder := &restaurant.Order{
			ProductId:   order.ProductUuid.String(),
			ProductName: order.Product.Name,
			Count:       order.Count,
		}
		apiTotalOrders = append(apiTotalOrders, apiOrder)
	}

	var offices []modelsCustomer.Office
	if err := s.db.Find(&offices).Error; err != nil {
		log.Printf("ошибка при получении офисов из базы данных: %v", err)
		return nil, fmt.Errorf("ошибка при получении списка офисов")
	}

	apiOrdersByOffice := make([]*restaurant.OrdersByOffice, 0)

	for _, office := range offices {
		apiOrdersByCompany := make([]*restaurant.Order, 0)
		var orders []models.Order
		s.db.Where(&models.Order{
			User: modelsCustomer.User{
				OfficeUuid: office.ID,
			},
		}).Find(&orders)
		for _, order := range orders {
			apiOrder := &restaurant.Order{
				ProductId:   order.ProductUuid.String(),
				ProductName: order.Product.Name,
				Count:       order.Count,
			}
			apiOrdersByCompany = append(apiOrdersByCompany, apiOrder)
		}
		apiOrdersByOffice = append(apiOrdersByOffice, &restaurant.OrdersByOffice{
			OfficeUuid:    office.ID.String(),
			OfficeName:    office.Name,
			OfficeAddress: office.Address,
			Result:        apiOrdersByCompany,
		})
	}

	return &restaurant.GetUpToDateOrderListResponse{
		TotalOrders:          apiTotalOrders,
		TotalOrdersByCompany: apiOrdersByOffice,
	}, nil
}
