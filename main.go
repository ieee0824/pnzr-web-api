package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ieee0824/pnzr-web/lib/config"
	"github.com/ieee0824/pnzr-web/api/deploy"
)

func main() {
	cfg := config.New()
	router := gin.Default()

	router.POST("/deploy", deploy.Deploy)

	if err := router.Run(cfg.Port); err != nil {
		panic(err)
	}
}