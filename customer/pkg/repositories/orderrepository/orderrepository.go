package orderrepositoryimport

import (
	"context"
	"github.com/comp1x/final-task/customer/pkg/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	customer "gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/customer"
	"google.golang.org/grpc/codes"
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
		s.logger.Error("CreateOrder: ", err, time.Now().UTC().UTC())
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

	Year, Month, Day := time.Now().UTC().Date()

	var menu models.Menu
	if err := s.db.WithContext(ctx).Where("year = ? AND month = ? AND day = ?", Year, int(Month), Day+1).First(&menu).Error; err != nil {
		s.logger.Errorf("GetActualMenu: ", err, time.Now().UTC())
		return nil, status.Error(codes.Unknown, err.Error())
	}

	var products []models.Product
	if err := s.db.WithContext(ctx).Where("id IN ?", []string(menu.ProductsUuids)).Find(&products).Error; err != nil {
		s.logger.Error("GetActualMenu: ", err, time.Now().UTC())
		return nil, status.Error(codes.Unknown, err.Error())
	}
	apiUnspecified := make([]*customer.Product, 0)
	apiSalads := make([]*customer.Product, 0)
	apiGarnishes := make([]*customer.Product, 0)
	apiMeats := make([]*customer.Product, 0)
	apiSoups := make([]*customer.Product, 0)
	apiDrinks := make([]*customer.Product, 0)
	apiDesserts := make([]*customer.Product, 0)
	for _, product := range products {
		switch product.Type {
		case 0:
			apiProduct := &customer.Product{
				Uuid:        product.ID.String(),
				Name:        product.Name,
				Description: product.Description,
				Type:        product.Type,
				Weight:      product.Weight,
				Price:       product.Price,
				CreatedAt:   timestamppb.New(product.CreatedAt),
			}
			apiUnspecified = append(apiUnspecified, apiProduct)
		case 1:
			apiProduct := &customer.Product{
				Uuid:        product.ID.String(),
				Name:        product.Name,
				Description: product.Description,
				Type:        product.Type,
				Weight:      product.Weight,
				Price:       product.Price,
				CreatedAt:   timestamppb.New(product.CreatedAt),
			}
			apiSalads = append(apiSalads, apiProduct)
		case 2:
			apiProduct := &customer.Product{
				Uuid:        product.ID.String(),
				Name:        product.Name,
				Description: product.Description,
				Type:        product.Type,
				Weight:      product.Weight,
				Price:       product.Price,
				CreatedAt:   timestamppb.New(product.CreatedAt),
			}
			apiGarnishes = append(apiGarnishes, apiProduct)
		case 3:
			apiProduct := &customer.Product{
				Uuid:        product.ID.String(),
				Name:        product.Name,
				Description: product.Description,
				Type:        product.Type,
				Weight:      product.Weight,
				Price:       product.Price,
				CreatedAt:   timestamppb.New(product.CreatedAt),
			}
			apiMeats = append(apiMeats, apiProduct)
		case 4:
			apiProduct := &customer.Product{
				Uuid:        product.ID.String(),
				Name:        product.Name,
				Description: product.Description,
				Type:        product.Type,
				Weight:      product.Weight,
				Price:       product.Price,
				CreatedAt:   timestamppb.New(product.CreatedAt),
			}
			apiSoups = append(apiSoups, apiProduct)
		case 5:
			apiProduct := &customer.Product{
				Uuid:        product.ID.String(),
				Name:        product.Name,
				Description: product.Description,
				Type:        product.Type,
				Weight:      product.Weight,
				Price:       product.Price,
				CreatedAt:   timestamppb.New(product.CreatedAt),
			}
			apiDrinks = append(apiDrinks, apiProduct)
		case 6:
			apiProduct := &customer.Product{
				Uuid:        product.ID.String(),
				Name:        product.Name,
				Description: product.Description,
				Type:        product.Type,
				Weight:      product.Weight,
				Price:       product.Price,
				CreatedAt:   timestamppb.New(product.CreatedAt),
			}
			apiDesserts = append(apiDesserts, apiProduct)
		}
	}

	s.logger.Println("GetActualMenu: ",
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
