package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"

	v1 "github.com/setcreed/store-core/api/store_service/v1"
)

func main() {
	client, err := grpc.DialContext(context.Background(),
		"localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatal(err)
	}
	paramStruct, err := structpb.NewStruct(map[string]interface{}{
		//"user_id>?": 1,
		//"user_id<?": 3,
		//"id": 1,
		//"user_name":     "test",
		//"user_password": "test",
		"user_id": 4,
	})
	if err != nil {
		log.Fatal(err)
	}
	params := &v1.SimpleParams{
		Params: paramStruct,
	}
	//req := &v1.QueryRequest{Name: "userList", Params: params}
	req := &v1.QueryRequest{Name: "deleteUser", Params: params}
	rsp := &v1.QueryResponse{}
	err = client.Invoke(context.Background(),
		"/store_service.v1.DBService/Exec", req, rsp)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range rsp.Result {
		fmt.Println(item.AsMap())
	}
}
