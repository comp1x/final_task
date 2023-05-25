package userrepository

import (
	"context"
	"final-task/customer/internal/models"
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/customer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type UserService struct {
	customer.UnimplementedUserServiceServer

	db *gorm.DB
}

func New(dbURL string) (*UserService, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к базе данных: %w", err)
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
		log.Printf("not type uuid.UUID: %v", err)
		return nil, fmt.Errorf("not type uuid.UUID")
	}

	user := &models.User{
		Name:       request.Name,
		OfficeUuid: UuidFormatFromString,
	}

	if err = s.db.Create(user).Error; err != nil {
		log.Printf("ошибка при создании юзера в базе данных: %v", err)
		return nil, fmt.Errorf("ошибка при создании юзера")
	}

	return &customer.CreateUserResponse{}, nil
}

func (s *UserService) GetUserList(
	ctx context.Context, request *customer.GetUserListRequest,
) (*customer.GetUserListResponse, error) {
	if err := request.ValidateAll(); err != nil {
		log.Printf("ошибка при валидации: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	UuidFormatFromString, err := uuid.Parse(request.OfficeUuid)
	if err != nil {
		log.Printf("ошибка при конвертации: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var CurrentOffice models.Office

	if err := s.db.WithContext(ctx).Find(&CurrentOffice).Where("id = ?", UuidFormatFromString).Error; err != nil {
		log.Printf("ошибка при получении имени офиса из базы данных: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var users []models.User

	if err := s.db.Find(&users).Where("office_uuid = ?", UuidFormatFromString).Error; err != nil {
		log.Printf("ошибка при получении списка юзеров из базы данных: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
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
