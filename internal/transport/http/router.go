package http

import "github.com/gin-gonic/gin"

func Router(h Handler) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1/swift-codes")
	{
		v1.GET("/:code", h.GetByCode)
		v1.GET("/country/:iso2", h.Country)
		v1.POST("", h.Create)
		v1.DELETE("/:code", h.Delete)
	}
	return r
}
