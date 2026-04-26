package main

import (
	"fmt"
	"url_shortner/routers"
	"url_shortner/store"
)

func main() {

	db := store.NewDB()

	r := routers.SetupRouter(db)

	fmt.Println("URL SHORTNER running on http://localhost:8080")
	r.Run(":8080")
}
