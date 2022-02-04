package main

import (
	"context"
	"fmt"
	"os"

	pb "github.com/omustardo/wikitree-api-client/go/proto"
	"github.com/omustardo/wikitree-api-client/go/wikiclient"
	"google.golang.org/protobuf/proto"
)

func main() {
	fmt.Println("Starting WikiTree client")

	client, err := wikiclient.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req := &pb.GetProfileRequest{
		Key: proto.String("Kennedy-21529"),
	}
	resp, err := client.GetProfile(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("== Proto Response:")
	fmt.Println(resp)
}
