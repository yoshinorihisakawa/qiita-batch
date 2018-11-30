package infrastructure

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/yoshinorihisakawa/qiita-batch/domain/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
	"github.com/jinzhu/gorm"
)

// アクセストークン取得
var qiitaToken = os.Getenv("QIITA_TOKEN")

type itemRepository struct {
	db *gorm.DB
}

type ItemRepository interface {
	FetchAll(fromDate, ToDate time.Time, i int) ([]*model.Item, error)
	Store(item []*model.Item) error
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{db}
}

func (itemRepository *itemRepository) FetchAll(fromDate, ToDate time.Time, i int) ([]*model.Item, error) {

	page := strconv.Itoa(i)
	perPage := "100"
	requestUrl := "https://qiita.com/api/v2/items?page=" + page + "&per_page=" + perPage

	// monthだけint型に型変換
	year, month, day := fromDate.Date()
	fDate := dateNum2String(year, int(month), day)

	year, month, day = ToDate.Date()
	tDate := dateNum2String(year, int(month), day)

	// 投稿の検索クエリを作成
	// 検索クエリ stocks:>NUM created:<=YYYY-MM-DD created:>YYYY-MM-DD
	// 指定日に投稿されたストック数30以上の記事を取得
	varParam := "&query=stocks:>10+created:>=" + fDate + "+created:<" + tDate

	endpointURL, err := url.Parse(requestUrl + varParam)
	if err != nil {
		return nil, errors.Errorf("parse error: %s", err)
	}

	b, err := json.Marshal(model.Item{})
	if err != nil {
		return nil, errors.Errorf("JSON marshal error: %s", err)
	}

	var resp = &http.Response{}
	// qiitaのアクセストークンがない場合はAuthorizationを付与しない
	if len(qiitaToken) > 0 {
		resp, err = http.DefaultClient.Do(&http.Request{
			URL:    endpointURL,
			Method: "GET",
			Header: http.Header{
				"Content-Type": {"application/json"},
				// https://qiita.com/api/v2/docs#%E8%AA%8D%E8%A8%BC%E8%AA%8D%E5%8F%AF
				"Authorization": {"Bearer " + qiitaToken},
			},
		})
	} else {
		resp, err = http.DefaultClient.Do(&http.Request{
			URL:    endpointURL,
			Method: "GET",
			Header: http.Header{
				"Content-Type": {"application/json"},
			},
		})
	}
	defer resp.Body.Close()

	// 200以外ならエラーとして処理する。
	if err != nil {
		return nil, errors.Errorf("response error: %s", err)
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Errorf("read error: %s", err)
	}

	var items []*model.Item
	if err := json.Unmarshal(b, &items); err != nil {
		return nil, errors.Errorf("JSON Unmarshal error: %s", err)
	}

	return items, nil
}

// 年月日の数値を文字列に変換
func dateNum2String(year int, month int, day int) string {
	return fmt.Sprintf("%d-%d-%d", year, month, day)
}
func (itemRepository *itemRepository) Store(items []*model.Item) error {
	for _, i := range items {
		err := itemRepository.db.Save(i).Error
		if err != nil {
			return errors.Errorf("save error: %s", err)
		}
	}
	return nil
}
