package main

import (
	"cine-api/config"
	"cine-api/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("Archivo .env no encontrado, usando variables del sistema")
	}

	// Conectar a MongoDB
	config.ConnectDB()

	// Configurar Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "servicio": "CineApp API"})
	})

	// Registrar rutas
	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("CineApp API iniciada en http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Error iniciando servidor:", err)
	}
}
