package main

import (
	"fmt"
	"url_shortner/routers"
	"url_shortner/store"
)

func main() {

	db := store.Newrdb()

	r := routers.SetupRouter(db)

	fmt.Println("URL SHORTNER running on http://localhost:8080")

	if db.IsHealthy() {
		fmt.Println("Redis :: UP 💚💚💚")
	} else {
		fmt.Println("Redis :: DOWN 🔴🔴🔴")
	}

	r.Run(":8080")
}
