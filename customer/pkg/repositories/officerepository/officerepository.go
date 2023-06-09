package officerepository

import (
	"context"
	"github.com/comp1x/final-task/customer/pkg/models"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/customer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type OfficeService struct {
	customer.UnimplementedOfficeServiceServer

	db *gorm.DB
}

func New(dbURL string) (*OfficeService, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	return &OfficeService{
		db: db,
	}, nil
}

func (s *OfficeService) CreateOffice(
	ctx context.Context, request *customer.CreateOfficeRequest,
) (*customer.CreateOfficeResponse, error) {
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	office := &models.Office{
		Name:    request.Name,
		Address: request.Address,
		//CreatedAt: TimeToTimestamp(time.Now()),
	}

	if err := s.db.Create(office).Error; err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	return &customer.CreateOfficeResponse{}, nil
}

func (s *OfficeService) GetOfficeList(
	ctx context.Context, request *customer.GetOfficeListRequest,
) (*customer.GetOfficeListResponse, error) {
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var offices []models.Office
	if err := s.db.Find(&offices).Error; err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	apiOffices := make([]*customer.Office, 0, len(offices))
	for _, office := range offices {
		apiOffice := &customer.Office{
			Uuid:      office.ID.String(),
			Name:      office.Name,
			Address:   office.Address,
			CreatedAt: timestamppb.New(office.CreatedAt),
		}
		apiOffices = append(apiOffices, apiOffice)
	}

	return &customer.GetOfficeListResponse{
		Result: apiOffices,
	}, nil
}
