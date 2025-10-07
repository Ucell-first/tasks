# tasks


Library for task management.

## Configuring

Library uses functional options pattern. Available configuration options (which should be passed to `New`):

| Option | Description                      |
| --- |----------------------------------|
| `WithContext(context.Context)` | Application's context, used for HTTP requests. |
| `WithLogger(logger.Logger)` | Logger to use. See `logger.Logger` interface. |
| `WithProvider(provider communication.Provider, topic string)` |                                  |
| `WithNumWorkers(numWorkers int)` |                                  |
| `WithNumWorkers(numWorkers int)` |                                  |
| `WithQueueSize(queueSize int)` |                                  |
| `WithRetryPolicy(retryPolicy models.RetryPolicy)` |                                  |


## Using

```go
func (d *domain) LaunchStartupTasks() error {
	
	var err error

	d.kafkaProvider, err = d.app.Communication.GetProvider(framework.Kafka)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrLaunchStartupTasks, err)
	}

	d.tasker, err = tasks.New(tasks.WithContext(ado.app.GetContext()),
		tasks.WithLogger(ado.app.Logger()),
		tasks.WithProvider(ado.kafkaProvider, "topic"),
		tasks.WithQueueSize(100),
		tasks.WithRetryPolicy(models.RetryPolicy{
		    BackoffCoefficient: 2.0,
			MaximumAttempts: 10,
		}),
	)
	if err != nil {
		return fmt.Errorf("%w:: %w", ErrLaunchStartupTasks, err)
	}
	
	

	err = d.tasker.RegisterHandler("check_status", d.checkStatus)
	if err != nil {
		return fmt.Errorf("%w:: %w", ErrLaunchStartupTasks, err)
	}
	
	return nil
}

func (d *domain) Shutdown() error {
	// Stop task processing
  d.tasker.Stop()
  
  return nil
}

```

```go
// task create example
func (d *domain) TaskCreate() error {
		err := d.tasker.Create(ado.app.GetContext(), "check_status", map[string]string{})
		if err != nil {
			d.logger.Logf("ERROR", "task create:%v", nil, err)
		}
	}
	
	return nil
}

```

```go
func (d *domain) checkStatus(params map[string]string) error {
	/** 
	 implement logic
	*/
	return nil
}

```# tasks
