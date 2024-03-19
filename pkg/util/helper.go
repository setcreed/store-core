package util

import "google.golang.org/protobuf/types/known/structpb"

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
