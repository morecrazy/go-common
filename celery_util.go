package common

import (
	"third/go-celery"
)

func PublishTask(celery_client *celery.Celery, task_name, queue string, args_param ...interface{}) error {

	var args []interface{}
	args = append(args, args_param...)
	task, err := celery.NewTask(task_name, args, nil)
	if nil != err {
		return &InternalError{InternalErrorCode, err}
	}

	publisher, err := celery.NewPublishing(task, queue)
	if nil != err {
		return &InternalError{InternalErrorCode, err}
	}

	err = celery_client.Publish(publisher)
	if nil != err {
		return &InternalError{InternalErrorCode, err}
	}

	return nil
}
