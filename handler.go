package tasks

import (
	"encoding/json"
	"fmt"

	comContext "github.com/mc2soft/framework/communication/context"
	errKafka "github.com/mc2soft/framework/errors"
	"gitlab.local.iti.domain/mc2/golibs/tasks/models"
)

func (t *Tasks) handleTask(ctx comContext.Context) error {
	var task models.Task

	err := json.NewDecoder(ctx.Body()).Decode(&task)
	if err != nil {
		return fmt.Errorf("%w: %w", errHandler, err)
	}

	if !t.AreConsumersActive.Load() {
		return errKafka.ErrKafkaDoNotSkipMessage
	}

	t.taskQueue <- task

	return nil
}
