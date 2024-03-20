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

func MakeParam(m map[string]interface{}) *v1.SimpleParams {
	paramStruct, _ := structpb.NewStruct(m)
	return &v1.SimpleParams{
		Params: paramStruct,
	}
}

func main() {
	client, err := grpc.DialContext(context.Background(),
		"localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	dbServiceClient := v1.NewDBServiceClient(client)

	// 执行事务
	tx, err := dbServiceClient.Tx(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	addUserParam := MakeParam(map[string]interface{}{
		"user_name":     "qqqq",
		"user_password": "123456",
	})
	err = tx.Send(&v1.TxRequest{
		Api:    "addUser",
		Params: addUserParam,
		Type:   "exec",
	})
	if err != nil {
		log.Fatal(err)
	}
	addUserRsp, err := tx.Recv()
	if err != nil {
		log.Fatal(err)
	}
	ret := addUserRsp.Result.AsMap()
	uid := ret["exec"].([]interface{})[1].(map[string]interface{})["user_id"]
	fmt.Println("用户ID是", uid)

	log.Fatal("意外退出")

	addOrderParam := MakeParam(map[string]interface{}{
		"user_id": uid, // 下订单的用户ID
	})
	err = tx.Send(&v1.TxRequest{Api: "addUserOrder", Params: addOrderParam, Type: "exec"})
	if err != nil {
		log.Fatal(err)
	}
	addOrderRsp, err := tx.Recv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(addOrderRsp.Result.AsMap())
	err = tx.CloseSend()
	fmt.Println("结束")
}
