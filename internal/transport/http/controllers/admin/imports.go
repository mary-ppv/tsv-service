package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tsv-service/internal/services/imports"
)

type ImportsController struct {
	svc *imports.Service
}

func NewImportsController(svc *imports.Service) *ImportsController {
	return &ImportsController{svc: svc}
}

func (c *ImportsController) ImportSourceData(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{"message": "todo"})
}
