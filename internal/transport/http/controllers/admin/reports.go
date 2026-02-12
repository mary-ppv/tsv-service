package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tsv-service/internal/services/reports"
)

type ReportsController struct {
	svc *reports.Service
}

func NewReportsController(svc *reports.Service) *ReportsController {
	return &ReportsController{svc: svc}
}

func (c *ReportsController) ListReports(ctx *gin.Context) {
	unit := ctx.Param("unit_guid")

	limit := 0
	if v := ctx.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}

	reps, err := c.svc.ListReports(ctx, unit, limit)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error(), "message": "can't get reports"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": reps})
}
