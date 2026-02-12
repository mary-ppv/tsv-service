package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"tsv-service/internal/transport/http/controllers/admin"
)

func NewRouter(
	db *sql.DB,
	unitsCtrl *admin.UnitsController,
	reportsCtrl *admin.ReportsController,
	importsCtrl *admin.ImportsController,
) *gin.Engine {

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(DBExecutor(db))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		adminAPI := api.Group("/admin")
		{
			adminAPI.GET(
				"/units/:unit_guid/records",
				unitsCtrl.ListRecords,
			)

			adminAPI.GET(
				"/units/:unit_guid/reports",
				reportsCtrl.ListReports,
			)

			adminAPI.POST(
				"/imports/source-data",
				importsCtrl.ImportSourceData,
			)
		}
	}

	return r
}
