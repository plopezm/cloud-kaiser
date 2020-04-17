package search

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	elastic "github.com/olivere/elastic/v7"
	"github.com/plopezm/cloud-kaiser/core/logger"
	"github.com/plopezm/cloud-kaiser/core/types"
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
			elastic.NewWildcardQuery("name", query),
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
		if err = json.Unmarshal(hit.Source, &task); err != nil {
			log.Println(err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *ElasticSearchRepository) InsertLog(ctx context.Context, taskExecutionLog types.TaskExecutionLog) error {
	_, err := r.client.Index().
		Index("logs").
		Type("JobLog").
		Id(fmt.Sprintf("%s:%s:%s:%s:%d", taskExecutionLog.JobName, taskExecutionLog.JobVersion,
			taskExecutionLog.TaskName, taskExecutionLog.TaskVersion,
			taskExecutionLog.Ts)).
		BodyJson(taskExecutionLog).
		Refresh("wait_for").
		Do(ctx)
	return err
}

func (r *ElasticSearchRepository) FindLogs(ctx context.Context, query string, field string, offset uint64, limit uint64) ([]types.TaskExecutionLog, error) {
	initialFields := []string{"jobName", "jobVersion", "taskName", "taskVersion", "line"}

	highlightFields := make([]*elastic.HighlighterField, 0)
	for _, field := range initialFields {
		highlightFields = append(highlightFields, elastic.NewHighlighterField(field))
	}
	result, err := r.client.Search().
		Index("logs").
		Query(
			elastic.NewBoolQuery().Must(
				elastic.NewBoolQuery().Must(
					elastic.NewBoolQuery().
						MinimumShouldMatch("1").
						Should(
							elastic.NewMultiMatchQuery(query, initialFields...).
								Type("cross_fields").
								Operator("and"),
							elastic.NewMultiMatchQuery(query, initialFields...).
								Type("phrase_prefix").
								Operator("and"),
						),
				),
			),
		).
		Highlight(
			elastic.NewHighlight().Fields(
				highlightFields...,
			),
		).
		TrackScores(true).
		From(int(offset)).
		Size(int(limit)).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	logs := make([]types.TaskExecutionLog, 0)
	for _, hit := range result.Hits.Hits {
		var log types.TaskExecutionLog
		if err = json.Unmarshal(hit.Source, &log); err != nil {
			logger.GetLogger().Error(err.Error())
		}
		logs = append(logs, log)
	}
	return logs, nil
}
