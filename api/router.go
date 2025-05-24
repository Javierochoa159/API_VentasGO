package api

import (
	"API_VentasGO/internal/metadata"
	"API_VentasGO/internal/sale"
	"API_VentasGO/internal/user"

	"github.com/gin-gonic/gin"
)

// InitRoutes registers all user CRUD endpoints on the given Gin engine.
// It initializes the storage, service, and handler, then binds each HTTP
// method and path to the appropriate handler function.
func InitRoutes(e *gin.Engine) {
	userStorage := user.NewLocalStorage()
	userService := user.NewService(userStorage, nil)
	saleStorage := sale.NewLocalStorage()
	saleService := sale.NewService(saleStorage, nil, nil)
	metadataStorage := metadata.NewLocalStorage()
	metadataService := metadata.NewService(metadataStorage)

	h := handler{
		userService:     userService,
		saleService:     saleService,
		metadataService: metadataService,
	}

	e.POST("/users", h.handleCreateUser)
	e.GET("/users/:id", h.handleReadUser)
	e.PATCH("/users/:id", h.handleUpdateUser)
	e.DELETE("/users/:id", h.handleDeleteUser)

	e.POST("/sales", h.handleCreateSale)
	e.GET("/sales", h.handleReadSale)
	e.PATCH("/sales/:id", h.handleUpdateSale)

}
