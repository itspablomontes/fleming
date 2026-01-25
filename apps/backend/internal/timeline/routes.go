package timeline

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	timeline := rg.Group("/timeline")
	{
		timeline.GET("", h.HandleGetTimeline)
		timeline.GET("/graph", h.HandleGetGraphData)

		timeline.GET("/events/:id", h.HandleGetEvent)
		timeline.POST("/events", h.HandleAddEvent)
		timeline.POST("/events/:id/correction", h.HandleCorrectEvent)
		timeline.DELETE("/events/:id", h.HandleDeleteEvent)

		timeline.POST("/events/:id/link", h.HandleLinkEvents)
		timeline.GET("/events/:id/related", h.HandleGetRelatedEvents)
		timeline.DELETE("/edges/:edgeId", h.HandleUnlinkEvents)

		timeline.GET("/events/:id/files/:fileId", h.HandleDownloadFile)
		timeline.GET("/events/:id/files/:fileId/key", h.HandleGetFileKey)
		timeline.POST("/events/:id/files/:fileId/share", h.HandleShareFile)

		timeline.POST("/events/:id/files/multipart/start", h.HandleStartMultipartUpload)
		timeline.PUT("/events/:id/files/multipart/part", h.HandleUploadMultipartPart)
		timeline.POST("/events/:id/files/multipart/complete", h.HandleCompleteMultipartUpload)
	}
}
