package main

import (
	"context"
	"fmt"
	"log"

	spec "article-openapi"
)

const server = "http://localhost:8088"

func main() {
	client, err := spec.NewClientWithResponses(server)
	if err != nil {
		log.Fatal(err)
	}

	item := spec.OrderItemTeaTableGreen
	price := 14
	resp, err := client.PutOrderIdWithResponse(context.Background(), "234578", spec.PutOrderIdJSONRequestBody{
		Item: &item, Price: &price,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.StatusCode())
}
