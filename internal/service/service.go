package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/DblMOKRQ/auth-service/internal/entity"
	"github.com/DblMOKRQ/auth-service/internal/token"
	auth "github.com/DblMOKRQ/auth-service/pkg/api"
	"github.com/DblMOKRQ/auth-service/pkg/logger"
	"go.uber.org/zap"

	// "golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	userExists = "user already exists"
)

type Repository interface {
	Register(*entity.User) (int64, error)
	Login(user *entity.User, t *token.JWTMaker) (string, error)
	ValideToken(t *token.JWTMaker, token string) error
}

type Service struct {
	repo   Repository
	logger *logger.Logger
	t      *token.JWTMaker
	auth.UnimplementedAuthServer
}

func NewService(repo Repository, logger *logger.Logger, t *token.JWTMaker) *Service {
	return &Service{repo: repo, logger: logger, t: t}
}
func (s *Service) Register(ctx context.Context, in *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	id, err := s.repo.Register(&entity.User{
		Username: in.GetUsername(),
		Password: in.Password,
	})
	if err != nil {
		if errors.Is(err, entity.ErrUserExists) {
			s.logger.Error(userExists, zap.Error(err))

			err = status.Errorf(codes.AlreadyExists, userExists)
			return nil, err
		}
		s.logger.Error("failed to register user", zap.Error(err))
		return nil, fmt.Errorf("serv.reg.failed to register user: %v", err)
	}

	s.logger.Info("user registered", zap.Int64("id", id))
	return &auth.RegisterResponse{Id: id}, nil
}

func (s *Service) Login(ctx context.Context, in *auth.LoginRequest) (*auth.LoginResponse, error) {
	token, err := s.repo.Login(&entity.User{
		Username: in.GetUsername(),
		Password: in.Password,
	}, s.t)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) || errors.Is(err, entity.ErrInvalidPassword) {
			s.logger.Error("Invalid username or password", zap.Error(err))

			err = status.Errorf(codes.NotFound, "Invalid username or password")
			return nil, err
		}
		s.logger.Error("failed to login user", zap.Error(err))

		return nil, fmt.Errorf("failed to login user: %v", err)
	}

	s.logger.Info("user logged in", zap.String("token", token))
	return &auth.LoginResponse{Token: token}, nil
}

func (s *Service) Validate(ctx context.Context, in *auth.ValidateRequest) (*auth.ValidateResponse, error) {
	err := s.repo.ValideToken(s.t, in.Token)
	if err != nil {
		s.logger.Error("Invalid token", zap.Error(err))
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token")
	}

	return &auth.ValidateResponse{Valid: true}, nil
}
