package main

import (
	"context"
	"fmt"
	"flag"
	"google.golang.org/grpc"

	pb "tokenmngr/proto"

)

func main() {

	
	serv_addr := flag.String("host", "localhost", "Server address")
	serv_port := flag.Int("port", 50051, "Server port")
	
	id := flag.String("id", "", "ID")
	name := flag.String("name", "", "Name")
	low := flag.Uint64("low", 0, "Lower-bound")
	mid := flag.Uint64("mid", 0, "Midpoint")
	high := flag.Uint64("high", 0, "Upper-bound")

	create := flag.Bool("create", false, "Create token")
	read := flag.Bool("read", false, "Read token")
	drop := flag.Bool("drop", false, "Drop token")
	write := flag.Bool("write", false, "Write token")

	flag.Parse() // Parse command line flags

	
	
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *serv_addr, *serv_port), grpc.WithInsecure())

	if err != nil {
		fmt.Println("Connection failed: %v", err)
		//log.Fatalf("Failed to connect to server: %v", err)
	}

	defer conn.Close()

	client := pb.NewTokenManagerClient(conn) // New Client
	
	ctx := context.Background() // Context

	
	switch {
	//Switch case to select operation
	case *create:

		output, err := client.Create(ctx, &pb.CreateRequest{Id: *id})
		
		if err != nil {
			fmt.Println("Error in Creation: %v", err)
		}

		fmt.Println(output.Success)

	case *drop:
		output, err := client.Drop(ctx, &pb.DropRequest{Id: *id})
		if err != nil {
			fmt.Println("Error in deleting: %v", err)
		}
		fmt.Println(output.Success)

	case *write:
		output, err := client.Write(ctx, &pb.WriteRequest{
			
			Id:    *id,
			Name:  *name,
			Low:   *low,
			Mid:   *mid,
			High:  *high,
		})

		if err != nil {
			fmt.Println("Error in writing: %v", err)
		}
		fmt.Println("Partial: %d", output.Partial)

	case *read:
		output, err := client.Read(ctx, &pb.ReadRequest{Id: *id})
		if err != nil {
			fmt.Println("Error in reading: %v", err)
		}
		fmt.Println("Final: %d", output.Final)
	

	default:
		fmt.Println("Operation not given")
	}
}
