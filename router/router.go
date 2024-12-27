package router

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(db *mongo.Database) *gin.Engine {
	r := gin.Default()
	//Belt-皮带
	BeltRouter(r, db)
	// Material-物料
	MaterialRouter(r, db)
	// Static-静态信息
	StaticRouter(r, db)

	return r
}
