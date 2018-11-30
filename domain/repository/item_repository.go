package repository

import (
	"time"
	"github.com/yoshinorihisakawa/qiita-batch/domain/model"
)

type ItemRepository interface {
	FetchAll(fromDate, ToDate time.Time, i int) ([]*model.Item, error)
	Store(item []*model.Item) error
}