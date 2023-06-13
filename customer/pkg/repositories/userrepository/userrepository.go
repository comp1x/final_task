package userrepository

import (
	"context"
	"github.com/comp1x/final-task/customer/pkg/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/customer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type UserService struct {
	customer.UnimplementedUserServiceServer
	db     *gorm.DB
	logger logrus.FieldLogger
}

func New(dbURL string, logger logrus.FieldLogger) (*UserService, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		logger.Error("New (UserService): ", err, time.Now().UTC())
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	return &UserService{
		db:     db,
		logger: logger,
	}, nil
}

func (s *UserService) CreateUser(
	ctx context.Context, request *customer.CreateUserRequest,
) (*customer.CreateUserResponse, error) {
	if err := request.ValidateAll(); err != nil {
		s.logger.Error("CreateUser: ", err, time.Now().UTC())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	UuidFormatFromString, err := uuid.Parse(request.OfficeUuid)

	if err != nil {
		s.logger.Error("CreateUser: ", err, time.Now().UTC())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user := models.User{
		Name:       request.Name,
		OfficeUuid: UuidFormatFromString,
	}

	if err = s.db.WithContext(ctx).Preload("Office").Create(&user).Error; err != nil {
		s.logger.Error("CreateUser: ", err, time.Now().UTC())
		return nil, status.Error(codes.Unknown, err.Error())
	}

	s.logger.Println("CreateUser: ", user, time.Now().UTC())

	return &customer.CreateUserResponse{}, nil
}

func (s *UserService) GetUserList(
	ctx context.Context, request *customer.GetUserListRequest,
) (*customer.GetUserListResponse, error) {
	if err := request.ValidateAll(); err != nil {
		s.logger.Error("GetUserList: ", err, time.Now().UTC())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	UuidFormatFromString, err := uuid.Parse(request.OfficeUuid)

	if err != nil {
		s.logger.Error("GetUserList: ", err, time.Now().UTC())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var CurrentOffice models.Office

	if err := s.db.WithContext(ctx).Where("id = ?", UuidFormatFromString).Find(&CurrentOffice).Error; err != nil {
		s.logger.Error("GetUserList: ", err, time.Now().UTC())
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	var users []models.User

	if err := s.db.WithContext(ctx).Where("office_uuid = ?", UuidFormatFromString).Find(&users).Error; err != nil {
		s.logger.Error("GetUserList: ", err, time.Now().UTC())
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	apiUsers := make([]*customer.User, 0, len(users))
	for _, user := range users {
		apiUser := &customer.User{
			Uuid:       user.ID.String(),
			Name:       user.Name,
			OfficeUuid: request.OfficeUuid,
			OfficeName: CurrentOffice.Name,
			CreatedAt:  timestamppb.New(user.CreatedAt),
		}
		apiUsers = append(apiUsers, apiUser)
	}

	s.logger.Println("GetUserList: ", apiUsers, time.Now().UTC())

	return &customer.GetUserListResponse{
		Result: apiUsers,
	}, nil
}
