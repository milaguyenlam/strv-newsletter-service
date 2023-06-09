package strvnewsletterservice

import (
	"log"

	"github.com/gin-gonic/gin"
	"strv.com/newsletter/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := gin.Default()

}
