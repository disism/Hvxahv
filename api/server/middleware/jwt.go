package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"hvxahv/api/server/httputils"
	"log"
	"strings"
)

// JWTAuth JWT authentication middleware for gin network framework,
// Check whether the Token is carried in the request, and verify whether the Token is correct,
// Will obtain the username by parsing the Token and add the username in the context and set the key to loginUser.
func JWTAuth(c *gin.Context) {
	ht := c.Request.Header.Get("Authorization")
	t := strings.Split(ht, "Bearer ")[1]
	if ht == "" {
		c.JSON(500, gin.H{
			"state": "500",
			"message": "Token is not carried in the request.",
		})
		c.Abort()
		return
	}
	_, _, err := JwtParseToken(t)
	if err != nil {
		c.JSON(500, gin.H{
			"state": "500",
			"message": "Login failed. Token is incorrect!",
		})
		c.Abort()
	} else {
		u, err := JwtParseUser(t)
		if err != nil {
			log.Println("通过 token 获取用户失败")
		}
		c.Set("loginUser", u.User)
		c.Next()
	}

}

func JwtParseToken(tokenString string) (*jwt.Token, *httputils.Claims, error) {
	Claims := &httputils.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, Claims,
		func(token *jwt.Token) (i interface{}, err error) {
			return httputils.K, nil
		})
	if err != nil {
		log.Println("解 Token 失败！")
	}
	return token, Claims, err
}

func JwtParseUser(tokenString string) (*httputils.Claims, error) {
	if tokenString == "" {
		log.Println("需要传 Token ")
	}
	Claims := &httputils.Claims{}
	_, err := jwt.ParseWithClaims(tokenString, Claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(httputils.K), nil
		})
	if err != nil {
		return nil, err
	}
	return Claims, err
}