package userrepository

import (
	"context"
	"github.com/comp1x/final-task/customer/pkg/models"
	"github.com/google/uuid"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/customer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserService struct {
	customer.UnimplementedUserServiceServer

	db *gorm.DB
}

func New(dbURL string) (*UserService, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	return &UserService{
		db: db,
	}, nil
}

func (s *UserService) CreateUser(
	ctx context.Context, request *customer.CreateUserRequest,
) (*customer.CreateUserResponse, error) {
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	UuidFormatFromString, err := uuid.Parse(request.OfficeUuid)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user := &models.User{
		Name:       request.Name,
		OfficeUuid: UuidFormatFromString,
	}

	if err = s.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &customer.CreateUserResponse{}, nil
}

func (s *UserService) GetUserList(
	ctx context.Context, request *customer.GetUserListRequest,
) (*customer.GetUserListResponse, error) {
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	UuidFormatFromString, err := uuid.Parse(request.OfficeUuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var CurrentOffice models.Office

	if err := s.db.WithContext(ctx).Where("id = ?", UuidFormatFromString).Find(&CurrentOffice).Error; err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	var users []models.User

	if err := s.db.WithContext(ctx).Where("office_uuid = ?", UuidFormatFromString).Find(&users).Error; err != nil {
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

	return &customer.GetUserListResponse{
		Result: apiUsers,
	}, nil
}
