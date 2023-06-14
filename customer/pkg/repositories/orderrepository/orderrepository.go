package orderrepositoryimport

import (
	"context"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/comp1x/final-task/customer/pkg/config"
	"github.com/comp1x/final-task/customer/pkg/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	customer "gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/customer"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/restaurant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type OrderService struct {
	customer.UnimplementedOrderServiceServer
	db     *gorm.DB
	logger logrus.FieldLogger
}

func New(dbURL string, logger logrus.FieldLogger) (*OrderService, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		logger.Error("New (OfficeService): ", err, time.Now().UTC())
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	return &OrderService{
		db:     db,
		logger: logger,
	}, nil
}

func (s *OrderService) CreateOrder(
	ctx context.Context, request *customer.CreateOrderRequest,
) (*customer.CreateOrderResponse, error) {
	if err := request.ValidateAll(); err != nil {
		s.logger.Error("CreateOrder: ", err, time.Now().UTC())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var Orders []*models.Order
	for _, OrderItem := range request.Salads {
		ProductUuid, err := uuid.Parse(OrderItem.ProductUuid)
		if err != nil {
			s.logger.Error("CreateOrder: ", err, time.Now().UTC())
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		UserUuid, err := uuid.Parse(request.UserUuid)
		if err != nil {
			s.logger.Error("CreateOrder: ", err, time.Now().UTC())
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		order := &models.Order{
			ProductUuid: ProductUuid,
			UserUuid:    UserUuid,
			Count:       int64(OrderItem.Count),
		}
		Orders = append(Orders, order)
	}

	for _, OrderItem := range request.Meats {
		ProductUuid, err := uuid.Parse(OrderItem.ProductUuid)
		if err != nil {
			s.logger.Error("CreateOrder: ", err, time.Now().UTC())
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		UserUuid, err := uuid.Parse(request.UserUuid)
		if err != nil {
			s.logger.Error("CreateOrder: ", err, time.Now().UTC())
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		order := &models.Order{
			ProductUuid: ProductUuid,
			UserUuid:    UserUuid,
			Count:       int64(OrderItem.Count),
		}
		Orders = append(Orders, order)
	}

	for _, OrderItem := range request.Soups {
		ProductUuid, err := uuid.Parse(OrderItem.ProductUuid)
		if err != nil {
			s.logger.Error("CreateOrder: ", err, time.Now().UTC())
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		UserUuid, err := uuid.Parse(request.UserUuid)
		if err != nil {
			s.logger.Error("CreateOrder: ", err, time.Now().UTC())
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		order := &models.Order{
			ProductUuid: ProductUuid,
			UserUuid:    UserUuid,
			Count:       int64(OrderItem.Count),
		}
		Orders = append(Orders, order)
	}

	for _, OrderItem := range request.Drinks {
		ProductUuid, err := uuid.Parse(OrderItem.ProductUuid)
		if err != nil {
			s.logger.Error("CreateOrder: ", err, time.Now().UTC())
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		UserUuid, err := uuid.Parse(request.UserUuid)
		if err != nil {
			s.logger.Error("CreateOrder: ", err, time.Now().UTC())
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		order := &models.Order{
			ProductUuid: ProductUuid,
			UserUuid:    UserUuid,
			Count:       int64(OrderItem.Count),
		}
		Orders = append(Orders, order)
	}

	for _, OrderItem := range request.Desserts {
		ProductUuid, err := uuid.Parse(OrderItem.ProductUuid)
		if err != nil {
			s.logger.Error("CreateOrder: ", err, time.Now().UTC())
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		UserUuid, err := uuid.Parse(request.UserUuid)
		if err != nil {
			s.logger.Error("CreateOrder: ", err, time.Now().UTC())
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		order := &models.Order{
			ProductUuid: ProductUuid,
			UserUuid:    UserUuid,
			Count:       int64(OrderItem.Count),
		}
		Orders = append(Orders, order)
	}
	for _, order := range Orders {
		if err := s.db.WithContext(ctx).Create(order).Error; err != nil {
			s.logger.Error("CreateOrder: ", err, time.Now().UTC())
			return nil, status.Error(codes.Unavailable, err.Error())
		}
	}

	s.logger.Println("CreateOffice: ", Orders)

	return &customer.CreateOrderResponse{}, nil
}

func (s *OrderService) GetActualMenu(
	ctx context.Context, request *customer.GetActualMenuRequest,
) (*customer.GetActualMenuResponse, error) {
	if err := request.ValidateAll(); err != nil {
		s.logger.Error("GetActualMenu: ", err, time.Now().UTC())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	cfg := config.Config{}

	if err := env.Parse(&cfg.Restaurant); err != nil {
		s.logger.Error("GetActualMenu: ", err, time.Now().UTC())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var conn *grpc.ClientConn

	conn, err := grpc.Dial(cfg.Restaurant.GRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.logger.Error("GetActualMenu: ", err, time.Now().UTC())
		return nil, err
	}

	yyyy, mm, dd := time.Now().Date()
	tomorrow := time.Date(yyyy, mm, dd+1, 10, 0, 0, 0, time.Now().Location())

	client := restaurant.NewMenuServiceClient(conn)

	requestGetMenu := &restaurant.GetMenuRequest{
		OnDate: timestamppb.New(tomorrow),
	}

	result, err := client.GetMenu(ctx, requestGetMenu)

	conn.Close()

	if err != nil {
		s.logger.Error("GetActualMenu: ", err, time.Now().UTC())
		return nil, fmt.Errorf("No menu on tomorrow or restaurant service is down.")
	}

	apiSalads := make([]*customer.Product, 0, len(result.Menu.Salads))
	for _, product := range result.Menu.Salads {
		apiSalad := &customer.Product{
			Uuid:        product.Uuid,
			Name:        product.Name,
			Description: product.Description,
			Type:        customer.CustomerProductType(product.Type),
			Weight:      product.Weight,
			Price:       product.Price,
			CreatedAt:   product.CreatedAt,
		}
		apiSalads = append(apiSalads, apiSalad)
	}

	apiGarnishes := make([]*customer.Product, 0, len(result.Menu.Garnishes))
	for _, product := range result.Menu.Garnishes {
		apiGarnish := &customer.Product{
			Uuid:        product.Uuid,
			Name:        product.Name,
			Description: product.Description,
			Type:        customer.CustomerProductType(product.Type),
			Weight:      product.Weight,
			Price:       product.Price,
			CreatedAt:   product.CreatedAt,
		}
		apiGarnishes = append(apiGarnishes, apiGarnish)
	}

	apiMeats := make([]*customer.Product, 0, len(result.Menu.Meats))
	for _, product := range result.Menu.Meats {
		apiMeat := &customer.Product{
			Uuid:        product.Uuid,
			Name:        product.Name,
			Description: product.Description,
			Type:        customer.CustomerProductType(product.Type),
			Weight:      product.Weight,
			Price:       product.Price,
			CreatedAt:   product.CreatedAt,
		}
		apiMeats = append(apiSalads, apiMeat)
	}

	apiSoups := make([]*customer.Product, 0, len(result.Menu.Soups))
	for _, product := range result.Menu.Soups {
		apiSoup := &customer.Product{
			Uuid:        product.Uuid,
			Name:        product.Name,
			Description: product.Description,
			Type:        customer.CustomerProductType(product.Type),
			Weight:      product.Weight,
			Price:       product.Price,
			CreatedAt:   product.CreatedAt,
		}
		apiSoups = append(apiSoups, apiSoup)
	}

	apiDrinks := make([]*customer.Product, 0, len(result.Menu.Drinks))
	for _, product := range result.Menu.Drinks {
		apiDrink := &customer.Product{
			Uuid:        product.Uuid,
			Name:        product.Name,
			Description: product.Description,
			Type:        customer.CustomerProductType(product.Type),
			Weight:      product.Weight,
			Price:       product.Price,
			CreatedAt:   product.CreatedAt,
		}
		apiDrinks = append(apiDrinks, apiDrink)
	}

	apiDesserts := make([]*customer.Product, 0, len(result.Menu.Desserts))
	for _, product := range result.Menu.Desserts {
		apiDessert := &customer.Product{
			Uuid:        product.Uuid,
			Name:        product.Name,
			Description: product.Description,
			Type:        customer.CustomerProductType(product.Type),
			Weight:      product.Weight,
			Price:       product.Price,
			CreatedAt:   product.CreatedAt,
		}
		apiDesserts = append(apiDesserts, apiDessert)
	}

	s.logger.Error("GetActualMenu: ",
		apiSalads,
		apiGarnishes,
		apiMeats,
		apiSoups,
		apiDrinks,
		apiDesserts,
		time.Now().UTC(),
	)

	return &customer.GetActualMenuResponse{
		Salads:    apiSalads,
		Garnishes: apiGarnishes,
		Meats:     apiMeats,
		Soups:     apiSoups,
		Drinks:    apiDrinks,
		Desserts:  apiDesserts,
	}, nil
}
