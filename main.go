package main

import (
	"context"
	"fmt"

	"github.com/princedraculla/go-microservice/application"
)

func main() {

	app := application.New()

	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("faild to start the server...! \n", err)
	}

}
