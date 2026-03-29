package main

import (
	"fmt"
	"os"

	"community-elderly-care-platform/cmd"
)

// @title 霍营街道社区养老服务平台 API
// @version 1.0
// @description 昌平区霍营街道社区养老信息分发平台后端 API
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
