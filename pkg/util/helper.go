package util

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/structpb"
)

func MapListToInterfaceList(m []map[string]interface{}) []interface{} {
	ret := make([]interface{}, len(m))
	for i, item := range m {
		ret[i] = item
	}
	return ret
}
func MapListToStructList(m []map[string]interface{}) ([]*structpb.Struct, error) {
	ret := make([]*structpb.Struct, len(m))
	iList := MapListToInterfaceList(m)
	vList, err := structpb.NewList(iList)
	if err != nil {
		return nil, err
	}
	for i, item := range vList.GetValues() {
		ret[i] = item.GetStructValue()
	}
	return ret, nil
}

func MapToStruct(m map[string]interface{}) (*structpb.Struct, error) {
	if m == nil {
		return nil, nil
	}
	s, err := structpb.NewStruct(m)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func ContextIsBroken(ctx context.Context) (codes.Code, bool) {
	if ctx.Err() != nil {
		fmt.Println("broken:", ctx.Err().Error())
		switch ctx.Err() {
		case context.Canceled:
			return codes.Canceled, true
		case context.DeadlineExceeded:
			return codes.DeadlineExceeded, true
		default:
			return codes.Unavailable, true
		}
	}
	return codes.OK, false
}
