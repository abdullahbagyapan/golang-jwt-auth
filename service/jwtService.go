package service

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	database "jwtauth/database"
	"jwtauth/model"
	"time"
)

const secretKey = "secret"

func Register(ctx *fiber.Ctx) (map[string]string, error) {
	user := new(model.User)

	err := ctx.BodyParser(&user)

	if err != nil {
		return nil, err
	}

	user.ID = uuid.New()

	database.Database.Table("users").Create(&user)

	token := generateToken(*user)

	return token, nil

}

func generateToken(user model.User) map[string]string {

	jwtware.New(jwtware.Config{SigningKey: []byte(secretKey)})

	claims := setClaimsByUser(user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, _ := token.SignedString([]byte("secret"))

	var dbToken model.Token

	dbToken.ID = uuid.New()
	dbToken.Token = t
	dbToken.UserId = user.ID

	database.Database.Table("tokens").Create(&dbToken)

	return map[string]string{"token": t}
}

func setClaimsByUser(user model.User) jwt.MapClaims {

	claims := jwt.MapClaims{
		"username": user.Username,
		"admin":    false,
		"exp":      time.Now().Add(time.Hour * 6).Unix(), // 6 HOUR
	}

	return claims
}

func CheckToken(ctx *fiber.Ctx) bool {

	exist, token := getTokenString(ctx)

	if !exist {
		return false
	}
	if !isTokenValid(token) {
		return false
	}
	return true

}

func getTokenString(ctx *fiber.Ctx) (bool, string) {

	token := string(ctx.Request().Header.Peek("Token")) // return key's value

	if token != "" {
		return true, token
	}
	return false, ""

}

func isTokenValid(tokenString string) bool {

	token := getToken(tokenString)

	if token != nil {
		return token.Valid
	}
	return false

}

func getToken(tokenString string) *jwt.Token {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil
	}
	return token
}

func LoginToken(ctx *fiber.Ctx) (map[string]string, error) {

	var request model.User
	var user model.User

	ctx.BodyParser(&request)

	result := database.Database.Table("users").First(&user, "username=? AND password =?", request.Username, request.Password)

	if result.Error != nil {
		return nil, result.Error
	}

	var tokenModel model.Token

	result = database.Database.Table("tokens").First(&tokenModel, "user_id=?", user.ID)

	if result.Error != nil {
		return nil, result.Error
	}

	isValid := isTokenValid(tokenModel.Token)

	if !isValid {
		return generateToken(user), nil
	}

	returnMap := make(map[string]string)

	returnMap["token"] = tokenModel.Token

	return returnMap, nil

}
