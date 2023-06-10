package menurepository

import (
	"context"
	"github.com/comp1x/final-task/restaurant/pkg/models"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/restaurant"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MenuService struct {
	restaurant.UnimplementedMenuServiceServer

	db *gorm.DB
}

func New(dbURL string) (*MenuService, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	return &MenuService{
		db: db,
	}, nil
}

func (s *MenuService) CreateMenu(
	ctx context.Context, request *restaurant.CreateMenuRequest,
) (*restaurant.CreateMenuResponse, error) {
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var ProductsUuids []string

	ProductsUuids = append(ProductsUuids, request.Salads...)
	ProductsUuids = append(ProductsUuids, request.Garnishes...)
	ProductsUuids = append(ProductsUuids, request.Meats...)
	ProductsUuids = append(ProductsUuids, request.Soups...)
	ProductsUuids = append(ProductsUuids, request.Drinks...)
	ProductsUuids = append(ProductsUuids, request.Desserts...)

	menu := &models.Menu{
		OnDate:          request.OnDate.AsTime(),
		OpeningRecordAt: request.OpeningRecordAt.AsTime(),
		ClosingRecordAt: request.ClosingRecordAt.AsTime(),
		Year:            request.OnDate.AsTime().Year(),
		Month:           int(request.OnDate.AsTime().Month()),
		Day:             request.OnDate.AsTime().Day(),
		ProductsUuids:   ProductsUuids,
	}

	if err := s.db.WithContext(ctx).Create(menu).Error; err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}
	return &restaurant.CreateMenuResponse{}, nil
}

func (s *MenuService) GetMenu(
	ctx context.Context, request *restaurant.GetMenuRequest,
) (*restaurant.GetMenuResponse, error) {
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	Year, Month, Day := request.GetOnDate().AsTime().Date()

	var menu models.Menu
	if err := s.db.WithContext(ctx).Where("year = ? AND month = ? AND day = ?", Year, int(Month), Day).First(&menu).Error; err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	var products []models.Product
	if err := s.db.WithContext(ctx).Where("id IN ?", []string(menu.ProductsUuids)).Find(&products).Error; err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	apiUnspecified := make([]*restaurant.Product, 0)
	apiSalads := make([]*restaurant.Product, 0)
	apiGarnishes := make([]*restaurant.Product, 0)
	apiMeats := make([]*restaurant.Product, 0)
	apiSoups := make([]*restaurant.Product, 0)
	apiDrinks := make([]*restaurant.Product, 0)
	apiDesserts := make([]*restaurant.Product, 0)
	for _, product := range products {
		switch product.Type {
		case 0:
			apiProduct := &restaurant.Product{
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
			apiProduct := &restaurant.Product{
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
			apiProduct := &restaurant.Product{
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
			apiProduct := &restaurant.Product{
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
			apiProduct := &restaurant.Product{
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
			apiProduct := &restaurant.Product{
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
			apiProduct := &restaurant.Product{
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

	apiMenu := &restaurant.Menu{
		Uuid:            menu.ID.String(),
		OnDate:          timestamppb.New(menu.OnDate),
		OpeningRecordAt: timestamppb.New(menu.OpeningRecordAt),
		ClosingRecordAt: timestamppb.New(menu.ClosingRecordAt),
		Salads:          apiSalads,
		Garnishes:       apiGarnishes,
		Meats:           apiMeats,
		Soups:           apiSoups,
		Drinks:          apiDrinks,
		Desserts:        apiDesserts,
		CreatedAt:       timestamppb.New(menu.CreatedAt),
	}

	return &restaurant.GetMenuResponse{
		Menu: apiMenu,
	}, nil
}
