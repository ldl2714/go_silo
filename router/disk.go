package router

import (
	"go_silo/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func DiskRouter(r *gin.Engine, db *mongo.Database) {
	r.GET("/disks", services.GetDisk(db))
	r.PUT("/disks", services.PutDisk(db))
}
