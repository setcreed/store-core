package data

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	v1 "github.com/setcreed/store-core/api/store_service/v1"
	"github.com/setcreed/store-core/cmd/app/config"
)

type StoreInterface interface {
	QueryByTableName(ctx context.Context, sqlConfig *config.SQLConfig, params *v1.SimpleParams) ([]map[string]interface{}, error)
	QueryBySql(ctx context.Context, sqlConfig *config.SQLConfig, params *v1.SimpleParams) ([]map[string]interface{}, error)
	Query(ctx context.Context, sqlConfig *config.SQLConfig, params *v1.SimpleParams) ([]map[string]interface{}, error)
}

type store struct {
	db *gorm.DB
}

func newStore(db *gorm.DB) StoreInterface {
	return &store{
		db: db,
	}
}

func (s *store) QueryByTableName(ctx context.Context, sqlConfig *config.SQLConfig, params *v1.SimpleParams) ([]map[string]interface{}, error) {
	dbResult := make([]map[string]interface{}, 0)
	db := s.db.Table(sqlConfig.Table)
	for k, v := range params.Params.AsMap() {
		db = db.Where(k, v)
	}
	db = db.Find(&dbResult)

	return dbResult, db.Error
}

// Query 设置sql优先，如果写了sql和table，优先执行sql
func (s *store) Query(ctx context.Context, sqlConfig *config.SQLConfig, params *v1.SimpleParams) ([]map[string]interface{}, error) {
	if sqlConfig.Sql == "" || sqlConfig.Table == "" {
		return nil, fmt.Errorf("error sql or table ")
	}
	if sqlConfig.Sql != "" {
		return s.QueryBySql(ctx, sqlConfig, params)
	}

	return s.QueryByTableName(ctx, sqlConfig, params)

}

func (s *store) QueryBySql(ctx context.Context, sqlConfig *config.SQLConfig, params *v1.SimpleParams) ([]map[string]interface{}, error) {
	dbResult := make([]map[string]interface{}, 0)
	db := s.db.Raw(sqlConfig.Sql, params.Params.AsMap()).Find(&dbResult)
	return dbResult, db.Error
}
