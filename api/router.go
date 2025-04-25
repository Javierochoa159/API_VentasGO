package api

import (
	"net/http"
	"parte3/internal/sale"
	"parte3/internal/user"

	"github.com/gin-gonic/gin"
)

// InitRoutes registers all user CRUD endpoints on the given Gin engine.
// It initializes the storage, service, and handler, then binds each HTTP
// method and path to the appropriate handler function.
func InitRoutes(e *gin.Engine) {
	userStorage := user.NewLocalStorage()
	userService := user.NewService(userStorage)
	saleStorage := sale.NewLocalStorage()
	saleService := sale.NewService(saleStorage)

	h := handler{
		userService: userService,
		saleService: saleService,
	}

	e.POST("/users", h.handleCreateUser)
	e.GET("/users/:id", h.handleReadUser)
	e.PATCH("/users/:id", h.handleUpdateUser)
	e.DELETE("/users/:id", h.handleDeleteUser)

	e.POST("/sales", h.handleCreateSale)
	e.GET("/sales/:id", h.handleReadSale)
	e.PATCH("/sales/:id", h.handleUpdateSale)
	e.DELETE("/sales/:id", h.handleDeleteSale)

	e.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
