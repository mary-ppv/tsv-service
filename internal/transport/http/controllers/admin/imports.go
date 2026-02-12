package admin

import (
	"net/http"
	"tsv-service/internal/transport/http/controllers"

	"github.com/gin-gonic/gin"
)

type ImportsController struct {
	controllers.Base
}

func NewImportsController(base controllers.Base) *ImportsController {
	return &ImportsController{
		Base: base,
	}
}

func (c *ImportsController) ImportSourceData(ctx *gin.Context) {
	_, err := c.ServiceProvider().ImportsService()
	if err != nil {
		c.Logger().Error("imports service missing", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//todo smth

	ctx.JSON(http.StatusOK, gin.H{
		"status": "import started",
	})
}
