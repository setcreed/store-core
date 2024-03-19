package controller

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1 "github.com/setcreed/store-core/api/store_service/v1"
	"github.com/setcreed/store-core/cmd/app/config"
	"github.com/setcreed/store-core/pkg/data"
	"github.com/setcreed/store-core/pkg/util"
)

type DBService struct {
	v1.UnimplementedDBServiceServer

	f   data.ShareDaoFactory
	cfg *config.Config
}

func NewDBService(cfg *config.Config, f data.ShareDaoFactory) *DBService {
	return &DBService{
		f:   f,
		cfg: cfg,
	}
}

func (db *DBService) Query(ctx context.Context, in *v1.QueryRequest) (*v1.QueryResponse, error) {
	sqlApi := db.cfg.FindSQLByName(in.Name)
	if sqlApi == nil {
		return nil, status.Error(codes.Unavailable, "error api name")
	}

	ret, err := db.f.Store().QueryByTableName(ctx, sqlApi.Table, in.Params)
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}
	// 把map 转化为 StructList
	list, err := util.MapListToStructList(ret)
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	return &v1.QueryResponse{
		Message: "success",
		Result:  list,
	}, nil
}
