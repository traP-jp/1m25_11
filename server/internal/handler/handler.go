package handler

import (
	"github.com/traP-jp/1m25_11/server/internal/repository"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo *repository.Repository

}

func New(repo *repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}

}

func (h *Handler) SetupRoutes(api *echo.Group) {
	api.Use(h.AuthMiddleware)
	{
		pingAPI := api.Group("/ping")
		{
			pingAPI.GET("", h.Ping)
		}

		stampAPI := api.Group("/stamps")
		{
			stampAPI.GET("/search", h.getSearch)
			stampAPI.GET("/ranking", h.getRanking)
			stampAPI.GET("", h.getStamps)
			stampAPI.GET("/:stampId", h.getDetails)
			stampAPI.POST("/:stampId/tags/:tagId", h.createStampTags)
			stampAPI.DELETE("/:stampId/tags/:tagId", h.deleteStampTags)
			stampAPI.GET("/:stampId/descriptions", h.getDescriptions)
			stampAPI.POST("/:stampId/descriptions", h.createDescriptions)
			stampAPI.PUT("/:stampId/descriptions", h.updateDescriptions)
			stampAPI.DELETE("/:stampId/descriptions", h.deleteDescriptions)
		}

		tagAPI := api.Group("/tags")
		{
			tagAPI.GET("", h.getTags)
			tagAPI.POST("", h.createTags)
			tagAPI.GET("/:tagId", h.getCertainStamps)
			tagAPI.PUT("/:tagId", h.updateTags)
			tagAPI.DELETE("/:tagId", h.deleteTags)
			tagAPI.GET("/:tagId/stamps", h.getCertainStamps)
		}
		creatorAPI := api.Group("/me")
		{
			creatorAPI.GET("", h.getCreatorDetails)
		}
		userAPI := api.Group("/users-list")
		{
			userAPI.GET("", h.getUsersList)
		}
		loginAPI := api.Group("/login")
		{
			loginAPI.GET("", h.login)
		}
		callbackAPI := api.Group("/callback")
		{
			callbackAPI.GET("", h.callback)
		}
	}


}
