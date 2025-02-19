package main

import (
	"fmt"
	"kings-house-back/API/config"
	"kings-house-back/API/database"
	"kings-house-back/API/handlers"
	"kings-house-back/API/ws"

	"os"


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

	hub := ws.NewHub()
	
	var secret = []byte(secretValue)

	fmt.Print(config.DBURL())

	// Configurar CORS
	/*corsConfig := cors.Config{
   		AllowOrigins:     []string{"https://kings-bed-sm.onrender.com"},
   		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    		ExposeHeaders:    []string{"Content-Length"},
    		AllowCredentials: true,
    		MaxAge:           12 * time.Hour,
	}*/

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
	router.Use(cors.New(corsConfig))

	//Crear
	router.POST("/auth/login", handlers.LoginHandler(db)) 													       // Endpoint para el Login
	router.POST("/users", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador", "gestor"), handlers.CrearUsuarioHandler(db))                   // Creacion de un Usuario                                                                                                  
	router.POST("/pedidos", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador", "gestor", "vendedor"), handlers.CrearPedidoHandler(db, hub)) // Creacion de pedido

	//Leer
	router.GET("/users", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador"), handlers.ListarUsuariosHandler(db))                                                       // Lectura de todos los Usuarios
	router.GET("/users/:id", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador", "gestor"), handlers.ObtenerUsuarioHandler(db))                                         // Lectura de un Usuario por ID
	router.GET("/pedidos", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador", "gestor"), handlers.ListarPedidosHandler(db))                                            // Lectura de todos los pedidos
	router.GET("/pedidos/:id", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador", "gestor", "vendedor"), handlers.ObtenerPedidoHandler(db))                            // Lectura por el ID de un pedido
	router.GET("/usuarios/:usuario_id/pedidos", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador", "gestor", "vendedor"), handlers.ListarPedidosPorUsuarioHandler(db)) // Lectura de todos los pedidos de un vendedor
	router.GET("/users/vendedores", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador", "gestor"), handlers.ListarVendedoresHandler(db))                                // Lectura de todos los vendedores
	router.GET("/reportes/vendedores/:vendedor_id/montos", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador", "gestor"), handlers.SumarMontosPorVendedor(db))          //Lectura de comision de un vendedor

	//Actualizar
	router.PUT("/users/:id", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador"), handlers.ActualizarUsuarioHandler(db))                 // Modificar un Usuario por ID
	router.PUT("/pedidos/:id", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador", "gestor"), handlers.ActualizarPedidoHandler(db, hub)) // Modificar un Pedido por ID
	router.PUT("/users/bank-data", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("vendedor"), handlers.ActualizarDatosBancariosHandler(db))         // Modificar datos bancarios desde un vendedor

	//Eliminar
	router.DELETE("/users/:id", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador"), handlers.EliminarUsuarioHandler(db))            // Eliminar un Usuario por ID
	router.DELETE("/pedidos/:id", handlers.AuthMiddleware(secret), handlers.RoleMiddleware("administrador", "gestor"), handlers.EliminarPedidoHandler(db)) // Eliminar un Pedido por ID

	//WebSocket
	router.GET("/ws", handlers.AuthMiddlewareQuery(secret), handlers.RoleMiddleware("administrador", "gestor"), handlers.WSHandler(hub)) // Conectar a WebSocket

	router.Run(":8080")
}
