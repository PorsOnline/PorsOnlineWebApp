package user

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/porseOnline/internal/user/domain"
	"github.com/porseOnline/internal/user/port"
	"github.com/porseOnline/pkg/logger"
)

var (
	ErrUserOnCreate           = errors.New("error on creating new user")
	ErrUserCreationValidation = errors.New("validation failed")
	ErrUserNotFound           = errors.New("user not found")
)

type service struct {
	repo port.Repo
}

func NewService(repo port.Repo) port.Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateUser(ctx context.Context, user domain.User) (domain.UserID, error) {
	if err := user.Validate(); err != nil {
		return 0, fmt.Errorf("%w %w", ErrUserCreationValidation, err)
	}

	userID, err := s.repo.Create(ctx, user)
	if err != nil {
		log.Println("error on creating new user : ", err.Error())
		return 0, ErrUserOnCreate
	}

	return userID, nil
}

func (s *service) GetUserByID(ctx context.Context, userID domain.UserID) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil || user.ID == 0 {
		return nil, ErrUserNotFound
	}

	return user, nil
}
func (s *service) GetUserByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil || user.ID == 0 {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *service) UpdateUser(ctx context.Context, user domain.User) error {
	err := s.repo.UpdateUser(ctx, user)
	if err != nil {
		logger.Error("error in update user", nil)
		return err
	}
	return nil
}

func (s *service) DeleteByID(ctx context.Context, userID domain.UserID) error{
	err:=s.repo.DeleteByID(ctx, userID)
	if err!=nil{
		logger.Error("can not delete user", nil)
		return err
	}
	return nil
}
