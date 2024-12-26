package router

// import (
// 	"go_silo/db"
// 	"go_silo/services"

// 	"github.com/gin-gonic/gin"
// )

// func Router() *gin.Engine {
// 	r := gin.Default()
// 	// - Material
// 	// 获取 material 表中数据
// 	r.GET("/materials", services.GetMaterial(db.Mongo))
// 	// 修改 material 表中数据
// 	r.PUT("/materials", services.UpdateMaterial(db.Mongo))
// 	// 新增 material 表中数据
// 	r.POST("/materials", services.AddMaterial(db.Mongo))
// 	// 删除 material 表中数据
// 	r.DELETE("/materials", services.DeleteMaterial(db.Mongo))

// 	r.GET("/statics", services.GetStatic(db.Mongo))
// 	r.PUT("/statics", services.UpdateStatic(db.Mongo))
// 	return r
// }

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(db *mongo.Database) *gin.Engine {
	r := gin.Default()

	// Material 路由
	MaterialRouter(r, db)

	// Static 路由
	StaticRouter(r, db)

	return r
}
