package main

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type File struct {
	ID          uint
	Filename    string
	Signature   string
	TenantID    string
	DateCreated string
	DateDeleted string
}

func main() {
	// Force log's color
	gin.ForceConsoleColor()

	router := gin.Default()

	router.GET("/dummyland/image", func(c *gin.Context) {
		// get tenantID from query parameter
		tenantID := c.Query("tenantID")

		db, err := sqlx.Open("sqlite3", "./dummy.db")
		defer db.Close()

		if err != nil {
			slog.Error("Failed to open database", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open database"})
			return
		}

		rows, err := db.Queryx("SELECT * FROM Files WHERE tenant_id = '" + tenantID + "'")
		if err != nil {
			slog.Error("Failed to query database", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
			return
		}
		defer rows.Close()

		var files []File
		for rows.Next() {
			var f File
			err := rows.StructScan(&f)
			if err != nil {
				slog.Error("Failed to scan database row", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan database row"})
				return
			}
			files = append(files, f)
		}

		// return files as an object in the JSON response
		c.JSON(http.StatusOK, gin.H{"files": files})
	})

	router.Run() // listen and serve on 0.0.0.0:8080
}
