package router

import (
	"go_silo/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProductionRouter(r *gin.Engine, db *mongo.Database) {
	r.GET("/productions", services.GetProduction(db))
}
