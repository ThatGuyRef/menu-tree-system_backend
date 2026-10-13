package routes

import (
	"menu-tree-system/backend/config"
	"menu-tree-system/backend/handlers"
	"menu-tree-system/backend/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(cfg *config.Config, menuHandler *handlers.MenuHandler) *gin.Engine {
	r := gin.Default()

	// Pasang CORS middleware
	r.Use(middleware.CORSMiddleware(cfg.CORSOrigins))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Grouping API Menus
	api := r.Group("/api")
	{
		menus := api.Group("/menus")
		{
			menus.GET("", menuHandler.GetMenus)
			menus.GET("/:id", menuHandler.GetMenu)
			menus.POST("", menuHandler.CreateMenu)
			menus.PUT("/:id", menuHandler.UpdateMenu)
			menus.DELETE("/:id", menuHandler.DeleteMenu)
			menus.PATCH("/:id/move", menuHandler.MoveMenu)
			menus.PATCH("/:id/reorder", menuHandler.ReorderMenu)
		}
	}

	return r
}