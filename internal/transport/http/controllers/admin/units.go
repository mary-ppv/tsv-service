package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tsv-service/internal/services/units"
)

type UnitsController struct {
	svc *units.Service
}

func NewUnitsController(svc *units.Service) *UnitsController {
	return &UnitsController{svc: svc}
}

func (c *UnitsController) ListRecords(ctx *gin.Context) {
	unit := ctx.Param("unit_guid")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))

	rows, total, err := c.svc.ListRecords(ctx, unit, page, limit)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error(), "message": "can't get records"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": rows,
		"meta": gin.H{"total": total, "page": page, "limit": limit},
	})
}
