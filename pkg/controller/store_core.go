package controller

import (
	"context"
	"fmt"
	"io"

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
	fmt.Printf("sqlApi  name:%v, table:%v, sql:%v\n", sqlApi.Name, sqlApi.Table, sqlApi.Sql)
	ret, err := db.f.Store().Query(ctx, sqlApi, in.Params)
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

func (db *DBService) Exec(ctx context.Context, in *v1.ExecRequest) (*v1.ExecResponse, error) {
	sqlApi := db.cfg.FindSQLByName(in.Name)
	if sqlApi == nil {
		return nil, status.Error(codes.Unavailable, "error api name")
	}
	rowsAffected, selectKey, err := db.f.Store().ExecBySql(ctx, sqlApi, in.Params)
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}
	selectKeyValue, err := util.MapToStruct(selectKey)
	if err != nil {
		selectKeyValue = nil
	}
	return &v1.ExecResponse{
		Message:      "success",
		RowsAffected: rowsAffected,
		Select:       selectKeyValue,
	}, nil
}

func (db *DBService) Tx(c v1.DBService_TxServer) error {
	tx := db.f.Store().GetDB().Begin()
	for {
		txRequest, err := c.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF")
				tx.Commit()
				return nil
			}
			fmt.Println("Error-0:", err)
			tx.Rollback()
			return err
		}
		sqlApi := db.cfg.FindSQLByName(txRequest.Api)
		if sqlApi == nil {
			tx.Rollback()
			return status.Error(codes.Unavailable, "error api name")
		}

		ret := make(map[string]interface{})
		if txRequest.Type == "query" {
			queryRet, err := db.f.Store().QueryBySql(context.Background(), sqlApi, txRequest.Params)
			if err != nil {
				fmt.Println("Error-1:", err)
				tx.Rollback()
				return err
			}
			ret["query"] = queryRet
		} else {
			af, sk, err := db.f.Store().ExecBySql(context.Background(), sqlApi, txRequest.Params)
			if err != nil {
				fmt.Println("Error-2:", err)
				tx.Rollback()
				return err
			}
			ret["exec"] = []interface{}{af, sk} //受影响的行、selectKey
		}
		m, _ := util.MapToStruct(ret)
		err = c.Send(&v1.TxResponse{
			Result: m,
		})
		if err != nil {
			fmt.Println("Error-3:", err)
			tx.Rollback()
			return err
		}
	}
}
