package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/honlyc/elasticsearch-dump/config"
	"github.com/honlyc/elasticsearch-dump/pkg/elasticsearch"
	"github.com/olivere/elastic/v7"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"os"
	"strings"
)

type Res struct {
	IndexedAt       string `json:"indexedAt"`
	RetweetUserID   string `json:"retweetUserId"`
	VipType         int    `json:"vipType"`
	RetweetID       int64  `json:"retweetId"`
	Label           string `json:"label"`
	UserID          string `json:"userId"`
	Platform        string `json:"platform"`
	RepostsCount    int    `json:"repostsCount"`
	CreatedAt       int64  `json:"createdAt"`
	ReplyStatusID   int    `json:"replyStatusId"`
	CurrentContent  string `json:"currentContent"`
	FavourCount     string `json:"favourCount"`
	Location        string `json:"location"`
	ID              string `json:"id"`
	Vip             bool   `json:"vip"`
	OriginalContent string `json:"originalContent"`
	Username        string `json:"username"`
}

func main() {
	client := elasticsearch.NewClient(config.CONFIG.Es.Cluster)
	testIndexName := config.CONFIG.Es.IndexName
	query := config.CONFIG.Es.Query
	fmt.Printf("query: %s\n", query)
	split := strings.Split(query, ",")
	queryStr := split[1]
	title := split[0]
	queryString := elastic.NewQueryStringQuery(queryStr).DefaultOperator("AND").DefaultField("entireContent")
	total, err := client.Count(testIndexName).Query(queryString).Do(context.TODO())

	if err != nil {
		panic(err)
	}

	fmt.Printf("total: %d\n", total)

	bar := pb.StartNew(int(total))

	// 1st goroutine sends individual hits to channel.
	hits := make(chan json.RawMessage)
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		defer close(hits)
		// Initialize scroller. Just don't call Do yet.
		scroll := client.Scroll(testIndexName).Size(config.CONFIG.Es.Size).Query(queryString)
		for {
			results, err := scroll.Do(context.TODO())
			if err == io.EOF {
				return nil // all results retrieved
			}
			if err != nil {
				return err // something went wrong
			}

			// Send the hits to the hits channel
			for _, hit := range results.Hits.Hits {
				select {
				case hits <- hit.Source:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}

		//err = scroll.Clear(context.TODO())
		//if err != nil {
		//	panic(err)
		//}
		//
		//_, err = scroll.Do(context.TODO())
		//if err == nil {
		//	panic("expected to fail")
		//}
		return nil
	})

	file := "./" + title + ".json"

	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0766)
	if nil != err {
		panic(err)
	}

	//创建一个Logger
	//参数1：日志写入目的地
	//参数2：每条日志的前缀
	//参数3：日志属性
	loger := log.New(logFile, "", 0)

	for i := 0; i < 10; i++ {
		g.Go(func() error {
			for hit := range hits {
				// Deserialize
				item := &Res{}
				err := json.Unmarshal(hit, &item)
				if err != nil {
					return err
				}

				// Do something with the product here, e.g. send it to another channel
				// for further processing.
				//_ = p
				itemJson, _ := json.Marshal(item)
				loger.Printf("%s", string(itemJson))
				bar.Increment()

				// Terminate early?
				select {
				default:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return nil
		})
	}

	// Check whether any goroutines failed.
	if err := g.Wait(); err != nil {
		panic(err)
	}

	// Done.
	bar.Finish()
}
