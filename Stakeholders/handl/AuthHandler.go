package handl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"stakeholders_service/model"
	"stakeholders_service/proto/auth"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type AuthHandler struct {
	auth.UnimplementedStakeholderServiceServer
	DatabaseConnection *gorm.DB
	Key                string
}

func ConvertToString(role int) string {
	switch role {
	case 0:
		return "administrator"
	case 1:
		return "author"
	case 2:
		return "tourist"
	default:
		return "unknown"
	}
}

/*func generateAccessToken(user model.User, person model.Person, key string) (string, error) {
	claims := jwt.MapClaims{
		"jti":      uuid.New().String(),
		"id":       user.ID,
		"username": user.Username,
		"personId": person.ID,
		"role":     ConvertToString(user.Role),
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}*/

func generateAccessToken(user model.User, person model.Person, key string) (string, error) {
	claims := jwt.MapClaims{
		"jti":      uuid.New().String(),
		"id":       user.ID,
		"username": user.Username,
		"personId": person.ID,
		"http://schemas.microsoft.com/ws/2008/06/identity/claims/role": ConvertToString(user.Role),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iss": "explorer",
		"aud": "explorer-front.com",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func createToken(claims jwt.MapClaims, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (h AuthHandler) LogIn(ctx context.Context, request *auth.RequestLogIn) (*auth.ResponseLogIn, error) {
	var user model.User
	if err := h.DatabaseConnection.Table(`stakeholders."Users"`).Where(`"Users"."Username" = ? and "Users"."IsActive" = true`, request.Username).First(&user).Error; err != nil {
		return nil, err
	}

	var person model.Person
	if err := h.DatabaseConnection.Table(`stakeholders."People"`).Where(`"People"."UserId" = ?`, user.ID).First(&person).Error; err != nil {
		return nil, err
	}

	tokenString, _ := generateAccessToken(user, person, h.Key)

	return &auth.ResponseLogIn{
		Id:          user.ID,
		AccessToken: tokenString,
	}, nil
}

func (h AuthHandler) RegisterTourist(ctx context.Context, request *auth.RequestRegister) (*auth.ResponseLogIn, error) {

	var user model.User = model.User{
		Username: request.Username,
		Password: request.Password,
		Role:     2,
		IsActive: false,
	}

	dbResult := h.DatabaseConnection.Table(`stakeholders."Users"`).Create(&user)

	if dbResult.Error != nil {
		return nil, dbResult.Error
	}

	var person model.Person = model.Person{
		UserID:  user.ID,
		Name:    request.Name,
		Surname: request.Surname,
		Email:   request.Email,
	}
	dbResultPerson := h.DatabaseConnection.Table(`stakeholders."People"`).Create(&person)

	if dbResultPerson.Error != nil {
		return nil, dbResultPerson.Error
	}
	tokenString, _ := generateAccessToken(user, person, h.Key)

	user.EmailVerificationToken = &tokenString
	dbResultToken := h.DatabaseConnection.Table(`stakeholders."Users"`).Save(user)

	if dbResultToken.Error != nil {
		return nil, dbResultToken.Error
	}

	var wallet model.Wallet = model.Wallet{
		UserId:  user.ID,
		Balance: 0,
	}
	dbResultWallet := h.DatabaseConnection.Table(`payments."Wallet"`).Create(&wallet)

	if dbResultWallet.Error != nil {
		return nil, dbResultWallet.Error
	}
	var userExperience model.UserExperience = model.UserExperience{
		UserID: user.ID,
		XP:     0,
		Level:  1,
	}

	dbResultUserExp := h.DatabaseConnection.Table(`encounters."UserExperience"`).Create(&userExperience)

	if dbResultUserExp.Error != nil {
		return nil, dbResultUserExp.Error
	}

	tokenStringAccess, _ := generateAccessToken(user, person, h.Key)

	return &auth.ResponseLogIn{
		Id:          user.ID,
		AccessToken: tokenStringAccess,
	}, nil
}

func (h AuthHandler) ActivateUser(ctx context.Context, request *auth.RequestActivateUser) (*auth.ResponseLogIn, error) {

	token, err := jwt.Parse(request.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Key), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	id, ok := claims["id"]
	if !ok {
		return nil, fmt.Errorf("invalid id in token")
	}

	var user model.User
	if err := h.DatabaseConnection.Table(`stakeholders."Users"`).First(&user, id).Error; err != nil {
		return nil, err
	}

	if user.EmailVerificationToken == nil || *user.EmailVerificationToken != request.Token {
		return nil, fmt.Errorf("invalid token")
	}

	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)

	if time.Now().After(expirationTime) {
		return nil, fmt.Errorf("token has expired")
	}

	user.IsActive = true
	user.EmailVerificationToken = nil
	dbResultToken := h.DatabaseConnection.Table(`stakeholders."Users"`).Save(user)

	if dbResultToken.Error != nil {
		return nil, dbResultToken.Error
	}

	var person model.Person
	if err := h.DatabaseConnection.Table(`stakeholders."People"`).Where(`"People"."UserId" = ?`, user.ID).First(&person).Error; err != nil {
		return nil, err
	}

	tokenStringAccess, _ := generateAccessToken(user, person, h.Key)

	return &auth.ResponseLogIn{
		Id:          user.ID,
		AccessToken: tokenStringAccess,
	}, nil
}

func (h AuthHandler) ChangePassword(ctx context.Context, request *auth.RequestChangePassword) (*auth.RequestActivateUser, error) {
	token, err := jwt.Parse(request.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Key), nil
	})
	if err != nil {
		return nil, err
	}

	if request.NewPassword != request.NewPasswordConfirm {
		return nil, errors.New("passwords do not match")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	userId, ok := claims["id"]
	if !ok {
		return nil, fmt.Errorf("invalid id in token")
	}

	var user model.User
	if err := h.DatabaseConnection.Table(`stakeholders."Users"`).Where(`"Users"."Id" = ?`, userId).First(&user).Error; err != nil {
		return nil, err
	}

	if user.ResetPasswordToken == nil || *user.ResetPasswordToken != request.Token {
		return nil, errors.New("invalid token")
	}

	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	if time.Now().After(expirationTime) {
		return nil, fmt.Errorf("token has expired")
	}

	user.Password = request.NewPassword
	user.ResetPasswordToken = nil

	dbResultToken := h.DatabaseConnection.Table(`stakeholders."Users"`).Save(user)
	if dbResultToken.Error != nil {
		return nil, dbResultToken.Error
	}

	return &auth.RequestActivateUser{
		Token: "Password successfully changed",
	}, nil
}

func (h AuthHandler) ChangePasswordRequest(ctx context.Context, request *auth.RequestChangePasswordRequest) (*auth.RequestActivateUser, error) {
	var person model.Person
	if err := h.DatabaseConnection.Table(`stakeholders."People"`).Where(`"People"."Email" = ?`, request.Email).First(&person).Error; err != nil {
		return nil, err
	}

	var user model.User
	if err := h.DatabaseConnection.Table(`stakeholders."Users"`).Where(`"Users"."Id" = ?`, person.UserID).First(&user).Error; err != nil {
		return nil, err
	}

	var claims = jwt.MapClaims{
		"jti":      uuid.New().String(),
		"id":       user.ID,
		"username": user.Username,
		"personId": person.ID,
		"role":     ConvertToString(user.Role),
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	jwt, err := createToken(claims, h.Key)
	if err != nil {
		return nil, err
	}

	user.ResetPasswordToken = &jwt
	dbResultToken := h.DatabaseConnection.Table(`stakeholders."Users"`).Save(user)
	if dbResultToken.Error != nil {
		return nil, dbResultToken.Error
	}

	return &auth.RequestActivateUser{
		Token: jwt,
	}, nil
}
func Conn() *nats.Conn {
	conn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

type Message struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func (h AuthHandler) ChangeRating(id int, userId int) {
	var applicationRating model.ApplicationRating
	if err := h.DatabaseConnection.Table(`stakeholders."ApplicationRatings"`).Where(`"ApplicationRatings"."UserId" = ?`, userId).First(&applicationRating).Error; err != nil {
		return
	}
	conn := Conn()
	if applicationRating.Grade < 8 {

		messageRec := Message{
			Id:   id,
			Body: "Failed",
		}
		data, err := json.Marshal(messageRec)
		if err != nil {
			log.Fatal(err)
		}
		errTours := conn.Publish("subTours", data)
		if errTours != nil {
			log.Fatal(errTours)
		}
	} else {

		applicationRating.Grade = 10
		if err := h.DatabaseConnection.Table(`stakeholders."ApplicationRatings"`).Save(&applicationRating).Error; err != nil {
			fmt.Println("Failed to update record:", err)
		}
		messageRec := Message{
			Id:   id,
			Body: "Success",
		}
		data, err := json.Marshal(messageRec)
		if err != nil {
			log.Fatal(err)
		}
		errTours := conn.Publish("subTours", data)
		if errTours != nil {
			log.Fatal(errTours)
		}
	}
	return
}
