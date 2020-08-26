package elasticsearch

import "github.com/olivere/elastic/v7"

func NewClient(url string) *elastic.Client {
	urls := []string{url}
	client, err := elastic.NewClient(elastic.SetURL(urls...))
	if err != nil {
		panic(err)
	}
	return client
}
