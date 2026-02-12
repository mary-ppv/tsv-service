package router

import (
	"github.com/gin-gonic/gin"

	admin "tsv-service/internal/transport/http/controllers/admin"
)

func NewRouter(
	unitsCtrl *admin.UnitsController,
	reportsCtrl *admin.ReportsController,
	importsCtrl *admin.ImportsController,
) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		adminAPI := api.Group("/admin")
		{
			adminAPI.GET("/units/:unit_guid/records", unitsCtrl.ListRecords)
			adminAPI.GET("/units/:unit_guid/reports", reportsCtrl.ListReports)
			adminAPI.POST("/imports/source-data", importsCtrl.ImportSourceData)
		}
	}

	return r
}
