package main

import(
	"log"
	"./route"
)

func main() {
	r := route.RegisterRoute()
	err := route.DeployRoute(r, 9000)

	if err != nil {
		log.Fatal(err)
	}
}
