package service

import (
	"context"

	"github.com/porseOnline/api/pb"
	"github.com/porseOnline/internal/user"
	"github.com/porseOnline/internal/user/domain"
	userPort "github.com/porseOnline/internal/user/port"
	"github.com/porseOnline/pkg/jwt"
	"github.com/porseOnline/pkg/time"

	jwt2 "github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	svc                   userPort.Service
	authSecret            string
	expMin, refreshExpMin uint
}

func NewUserService(svc userPort.Service, authSecret string, expMin, refreshExpMin uint) *UserService {
	return &UserService{
		svc:           svc,
		authSecret:    authSecret,
		expMin:        expMin,
		refreshExpMin: refreshExpMin,
	}
}

var (
	ErrUserCreationValidation = user.ErrUserCreationValidation
	ErrUserOnCreate           = user.ErrUserOnCreate
	ErrUserNotFound           = user.ErrUserNotFound
)

func (s *UserService) SignUp(ctx context.Context, req *pb.UserSignUpFirstRequest) (*pb.UserSignUpFirstResponse, error) {
	userID, err := s.svc.CreateUser(ctx, domain.User{
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		Phone:     domain.Phone(req.GetPhone()),
	})

	if err != nil {
		return nil, err
	}

	return &pb.UserSignUpFirstResponse{
		UserId: uint64(userID),
	}, nil
}
func (s *UserService) SignUpCodeVerification(ctx context.Context, req *pb.UserSignUpSecondRequest) (*pb.UserSignUpSecondResponse, error) {
	_, err := s.svc.GetUserByID(ctx, domain.UserID(req.GetUserId()))
	if err != nil {
		return nil, err
	}

	accessToken, err := jwt.CreateToken([]byte(s.authSecret), &jwt.UserClaims{
		RegisteredClaims: jwt2.RegisteredClaims{
			ExpiresAt: jwt2.NewNumericDate(time.AddMinutes(s.expMin, true)),
		},
		UserID: uint(req.GetUserId()),
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.CreateToken([]byte(s.authSecret), &jwt.UserClaims{
		RegisteredClaims: jwt2.RegisteredClaims{
			ExpiresAt: jwt2.NewNumericDate(time.AddMinutes(s.refreshExpMin, true)),
		},
		UserID: uint(req.GetUserId()),
	})

	if err != nil {
		return nil, err
	}

	return &pb.UserSignUpSecondResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (s *UserService) GetByID(ctx context.Context, id uint) (*pb.User, error) {
	user, err := s.svc.GetUserByID(ctx, domain.UserID(id))
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id:        uint64(user.ID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     string(user.Phone),
	}, nil
}
