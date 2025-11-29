package config

import (
	"github.com/gin-contrib/cors"
)

func GetCORSConfig() cors.Config {
	return cors.Config{
		AllowOrigins:     []string{"https://monkreflections.com", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
}
