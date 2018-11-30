package usecase

import (
	"time"
	"github.com/yoshinorihisakawa/qiita-batch/domain/model"
	"github.com/yoshinorihisakawa/qiita-batch/domain/service"
)

type itemUseCase struct {
	ItemService service.ItemService
}

type ItemUseCase interface {
	GetQiitaData(fromDate, ToDate time.Time) ([]*model.Item, error)
	StoreQiitaData(items []*model.Item) error
}

func NewItemUseCase(itemService service.ItemService) ItemUseCase {
	return &itemUseCase{ItemService:itemService}
}

func (itemUseCase *itemUseCase) GetQiitaData(fromDate, ToDate time.Time) ([]*model.Item, error){
	// Qiitaから時間を指定してデータを取得
	items,err := itemUseCase.ItemService.Get(fromDate, ToDate)
	if err != nil {
		return nil,err
	}
	return items,nil
}

func (itemUseCase *itemUseCase) StoreQiitaData(items []*model.Item) error {
	err := itemUseCase.ItemService.Store(items)
	if err != nil {
		return err
	}
	return nil
}
