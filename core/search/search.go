package search

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"github.com/plopezm/cloud-kaiser/core/types"
	"log"
)

type ElasticSearchRepository struct {
	client *elastic.Client
}

func NewElasticSearch(url string) (*ElasticSearchRepository, error) {
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}
	return &ElasticSearchRepository{
		client: client,
	}, nil
}

func (r *ElasticSearchRepository) Close() {
}

func (r *ElasticSearchRepository) InsertTask(ctx context.Context, task types.JobTask) error {
	_, err := r.client.Index().
		Index("tasks").
		Type("JobTask").
		Id(fmt.Sprintf("%s:%s", task.Name, task.Version)).
		BodyJson(task).
		Refresh("wait_for").
		Do(ctx)
	return err
}

func (r *ElasticSearchRepository) FindTasks(ctx context.Context, query string, offset uint64, limit uint64) ([]types.JobTask, error) {
	result, err := r.client.Search().
		Index("tasks").
		Query(
			elastic.NewMultiMatchQuery(query, "body").
				Fuzziness("3").
				PrefixLength(1).
				CutoffFrequency(0.0001),
		).
		From(int(offset)).
		Size(int(limit)).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	tasks := []types.JobTask{}
	for _, hit := range result.Hits.Hits {
		var task types.JobTask
		if err = json.Unmarshal(*hit.Source, &task); err != nil {
			log.Println(err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
