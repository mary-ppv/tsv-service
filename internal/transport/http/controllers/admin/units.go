package admin

import (
	"net/http"
	"strconv"
	"tsv-service/internal/transport/http/controllers"

	"github.com/gin-gonic/gin"
)

type UnitsController struct {
	controllers.Base
}

func NewUnitsController(base controllers.Base) *UnitsController {
	return &UnitsController{base}
}

// GET /api/admin/units/:unit_guid/records
func (c *UnitsController) ListRecords(ctx *gin.Context) {
	svc, err := c.ServiceProvider().UnitsService()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	unit := ctx.Param("unit_guid")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))

	rows, total, err := svc.ListRecords(ctx.Request.Context(), unit, page, limit)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error(), "message": "can't get records"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": rows,
		"meta": gin.H{"total": total, "page": page, "limit": limit},
	})
}
