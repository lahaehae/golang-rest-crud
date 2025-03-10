package main

import (
	"context"
	"fmt"
	"log"

	//"os/user"
	"sync"
	"time"

	pb "github.com/lahaehae/crud_project/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)



func main() {


	conn, err := grpc.NewClient(":9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server")
	}
	defer conn.Close()

	c := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(),10 * time.Second)
	defer cancel()

	var wg sync.WaitGroup
	numRequests := 1000


	for i:= 0; i < numRequests; i++ {
		wg.Add(1)
		go func(id int){
			defer wg.Done()
			userCreateResp, err := c.CreateUser(ctx, &pb.CreateUserRequest{
				Name: fmt.Sprintf("User%d", id),
				Email: fmt.Sprintf("user%d@mail.com", id),
			})
			if err != nil{
				log.Printf("Failed to create user %d: %v", id, err)
				return
			}
			log.Printf("Created user: %v", userCreateResp)
		}(i)
	}

	wg.Wait()
	log.Println("All requests completed")
}
