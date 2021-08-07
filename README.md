# Go Code Generation from OpenAPI spec

One of the nicest features of Go is the power of code generation. `go generate` command serves as a Swish knife allowing you to generate enums, mocks and stubs. In this article, we will employ this feature to generate a Go code from OpenAPI specification. OpenAPI specification is a modern industrial standard for REST API. This standard has fantastic tooling support and allows you to conveniently render and validate the spec. We are going to befriend the power of Go code generation with the elegance and clarity of the OpenAPI specification. In this way, you don't have to manually update the Go boilerplate code after every change in the spec. You also ensure that your docs and your code are a single entity, as your code is being begotten from the docs.

Let's start dead-simple: we have a service that accepts order requests. Let's declare endpoint `order/10045234` that accepts PUT requests, where `10045234` is an ID of a particular order. We expect to receive an order as a JSON payload in the following format.

```json
    {"item":  "Tea Table Green", "price":  106}
```

How can describe this endpoint in the OpenAPI spec?

The skeleton of the spec looks the following:

```yaml
    openapi: 3.0.3
    info: 
        title: Go and OpenAPI Are Friends Forever
        description: Let's do it
        version: 1.0.0
    paths:
    components:
```

First, we need to describe the response body - `order` and put this description under section `components`. In OpenAPI the order object looks the following:

```yaml
components: 
  schemas: 
    Order:
      type: object
      properties: 
        item:
          type: string
        price:
          type: integer
```

It's a bit more verbose, but it is also very expressive. We can easily add validations rules on top of it. For example, we can enlist all possible items:

```yaml
item:
  type: string
  enum:
    - Tea Table Green
    - Tea Table Red
```

OpenAPI specification is very rich and lets one easily specify that price must be a positive integer or that id must be in UUID format. We also can specify which fields are mandatory.

Now, let's specify our endpoint under section `paths`.

```yaml
paths:
"/order/{id}":
    put:
      summary: Create an order
      parameters:
        - in: path
          description: Order ID
          name: id
          schema: 
            type: string
      requestBody:
        required: true
        content: 
          application/json:
            schema: 
              $ref: "#/components/schemas/Order"
```

The above part of the specification states that we expect a PUT request on the endpoint "/order/id" with the `Order` object as payload. Finally, let's add an expected response to the endpoint.

```yaml
responses: 
    "201":
      description: The order was successfully created.
```

## Code Generation

At this point we have the complete specification:

```yaml
openapi: 3.0.3
info:
  title: Title
  description: Title
  version: 1.0.0
components:
  schemas:
    Order:
      type: object
      properties:
        item:
          type: string
          enum:
            - Tea Table Green
            - Tea Table Red
        id:
          type: string
        price:
          type: integer
paths:
  "/order/{id}":
    put:
      summary: Create an order
      parameters:
        - in: path
          description: Order ID
          name: id
          schema: 
            format: string
      requestBody:
        required: true
        content: 
          application/json:
            schema: 
              $ref: "#/components/schemas/Order"
      responses: 
        "201":
          description: The order was successfully created.
```

Now, the fun part: let the machine do its job and generate the code. To do so we will employ an awesome Go lib [oapi-codegen](https://github.com/deepmap/oapi-codegen). Download it with `go get` and add the following command to your Makefile:

```makefile
gen:
    oapi-codegen -package spec spec.yaml > gen.go
```

Now, let's give it a try and see what we got from the generation. For the client side, we have `NewClientWithResponses` function that returns a ready-to-use client object. The generator also creates order object as `PutOrderIdJsonRequestBody` and enum values for `item` properties. You can notice that `price` and `item` expects pointers, as those fields are optional in our OpenAPI spec.

```go
const server = "localhost:8080"

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
```

In less than 15 lines of hand-written code, we got a full-fledged client of our API! The OpenAPI standard has a tremendous palette of tooling and the code be generated for a variety of languages.

And for the server side, we got interface `ServerInterface` with a single method `PutOrderId` to implement. Let's declare a type `server` that conforms to this interface. By default, the code is generated for router `echo`, which can be downloaded as `go get github.com/labstack/echo`.

```go
type server struct {}

func (s server) PutOrderId(c echo.Context, id string) error {
   var req spec.PutOrderIdJSONRequestBody
   if err := c.Bind(&req); err != nil {
      log.Fatal(err)
   }
   log.Printf("id: %v, req: %v", id, req)
   return nil
}
```

Now all we need to do is to register this handler:

```go
func main() {
   e := echo.New()
   spec.RegisterHandlers(e, server{})

   e.Logger.Fatal(e.Start(address))
}
```

At this stage, we have a fully working client and server with a minimum amount of hand-written code. There is a ton of cool features we can add on top of that, for example, `echo` offers a middleware that will automatically validate the request payload and return detailed validation error. Overall, the usage of OpenAPI specification in combination with `oapi-codegen` generator and `echo` framework can radically increase the speed of web development by removing the necessity for hand-written boiler-plate code. 
