package main

import (
	"fmt"
	"kings-house-back/API/config"
	"kings-house-back/API/database"
	"kings-house-back/API/handlers"
	"os"

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

	secretValue := os.Getenv("JWT_SECRET")
	log.Println("DEBUG JWT_SECRET:", secretValue)

	var secret = []byte(secretValue) // ahora sí asignas el contenido real

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

	router := gin.Default()
	router.Use(cors.New(corsConfig)) //

	//Crear
	router.POST("/auth/login", handlers.LoginHandler(db))
	router.POST("/users", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador", "gestor"), handlers.CrearUsuarioHandler(db))              // Creacion de un Usuario                                                                                                   // Endpoint para el Login
	router.POST("/pedidos", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador", "gestor", "vendedor"), handlers.CrearPedidoHandler(db)) // Creacion de pedido

	//Leer
	router.GET("/users", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador"), handlers.ListarUsuariosHandler(db))                                                       // Lectura de todos los Usuarios
	router.GET("/users/:id", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador", "gestor"), handlers.ObtenerUsuarioHandler(db))                                         // Lectura de un Usuario por ID
	router.GET("/pedidos", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador", "gestor"), handlers.ListarPedidosHandler(db))                                            // Lectura de todos los pedidos
	router.GET("/pedidos/:id", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador", "gestor", "vendedor"), handlers.ObtenerPedidoHandler(db))                            // Lectura por el ID de un pedido
	router.GET("/usuarios/:usuario_id/pedidos", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador", "gestor", "vendedor"), handlers.ListarPedidosPorUsuarioHandler(db)) // Lectura de todos los pedidos de un vendedor
	router.GET("/users/vendedores", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador", "gestor"), handlers.ListarVendedoresHandler(db))                                // Lectura de todos los vendedores

	//Actualizar
	router.PUT("/users/:id", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador"), handlers.ActualizarUsuarioHandler(db))            // Modificar un Usuario por ID
	router.PUT("/pedidos/:id", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador", "gestor"), handlers.ActualizarPedidoHandler(db)) // Modificar un Pedido por ID

	//Eliminar
	router.DELETE("/users/:id", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador"), handlers.EliminarUsuarioHandler(db))            // Eliminar un Usuario por ID
	router.DELETE("/pedidos/:id", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador", "gestor"), handlers.EliminarPedidoHandler(db)) //Eliminar un Pedido por ID

	router.Run(":8080")
}
