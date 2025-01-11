package service

import (
	"materials/internal/apperrors"
	"materials/internal/config"
	"materials/internal/dto"
	dao "materials/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	dbUser UserStorage
	cfg    config.Config
	log    zap.Logger
}

type UserStorage interface {
	CreateUser(dto dto.RegistrationUser) error
	GetUserByName(name string) (*dao.User, error)
}

type CustomClaims struct {
	UserID   int64  `json:"userID"`
	Username string `json:"userName"`
	jwt.RegisteredClaims
}

type GenerateJWTProps struct {
	Secret   []byte
	Exprires time.Time
	UserID   int64
	Username string
}

func NewUserService(
	userStorage UserStorage,
	cfg *config.Config,
	log *zap.Logger,
) *UserService {
	return &UserService{
		dbUser: userStorage,
		cfg:    *cfg,
		log:    *log,
	}
}

func (u *UserService) Registration(registrationDTO dto.RegistrationUser) (*string, *string, error) {
	hashedPassword, err := hashPassword(registrationDTO.Password, u.cfg.Cost)
	if err != nil {
		return nil, nil, apperrors.ErrHashPassword
	}

	registrationDTO.Password = hashedPassword

	err = u.dbUser.CreateUser(registrationDTO)
	if err != nil {
		u.log.Error(err.Error())
		return nil, nil, apperrors.ErrDBQuery
	}

	createdUser, err := u.dbUser.GetUserByName(registrationDTO.Username)
	if err != nil {
		u.log.Error(err.Error())
		return nil, nil, apperrors.ErrDBQuery
	}

	accessTime, err := time.ParseDuration(u.cfg.AccessLive)
	if err != nil {
		u.log.Error(err.Error())
		return nil, nil, err
	}

	JWTAccessProps := GenerateJWTProps{
		Secret:   []byte(u.cfg.JWTSecret),
		Exprires: time.Now().Add(accessTime),
		UserID:   int64(createdUser.ID),
		Username: createdUser.Username,
	}

	accessToken, err := generateJWT(JWTAccessProps)
	if err != nil {
		u.log.Error(err.Error())
		return nil, nil, apperrors.ErrJWTGeneration
	}

	refreshTime, err := time.ParseDuration(u.cfg.RefreshLive)
	if err != nil {
		u.log.Error(err.Error())
		return nil, nil, err
	}

	JWTRefreshProps := GenerateJWTProps{
		Secret:   []byte(u.cfg.JWTSecret),
		Exprires: time.Now().Add(refreshTime),
		UserID:   int64(createdUser.ID),
		Username: createdUser.Username,
	}

	refreshToken, err := generateJWT(JWTRefreshProps)
	if err != nil {
		u.log.Error(err.Error())
		return nil, nil, apperrors.ErrJWTGeneration
	}

	return &refreshToken, &accessToken, nil
}

func (u *UserService) Login(loginDTO dto.LoginUser) (string, string, error) {
	curUser, err := u.dbUser.GetUserByName(loginDTO.Username)
	if err != nil {
		u.log.Error(err.Error())
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(curUser.Password), []byte(loginDTO.Password))
	if err != nil {
		return "", "", apperrors.ErrInvalidPassword
	}

	accessTime, err := time.ParseDuration(u.cfg.AccessLive)
	if err != nil {
		u.log.Error(err.Error())
		return "", "", err
	}

	JWTAccessProps := GenerateJWTProps{
		Secret:   []byte(u.cfg.JWTSecret),
		Exprires: time.Now().Add(accessTime),
		UserID:   int64(curUser.ID),
		Username: curUser.Username,
	}

	accessToken, err := generateJWT(JWTAccessProps)
	if err != nil {
		u.log.Error(err.Error())
		return "", "", apperrors.ErrJWTGeneration
	}

	refreshTime, err := time.ParseDuration(u.cfg.RefreshLive)
	if err != nil {
		u.log.Error(err.Error())
		return "", "", apperrors.ErrJWTGeneration
	}

	JWTRefreshProps := GenerateJWTProps{
		Secret:   []byte(u.cfg.JWTSecret),
		Exprires: time.Now().Add(refreshTime),
		UserID:   int64(curUser.ID),
		Username: curUser.Username,
	}

	refreshToken, err := generateJWT(JWTRefreshProps)
	if err != nil {
		u.log.Error(err.Error())
		return "", "", apperrors.ErrJWTGeneration
	}

	return accessToken, refreshToken, nil
}

func hashPassword(password string, cost int) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func generateJWT(props GenerateJWTProps) (string, error) {
	claims := &CustomClaims{
		UserID:   props.UserID,
		Username: props.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(props.Exprires),
			Issuer:    "exampleIssuer",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(props.Secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
