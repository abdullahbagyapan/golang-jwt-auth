package service

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"jwtauth/model"
	"time"
)

const secretKey = "secret"

var list = map[string]string{} // temporary db

var tokenList = make(map[string]string, 10) // // temporary db

func GenerateToken(ctx *fiber.Ctx) map[string]string {
	user := new(model.User)

	err := ctx.BodyParser(&user)

	if err != nil {
		return nil
	}

	jwtware.New(jwtware.Config{SigningKey: []byte(secretKey)})

	claims := setClaimsByUser(*user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, _ := token.SignedString([]byte("secret"))

	list[user.Username] = user.Password
	tokenList[user.Username] = t

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
	if !checkClaimUser(token) {
		return false
	}
	return true

}

func checkClaimUser(tokenString string) bool {
	//token := getToken(tokenString)
	//
	//claimsMap := token.Claims.(jwt.MapClaims)
	//
	//username := claimsMap["username"]

	//TODO:	check from db and compare request token and user token
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
	//claimsMap := token.Claims.(jwt.MapClaims)
	//
	//username := claimsMap["username"]

	return token
}
