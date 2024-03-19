package data

import "gorm.io/gorm"

type ShareDaoFactory interface {
	Store() StoreInterface
}

type shareDaoFactory struct {
	db *gorm.DB
}

func (s *shareDaoFactory) Store() StoreInterface {
	return newStore(s.db)
}

func NewShareDaoFactory(db *gorm.DB) (ShareDaoFactory, error) {
	return &shareDaoFactory{
		db: db,
	}, nil
}
