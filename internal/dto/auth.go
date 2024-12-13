package dto

import "github.com/leonardonicola/golerplate/internal/domain/entity"

type RegisterUserDTO struct {
	FullName string `json:"full_name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	CPF      string `json:"cpf" binding:"required"`
	Age      int    `json:"age" binding:"required,min=18,max=150"`
	Password string `json:"password" binding:"required,min=6"`
}

type ErrorResponseDTO struct {
	Message string `json:"message"`
}

type TokenResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequestDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RegisterResponseDTO struct {
	User entity.User `json:"user"`
}
