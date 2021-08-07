package main

import (
	"log"

	spec "article-openapi"

	"github.com/labstack/echo/v4"
)

const address = ":8088"

type server struct {
}

func (s server) PutOrderId(c echo.Context, id string) error {
	var req spec.PutOrderIdJSONRequestBody
	if err := c.Bind(&req); err != nil {
		log.Fatal(err)
	}
	log.Printf("id: %v, req: %v", id, req)
	return nil
}

func main() {
	e := echo.New()
	spec.RegisterHandlers(e, server{})

	e.Logger.Fatal(e.Start(address))
}
