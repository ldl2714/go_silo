package services

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetProduction(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
