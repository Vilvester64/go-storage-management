package services

import (
	"errors"
	"net/mail"
	"strings"

	"storage-manager-api/internal/config"
	"storage-manager-api/internal/models"
	"storage-manager-api/internal/repository"
	"storage-manager-api/internal/utils"
)

var (
	ErrInvalidInput       = errors.New("email dan password wajib diisi dengan format yang valid")
	ErrEmailAlreadyExists = errors.New("email sudah terdaftar")
	ErrInvalidCredentials = errors.New("email atau password salah")
)

// AuthService mendefinisikan kontrak logika bisnis untuk autentikasi.
type AuthService interface {
	Register(email, password string) (*models.User, error)
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

// NewAuthService membuat instance baru AuthService.
func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// Register memvalidasi input, memastikan email unik, melakukan hashing password, dan menyimpan user baru.
func (s *authService) Register(email, password string) (*models.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || len(password) < 6 {
		return nil, ErrInvalidInput
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, ErrInvalidInput
	}

	existingUser, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    email,
		Password: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login memverifikasi kredensial user dan mengembalikan token JWT.
func (s *authService) Login(email, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return "", ErrInvalidCredentials
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrInvalidCredentials
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(user.ID, user.Email, s.cfg.JWTSecret, s.cfg.JWTExpirationHours)
	if err != nil {
		return "", err
	}

	return token, nil
}
