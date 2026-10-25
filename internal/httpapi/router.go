package httpapi

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(ctrl *Controller) *gin.Engine {
	r := gin.Default()

	r.GET("/:code", ctrl.Redirect)
	r.POST("/links", ctrl.MakeLink)
	r.GET("/:code/clicks", ctrl.GetClicks)
	r.DELETE("/:code", ctrl.DeleteLink)

	return r
}
