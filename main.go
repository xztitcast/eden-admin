package main

import (
	"eden-admin/app/admin/router"
)

func main() {
	r := router.InitRouter()
	r.Run(":8080")
}
