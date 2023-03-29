package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	api "web/hello"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client api.HelloServiceClient

func init() {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Did not connect: %v\n", err)
	}
	client = api.NewHelloServiceClient(conn)
	fmt.Println("init grpc client")
}

func main() {

	route := gin.Default()

	route.GET("/:name", func(c *gin.Context) {
		name := c.Param("name")

		resp, err := client.SayHello(
			context.Background(),
			&api.HelloReq{Name: name},
		)
		if err != nil {
			fmt.Printf("rpc err %+v\n", err)
		}

		fmt.Println(resp.GetResult())

		c.JSON(200, gin.H{"resp": resp.GetResult()})
	})

	apiSrv := &http.Server{
		Addr:    ":50052",
		Handler: route,
	}
	if err := apiSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

}
