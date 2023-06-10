package statisticsrepository

import (
	"context"
	"github.com/comp1x/final-task/restaurant/pkg/models"
	"github.com/google/uuid"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/statistics"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sort"
)

type StatisticsService struct {
	statistics.UnimplementedStatisticsServiceServer
	db *gorm.DB
}

func New(dbURL string) (*StatisticsService, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
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

	timeStart := request.GetStartDate().AsTime()
	timeEnd := request.GetEndDate().AsTime()

	var orders []models.Order
	if err := s.db.WithContext(ctx).Where("created_at BETWEEN ? AND ?", timeStart, timeEnd).Find(&orders).Error; err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}
	var profit float64

	for _, order := range orders {
		var product *models.Product
		if err := s.db.WithContext(ctx).Select("price").First(&product, order.ProductUuid).Error; err != nil {
			return nil, status.Error(codes.Unavailable, err.Error())
		}
		profit += product.Price * float64(order.Count)
	}
	return &statistics.GetAmountOfProfitResponse{
		Profit: profit,
	}, nil
}

func (s *StatisticsService) TopProducts(
	ctx context.Context, request *statistics.TopProductsRequest,
) (*statistics.TopProductsResponse, error) {

	var orders []models.Order

	if request.StartDate != nil && request.EndDate != nil {
		timeStart := request.GetStartDate().AsTime()
		timeEnd := request.GetEndDate().AsTime()

		if err := s.db.WithContext(ctx).Where("created_at BETWEEN ? AND ?", timeStart, timeEnd).Find(&orders).Error; err != nil {
			return nil, status.Error(codes.Unavailable, err.Error())
		}
	} else if request.StartDate != nil && request.EndDate == nil {
		timeStart := request.GetStartDate().AsTime()

		if err := s.db.WithContext(ctx).Where("created_at > ?", timeStart).Find(&orders).Error; err != nil {
			return nil, status.Error(codes.Unavailable, err.Error())
		}
	} else if request.EndDate != nil && request.StartDate == nil {
		timeEnd := request.GetStartDate().AsTime()

		if err := s.db.WithContext(ctx).Where("created_at < ?", timeEnd).Find(&orders).Error; err != nil {
			return nil, status.Error(codes.Unavailable, err.Error())
		}
	} else if request.EndDate == nil && request.StartDate == nil {
		if err := s.db.WithContext(ctx).Find(&orders).Error; err != nil {
			return nil, status.Error(codes.Unavailable, err.Error())
		}
	}

	productsMap := make(map[uuid.UUID]int64)

	for _, order := range orders {
		productsMap[order.ProductUuid] += order.Count
	}

	var productCounts []ProductCount
	for key, value := range productsMap {
		productCounts = append(productCounts, ProductCount{
			key,
			value,
		})
	}

	sort.Sort(ByCountDesc(productCounts))

	apiProducts := make([]*statistics.Product, 0, len(productsMap))
	for _, productWithCount := range productCounts {
		var product *models.Product
		if err := s.db.WithContext(ctx).Select("name", "type").First(&product, productWithCount.ProductUUID).Error; err != nil {
			return nil, status.Error(codes.Unavailable, err.Error())
		}
		apiProduct := &statistics.Product{
			Uuid:        productWithCount.ProductUUID.String(),
			Name:        product.Name,
			Count:       productWithCount.Count,
			ProductType: statistics.StatisticsProductType(product.Type),
		}
		apiProducts = append(apiProducts, apiProduct)
	}

	return &statistics.TopProductsResponse{
		Result: apiProducts,
	}, nil
}

type ProductCount struct {
	ProductUUID uuid.UUID
	Count       int64
}

type ByCountDesc []ProductCount

func (a ByCountDesc) Len() int           { return len(a) }
func (a ByCountDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCountDesc) Less(i, j int) bool { return a[i].Count > a[j].Count }
