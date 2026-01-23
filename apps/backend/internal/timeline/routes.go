package timeline

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	timeline := rg.Group("/timeline")
	{
		timeline.GET("", h.HandleGetTimeline)
		timeline.GET("/graph", h.HandleGetGraphData)

		timeline.GET("/events/:id", h.HandleGetEvent)
		timeline.POST("/events", h.HandleAddEvent)
		timeline.DELETE("/events/:id", h.HandleDeleteEvent)

		timeline.POST("/events/:id/link", h.HandleLinkEvents)
		timeline.GET("/events/:id/related", h.HandleGetRelatedEvents)
		timeline.DELETE("/edges/:edgeId", h.HandleUnlinkEvents)
	}
}
