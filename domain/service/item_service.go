package service

import (
	"github.com/yoshinorihisakawa/qiita-batch/domain/model"
	"github.com/yoshinorihisakawa/qiita-batch/domain/repository"
	"time"
)

type itemService struct {
	ItemRepository repository.ItemRepository
}

type ItemService interface {
	Get(fromDate, ToDate time.Time) ([]*model.Item, error)
	Store(item []*model.Item) error
}

func NewItemService(itemRepository repository.ItemRepository) ItemService {
	return &itemService{ItemRepository: itemRepository}
}

func (itemService *itemService) Get(fromDate, ToDate time.Time) ([]*model.Item, error) {
	var items []*model.Item

	// 100ページまでしか取得できないので最大値の100を指定
	for i := 1; i <= 100; i++ {
		result, err := itemService.ItemRepository.FetchAll(fromDate, ToDate, i)
		if err != nil {
			// TODO:リトライとか考える
			return nil, err
		}

		// 全件取得した時点で返却
		if len(result) == 0 {
			return items, nil
		}

		items = append(items, result...)

	}
	return items, nil
}

func (itemService *itemService) Store(item []*model.Item) error {
	return itemService.ItemRepository.Store(item)
}
