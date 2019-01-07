package search

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"github.com/plopezm/cloud-kaiser/core/logger"
	"github.com/plopezm/cloud-kaiser/core/types"
	"log"
	"time"
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

func (r *ElasticSearchRepository) InsertTask(ctx context.Context, task types.Task) error {
	_, err := r.client.Index().
		Index("tasks").
		Type("JobTask").
		Id(fmt.Sprintf("%s:%s", task.Name, task.Version)).
		BodyJson(task).
		Refresh("wait_for").
		Do(ctx)
	return err
}

func (r *ElasticSearchRepository) FindTasks(ctx context.Context, query string, offset uint64, limit uint64) ([]types.Task, error) {
	result, err := r.client.Search().
		Index("tasks").
		Query(
			elastic.NewMultiMatchQuery(query, "name").
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
	tasks := []types.Task{}
	for _, hit := range result.Hits.Hits {
		var task types.Task
		if err = json.Unmarshal(*hit.Source, &task); err != nil {
			log.Println(err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *ElasticSearchRepository) InsertLog(ctx context.Context, jobname string, jobversion string, taskname string, taskversion string, logline string) error {
	_, err := r.client.Index().
		Index("logs").
		Type("joblog").
		Id(fmt.Sprintf("%s:%s:%s:%s:%d", jobname, jobversion, taskname, taskversion, time.Now().UnixNano())).
		BodyJson(map[string]interface{}{
			"job":         jobname,
			"jobVersion":  jobversion,
			"task":        taskname,
			"taskVersion": taskversion,
			"msg":         logline,
			"ts":          time.Now().UnixNano(),
		}).
		Refresh("wait_for").
		Do(ctx)
	return err
}

func (r *ElasticSearchRepository) FindLogs(ctx context.Context, query string, fields []string, offset uint64, limit uint64) ([]interface{}, error) {
	result, err := r.client.Search().
		Index("logs").
		Query(
			elastic.NewMultiMatchQuery(query, fields...).
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
	logs := make([]interface{}, 0)
	for _, hit := range result.Hits.Hits {
		var log interface{}
		if err = json.Unmarshal(*hit.Source, &log); err != nil {
			logger.GetLogger().Error(err.Error())
		}
		logs = append(logs, log)
	}
	return logs, nil
}
