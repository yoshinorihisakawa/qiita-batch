package main

import (
	"fmt"
	"github.com/yoshinorihisakawa/qiita-batch/domain/service"
	"github.com/yoshinorihisakawa/qiita-batch/infrastructure"
	"github.com/yoshinorihisakawa/qiita-batch/usecase"
	"time"
	"os"
	log "github.com/sirupsen/logrus"
	"github.com/yoshinorihisakawa/qiita-batch/conf"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func main() {
	// conf読み取り
	conf.ReadConf()

	var log = log.New()
	// 依存関係を定義
	db := infrastructure.NewMySqlDB()
	dataBase := infrastructure.NewItemRepository(db)
	service := service.NewItemService(dataBase)
	usecase := usecase.NewItemUseCase(service)

	// 開始日時を指定
	fromDate := time.Date(2018, 5, 30, 0, 0, 0, 0, time.Local)

	// 月末を取得
	toDate := fromDate.AddDate(0, 1, -1 * fromDate.Day())

	// 周期の時刻を決める。
	t := time.NewTicker(3 * time.Second)

	// ループ処理開始
	for {
		select {
		case <- t.C:
			log.Info(fmt.Sprintf("start：%d-%d-%d - %d-%d-%d\n", fromDate.Year(), int(fromDate.Month()), fromDate.Day(),
				toDate.Year(), int(toDate.Month()), toDate.Day()))

			// 開始日からその月のデータを全て取得
			qiitaData, err := usecase.GetQiitaData(fromDate, toDate)
			if err != nil {
				log.Error(err)
			}

			// データを保存 インサートアップデート
			if err := usecase.StoreQiitaData(qiitaData); err != nil {
				log.Error(err)
			}

			// ファイルにログを主力
			log.Info(fmt.Sprintf("end：%d-%d-%d - %d-%d-%d\n", fromDate.Year(), int(fromDate.Month()), fromDate.Day(),
				toDate.Year(), int(toDate.Month()), toDate.Day()))

			// 次の開始日を決定
			fromDate, toDate = getNextPeriod(toDate)

			// 現在時刻をすぎたら終了
			if fromDate.Unix() > time.Now().Unix() {
				log.Info("バッチ終了")
				t.Stop()
				return
			}
		}
	}
}
func getNextPeriod(t time.Time) (fromDate, toDate time.Time) {
	// 一日足した時刻
	fromDate = t.AddDate(0, 0, 1)

	// 月末
	toDate = fromDate.AddDate(0, 1, -1 * fromDate.Day())
	return
}
