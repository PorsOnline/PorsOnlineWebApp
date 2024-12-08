package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/porseOnline/api/pb"
	codeVerficationDomain "github.com/porseOnline/internal/codeVerification/domain"
	codeVerificationPort "github.com/porseOnline/internal/codeVerification/port"
	"github.com/porseOnline/internal/user"
	"github.com/porseOnline/internal/user/domain"
	userPort "github.com/porseOnline/internal/user/port"
	"github.com/porseOnline/pkg/helper"
	"github.com/porseOnline/pkg/jwt"
	helperTime "github.com/porseOnline/pkg/time"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"

	jwt2 "github.com/golang-jwt/jwt/v5"
)

var (
	ErrPasswordNotMatch = errors.New("not match password")
)

type UserService struct {
	svc                    userPort.Service
	authSecret             string
	expMin, refreshExpMin  uint
	codeVerficationServise codeVerificationPort.Service
}

func NewUserService(svc userPort.Service, authSecret string, expMin, refreshExpMin uint, codeVerificationSvc codeVerificationPort.Service) *UserService {
	return &UserService{
		svc:                    svc,
		authSecret:             authSecret,
		expMin:                 expMin,
		refreshExpMin:          refreshExpMin,
		codeVerficationServise: codeVerificationSvc,
	}
}

var (
	ErrUserCreationValidation = user.ErrUserCreationValidation
	ErrUserOnCreate           = user.ErrUserOnCreate
	ErrUserNotFound           = user.ErrUserNotFound
)

type SignUpFirstResponseWrapper struct {
	RequestTimestamp int64                       `json:"requestTimestamp"`
	Data             *pb.UserSignUpFirstResponse `json:"data"`
}
type SignUpSecondResponseWrapper struct {
	RequestTimestamp int64                        `json:"requestTimestamp"`
	Data             *pb.UserSignUpSecondResponse `json:"data"`
}

func (s *UserService) SignUp(ctx context.Context, req *pb.UserSignUpFirstRequest) (*SignUpFirstResponseWrapper, error) {
	userID, err := s.svc.CreateUser(ctx, domain.User{
		FirstName:    req.GetFirstName(),
		LastName:     req.GetLastName(),
		Phone:        domain.Phone(req.GetPhone()),
		Email:        domain.Email(req.GetEmail()),
		PasswordHash: req.GetPassword(),
		NationalCode: domain.NationalCode(req.GetNationalCode()),
		BirthDate:    req.GetBirthdate().AsTime(),
		City:         req.GetCity(),
		Gender:       req.GetGender(),
	})

	if err != nil {
		return nil, err
	}

	code := strconv.Itoa(helper.GetRandomCode())

	s.codeVerficationServise.Send(ctx, codeVerficationDomain.NewCodeVerification(userID, fmt.Sprint(code), codeVerficationDomain.CodeVerificationTypeEmail, true, time.Minute*2))

	// go helper.SendEmail(req.GetEmail())
	// go helper.SendEmail(req.GetEmail(), code)
	response := &SignUpFirstResponseWrapper{
		RequestTimestamp: time.Now().Unix(),
		Data: &pb.UserSignUpFirstResponse{
			UserId: uint64(userID),
		},
	}

	return response, nil
}

func (s *UserService) SignUpCodeVerification(ctx context.Context, req *pb.UserSignUpSecondRequest) (*SignUpSecondResponseWrapper, error) {
	_, err := s.svc.GetUserByID(ctx, domain.UserID(req.GetUserId()))
	if err != nil {
		return nil, err
	}
	ok, err := s.codeVerficationServise.CheckUserCodeVerificationValue(ctx, domain.UserID(req.GetUserId()), req.GetCode())
	if err != nil {
		return nil, err
	}
	if ok {

		accessToken, err := jwt.CreateToken([]byte(s.authSecret), &jwt.UserClaims{
			RegisteredClaims: jwt2.RegisteredClaims{
				ExpiresAt: jwt2.NewNumericDate(helperTime.AddMinutes(s.expMin, true)),
			},
			UserID: uint(req.GetUserId()),
		})
		if err != nil {
			return nil, err
		}

		refreshToken, err := jwt.CreateToken([]byte(s.authSecret), &jwt.UserClaims{
			RegisteredClaims: jwt2.RegisteredClaims{
				ExpiresAt: jwt2.NewNumericDate(helperTime.AddMinutes(s.refreshExpMin, true)),
			},
			UserID: uint(req.GetUserId()),
		})

		if err != nil {
			return nil, err
		}

		response := &SignUpSecondResponseWrapper{
			RequestTimestamp: time.Now().Unix(), // Get current UNIX timestamp
			Data: &pb.UserSignUpSecondResponse{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
		}
		return response, nil
	} else {
		return nil, nil
	}

}
func (s *UserService) SignIn(ctx context.Context, req *pb.UserSignInRequest) (*SignUpSecondResponseWrapper, error) {
	user, err := s.svc.GetUserByEmail(ctx, domain.Email(req.GetEmail()))
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.GetPassword()))
	if err != nil {
		return nil, ErrPasswordNotMatch
	}
	accessToken, err := jwt.CreateToken([]byte(s.authSecret), &jwt.UserClaims{
		RegisteredClaims: jwt2.RegisteredClaims{
			ExpiresAt: jwt2.NewNumericDate(helperTime.AddMinutes(s.expMin, true)),
		},
		UserID: uint(user.ID),
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.CreateToken([]byte(s.authSecret), &jwt.UserClaims{
		RegisteredClaims: jwt2.RegisteredClaims{
			ExpiresAt: jwt2.NewNumericDate(helperTime.AddMinutes(s.refreshExpMin, true)),
		},
		UserID: uint(user.ID),
	})

	if err != nil {
		return nil, err
	}

	response := &SignUpSecondResponseWrapper{
		RequestTimestamp: time.Now().Unix(), // Get current UNIX timestamp
		Data: &pb.UserSignUpSecondResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}
	return response, nil

}

func (s *UserService) GetByID(ctx context.Context, id uint) (*pb.User, error) {
	user, err := s.svc.GetUserByID(ctx, domain.UserID(id))
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id:                uint64(user.ID),
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Phone:             string(user.Phone),
		Email:             string(user.Email),
		PasswordHash:      user.PasswordHash,
		NationalCode:      string(user.NationalCode),
		BirthDate:         timestamppb.New(user.BirthDate), // Converts time.Time to protobuf Timestamp
		City:              user.City,
		Gender:            user.Gender,
		SurveyLimitNumber: int32(user.SurveyLimitNumber), // Protobuf may require int32 instead of int
		CreatedAt:         timestamppb.New(user.CreatedAt),
		DeletedAt:         timestamppb.New(user.DeletedAt), // Handle DeletedAt if needed
		UpdatedAt:         timestamppb.New(user.UpdatedAt),
		Balance:           int32(user.Balance),
	}, nil
}
