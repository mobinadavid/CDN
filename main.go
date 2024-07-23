package main

import (
	"cdn/cmd"
	"log"
)

// @title CDN
// @version 1.0
// @description This is a swagger server for documentation of cdn app.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host cdn.omaxplatform.com
// @BasePath /app/api/v1
func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
