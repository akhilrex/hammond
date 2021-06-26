package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akhilrex/hammond/controllers"
	"github.com/akhilrex/hammond/db"
	"github.com/akhilrex/hammond/service"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var err error
	db.DB, err = db.Init()
	if err != nil {
		fmt.Println("status: ", err)
	} else {
		db.Migrate()
	}
	r := gin.Default()

	r.Use(setupSettings())
	r.Use(gin.Recovery())
	r.Use(location.Default())
	r.Use(static.Serve("/", static.LocalFile("./dist", true)))
	r.NoRoute(func(c *gin.Context) {
		//fmt.'Println(c.Request.URL.Path)
		c.File("dist/index.html")
	})
	router := r.Group("/api")

	dataPath := os.Getenv("DATA")

	router.Static("/assets/", dataPath)

	controllers.RegisterAnonController(router)
	controllers.RegisterAnonMasterConroller(router)
	controllers.RegisterSetupController(router)

	router.Use(controllers.AuthMiddleware(true))
	controllers.RegisterUserController(router)
	controllers.RegisterMastersController(router)
	controllers.RegisterAuthController(router)
	controllers.RegisterVehicleController(router)
	controllers.RegisterFilesController(router)
	controllers.RegisteImportController(router)

	go assetEnv()
	go intiCron()

	r.Run(":3000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
func setupSettings() gin.HandlerFunc {
	return func(c *gin.Context) {

		setting := db.GetOrCreateSetting()
		c.Set("setting", setting)
		c.Writer.Header().Set("X-Clacks-Overhead", "GNU Terry Pratchett")

		c.Next()
	}
}

func intiCron() {

	//gocron.Every(uint64(checkFrequency)).Minutes().Do(service.DownloadMissingEpisodes)
	gocron.Every(2).Days().Do(service.CreateBackup)
	<-gocron.Start()
}

func assetEnv() {
	log.Println("Config Dir: ", os.Getenv("CONFIG"))
	log.Println("Assets Dir: ", os.Getenv("DATA"))
}
