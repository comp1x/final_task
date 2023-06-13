package officerepository

import (
	"context"
	"github.com/comp1x/final-task/customer/pkg/models"
	"github.com/sirupsen/logrus"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/customer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type OfficeService struct {
	customer.UnimplementedOfficeServiceServer
	db     *gorm.DB
	logger logrus.FieldLogger
}

func New(dbURL string, logger logrus.FieldLogger) (*OfficeService, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		logger.Error("New (OfficeService): ", err, time.Now().UTC())
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	return &OfficeService{
		db:     db,
		logger: logger,
	}, nil
}

func (s *OfficeService) CreateOffice(
	ctx context.Context, request *customer.CreateOfficeRequest,
) (*customer.CreateOfficeResponse, error) {
	if err := request.ValidateAll(); err != nil {
		s.logger.Error("CreateOffice: ", err, time.Now().UTC())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	office := models.Office{
		Name:    request.Name,
		Address: request.Address,
	}

	if err := s.db.WithContext(ctx).Create(&office).Error; err != nil {
		s.logger.Error("CreateOffice: ", err, time.Now().UTC())
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	s.logger.Println("CreateOffice: ", office)

	return &customer.CreateOfficeResponse{}, nil
}

func (s *OfficeService) GetOfficeList(
	ctx context.Context, request *customer.GetOfficeListRequest,
) (*customer.GetOfficeListResponse, error) {
	if err := request.ValidateAll(); err != nil {
		s.logger.Error("CreateOffice: ", err, time.Now().UTC())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var offices []models.Office
	if err := s.db.WithContext(ctx).Find(&offices).Error; err != nil {
		s.logger.Error("CreateOffice: ", err, time.Now().UTC())
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

	s.logger.Println("GetOfficeList: ", apiOffices, time.Now().UTC())

	return &customer.GetOfficeListResponse{
		Result: apiOffices,
	}, nil
}
