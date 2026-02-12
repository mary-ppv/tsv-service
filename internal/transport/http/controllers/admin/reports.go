package admin

import (
	"net/http"
	"strconv"

	"tsv-service/internal/transport/http/controllers"

	"github.com/gin-gonic/gin"
)

type ReportsController struct {
	controllers.Base
}

func NewReportsController(base controllers.Base) *ReportsController {
	return &ReportsController{base}
}

// GET /api/admin/units/:unit_guid/reports?limit=10
func (c *ReportsController) ListReports(ctx *gin.Context) {

	svc, err := c.ServiceProvider().ReportsService()
	if err != nil {
		c.Logger().Error("reports service not available", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	unit := ctx.Param("unit_guid")

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))

	list, err := svc.ListReports(ctx.Request.Context(), unit, limit)
	if err != nil {
		c.Logger().Error("list reports failed", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}
