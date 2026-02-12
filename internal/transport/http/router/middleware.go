package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"tsv-service/internal/repository/driver"
)

func DBExecutor(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := driver.ExecutorToContext(c.Request.Context(), db)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
