package presentation

import (
	"log"
	"net/http"
	"strings"

	"amalhanaja.com/user-service/data"
	"amalhanaja.com/user-service/domain"

	"github.com/gin-gonic/gin"
)

// Start starting controller
func Start() {
	router := gin.Default()
	router.Use(CORSMiddleware())
	router.Use(gin.Recovery())
	router.Use(gin.Recovery())
	router.PUT("/register", RegisterController)
	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	router.PATCH("/checkAvailableEmail", CheckAvailableEmailAddress)
	router.POST("/activate", ActivateAccount)
	router.POST("/login", Login)
	me := router.Group("/me")
	me.Use(JWTAuthMiddleware())
	{
		me.GET("/aku", Profile)
	}
	router.Run(":8001")
}

// CORSMiddleware ...
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// JWTAuthMiddleware ...
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		useCase := domain.NewUserUseCase(data.GetInstance())
		authHeaderKey := "Authorization"
		authBearer := "Bearer"
		tokens := strings.Split(c.Request.Header.Get(authHeaderKey), " ")
		if len(tokens) < 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		if tokens[0] != authBearer {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		log.Println(tokens[1])
		user, err := useCase.RetrieveUserByToken(tokens[1])
		if err != nil {
			log.Panic(err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set("user", user)
		c.Next()
	}
}
