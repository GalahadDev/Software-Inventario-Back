package main

import (
	"fmt"
	"kings-house-back/API/config"
	"kings-house-back/API/database"
	"kings-house-back/API/handlers"

	//"kings-house-back/API/models"
	"log"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// Conectar a la base de datos
	db, err := database.OpenGormDB()
	if err != nil {
		log.Fatalf("Error al conectarse a la BD: %v", err)
	}

	//db.AutoMigrate(&models.Usuario{}, &models.Pedido{})

	fmt.Print(config.DBURL())

	// Configurar CORS
	corsConfig := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	fmt.Printf("db: %v\n", db)

	router := gin.Default()
	router.Use(cors.New(corsConfig))

	//Crear
	router.POST("/users", handlers.CrearUsuarioHandler(db)) // Creacion de un Usuario
	router.POST("/auth/login", handlers.LoginHandler(db))   // Endpoint para el Login

	//Leer
	router.GET("/users", handlers.ListarUsuariosHandler(db))     // Lectura de todos los Usuarios
	router.GET("/users/:id", handlers.ObtenerUsuarioHandler(db)) // Lectura de un Usuario por ID

	//Actualizar
	router.PUT("/users/:id", handlers.ActualizarUsuarioHandler(db)) // Modificar un Uusario por ID

	//Eliminar
	router.DELETE("/users/:id", handlers.EliminarUsuarioHandler(db)) // Eliminar un Usuario por ID

	router.Run(":8080")
}
