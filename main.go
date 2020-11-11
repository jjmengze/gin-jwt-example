package main

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	r := RoutersInit()
	if err := http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		log.Fatal(err)
	}
}

func RoutersInit() *gin.Engine {
	var router = gin.Default()
	user := router.Group("/user").Use(Auth())
	router.POST("/login", login)
	user.GET("/hello", userHello)
	return router
}

type UserClaims struct {
	Username string
	jwt.StandardClaims
}

type JWT struct {
	SigningKey []byte
}

func NewJWT() *JWT {
	return &JWT{
		[]byte("IAMAGOODKEY"),
	}
}

func (j *JWT) GenerateToken(user UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(t string) (*UserClaims, error) {
	if !strings.Contains(t, "Bearer ") {
		return nil, errors.New("Authorization filed has to contain Bearer key word.")
	}
	input := strings.Split(strings.TrimSpace(t), "Bearer ")
	token, err := jwt.ParseWithClaims(input[1], &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if v, ok := err.(*jwt.ValidationError); ok {
			if v.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("Error Malformed token")
			} else if v.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, errors.New("Token is expired")
			} else if v.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("Token is not active")
			} else {
				return nil, errors.New("Couldn't handle token")
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, errors.New("Couldn't handle token")
	} else {
		return nil, errors.New("Couldn't handle token\"")
	}
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "who are  you ?"})
			c.Abort()
			return
		}
		j := NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something error happened"})
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"exp,omitempty"`
}

func login(c *gin.Context) {
	var req LoginRequest
	_ = c.ShouldBindJSON(&req)
	if req.Name == "" || req.Password == "" {
		c.JSON(http.StatusUnauthorized, "please don't empty user name/password")
		return
	}
	GenerateTokenForUser(c, &req)
}

func GenerateTokenForUser(c *gin.Context, u *LoginRequest) {
	j := NewJWT()
	claims := UserClaims{
		Username: u.Name,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,    // effective on a time
			ExpiresAt: time.Now().Unix() + 60*60*2, // Expires on/from 2 hr .
			Issuer:    "jason",                     // Issuer
		},
	}
	token, err := j.GenerateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "GenToken error")
		return
	}
	c.JSON(http.StatusOK, LoginResponse{Token: token, ExpiresAt: claims.ExpiresAt * 1000})
	return
}

func userHello(c *gin.Context) {
	c.String(http.StatusOK, "welcome you who has right token")
}
