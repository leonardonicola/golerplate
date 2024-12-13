package entity

import (
	"errors"
	"regexp"
	"time"

	"github.com/leonardonicola/golerplate/pkg/constants"
)

var (
	ErrInvalidEmail = errors.New(constants.ErrMsgInvalidEmail)
	ErrInvalidCPF   = errors.New(constants.ErrMsgInvalidCPF)
	ErrInvalidAge   = errors.New(constants.ErrMsgInvalidAge)
	ErrInvalidName  = errors.New(constants.ErrMsgInvalidName)
)

type User struct {
	Age       uint8      `json:"age" db:"age"`
	ID        string     `json:"id" db:"id, primarykey"`
	FullName  string     `json:"full_name" db:"full_name"`
	Password  string     `json:"-" db:"password"`
	Email     string     `json:"email" db:"email"`
	CPF       string     `json:"cpf" db:"cpf"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

func NewUser(fullname, email, cpf, password string, age uint8) (*User, error) {
	u := &User{
		FullName:  fullname,
		Email:     email,
		CPF:       cpf,
		Age:       age,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) Validate() error {
	if err := u.validateEmail(); err != nil {
		return err
	}
	if err := u.validateCPF(); err != nil {
		return err
	}
	if err := u.validateAge(); err != nil {
		return err
	}
	return nil
}

func (u *User) validateEmail() error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return ErrInvalidEmail
	}
	return nil
}

func (u *User) validateCPF() error {
	cpf := u.cleanCPF(u.CPF)
	if !u.isValidCPF(cpf) {
		return ErrInvalidCPF
	}
	return nil
}

func (u *User) validateAge() error {
	if u.Age < 0 || u.Age > 150 {
		return ErrInvalidAge
	}
	return nil
}

// cleanCPF removes any non-digit characters from the CPF
func (u *User) cleanCPF(cpf string) string {
	regex := regexp.MustCompile(`[^0-9]`)
	return regex.ReplaceAllString(cpf, "")
}

func (u *User) isValidCPF(cpf string) bool {
	if len(cpf) != 11 {
		return false
	}

	// Check if all digits are the same
	allEqual := true
	for i := 1; i < 11; i++ {
		if cpf[i] != cpf[0] {
			allEqual = false
			break
		}
	}
	if allEqual {
		return false
	}

	// First digit validation
	sum := 0
	for i := 0; i < 9; i++ {
		num := int(cpf[i] - '0')
		sum += num * (10 - i)
	}
	remainder := sum % 11
	if remainder < 2 {
		if int(cpf[9]-'0') != 0 {
			return false
		}
	} else {
		if int(cpf[9]-'0') != 11-remainder {
			return false
		}
	}

	// Second digit validation
	sum = 0
	for i := 0; i < 10; i++ {
		num := int(cpf[i] - '0')
		sum += num * (11 - i)
	}
	remainder = sum % 11
	if remainder < 2 {
		if int(cpf[10]-'0') != 0 {
			return false
		}
	} else {
		if int(cpf[10]-'0') != 11-remainder {
			return false
		}
	}

	return true
}
