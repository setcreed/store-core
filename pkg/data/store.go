package data

import (
	"context"
	v1 "github.com/setcreed/store-core/api/store_service/v1"
	"gorm.io/gorm"
)

type StoreInterface interface {
	QueryByTableName(ctx context.Context, tableName string, params *v1.SimpleParams) ([]map[string]interface{}, error)
}

type store struct {
	db *gorm.DB
}

func newStore(db *gorm.DB) StoreInterface {
	return &store{
		db: db,
	}
}

func (s *store) QueryByTableName(ctx context.Context, tableName string, params *v1.SimpleParams) ([]map[string]interface{}, error) {
	dbResult := make([]map[string]interface{}, 0)
	db := s.db.Table(tableName)
	//for k, v := range params.Params.AsMap() {
	//	db = db.Where(k, v)
	//}
	db = db.Find(&dbResult)

	return dbResult, db.Error
}
