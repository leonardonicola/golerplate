package service

import (
	"context"
	"errors"

	"github.com/leonardonicola/golerplate/internal/domain/entity"
	"github.com/leonardonicola/golerplate/internal/dto"
	"github.com/leonardonicola/golerplate/internal/infra/repository"
	"github.com/leonardonicola/golerplate/pkg/constants"
	"github.com/leonardonicola/golerplate/pkg/util"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type UserService interface {
	Create(ctx context.Context, dto dto.RegisterUserDTO) (*entity.User, error)
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByCPF(ctx context.Context, cpf string) (*entity.User, error)
	Authenticate(ctx context.Context, email, password string) (*entity.User, error)
}

type userService struct {
	repo   repository.UserRepository
	tracer oteltrace.Tracer
}

func NewUserService(r repository.UserRepository) *userService {
	return &userService{
		repo:   r,
		tracer: otel.Tracer(constants.TRACER_NAME),
	}
}

func (s *userService) Create(ctx context.Context, dto dto.RegisterUserDTO) (*entity.User, error) {
	ctx, span := s.tracer.Start(ctx, "CreateUser", oteltrace.WithAttributes(attribute.String("email", dto.Email)))
	defer span.End()

	ctx, hashSpan := s.tracer.Start(ctx, "HashPassword")
	hashedPw, err := util.HashPassword(dto.Password)
	hashSpan.End()
	if err != nil {
		return nil, err
	}
	ctx, entitySpan := s.tracer.Start(ctx, "CreateEntity")
	user, err := entity.NewUser(
		dto.FullName,
		dto.Email,
		dto.CPF,
		hashedPw,
		uint8(dto.Age),
	)
	entitySpan.End()

	if err != nil {
		return nil, err
	}

	user, err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*entity.User, error) {
	user, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetByCPF(ctx context.Context, cpf string) (*entity.User, error) {
	user, err := s.repo.GetByCPF(ctx, cpf)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Authenticate(ctx context.Context, email, password string) (*entity.User, error) {

	user, err := s.repo.GetByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	if !util.CheckPasswordEquality(password, user.Password) {
		return nil, errors.New(constants.ErrMsgInvalidCredentials)
	}

	user.Password = ""

	return user, nil
}
