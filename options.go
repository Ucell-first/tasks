package tasks

import (
	"context"

	"github.com/mc2soft/framework/communication"
	"gitlab.local.iti.domain/mc2/golibs/tasks/logger"
	"gitlab.local.iti.domain/mc2/golibs/tasks/models"
)

type options struct {
	ctx         context.Context
	logger      logger.Logger
	provider    communication.Provider
	topic       string
	numWorkers  int
	queueSize   int
	retryPolicy models.RetryPolicy
}

// Option is an interface for configuration options.
type Option interface {
	apply(o *options)
}

type contextOption struct {
	ctx context.Context
}

func (co *contextOption) apply(o *options) {
	o.ctx = co.ctx
}

// WithContext sets context for translations library.
func WithContext(ctx context.Context) Option {
	return &contextOption{ctx: ctx}
}

type loggerOption struct {
	lgr logger.Logger
}

func (lo loggerOption) apply(o *options) {
	o.logger = lo.lgr
}

// WithLogger sets logger implementation.
func WithLogger(logger logger.Logger) Option {
	return &loggerOption{lgr: logger}
}

type providerOption struct {
	provider communication.Provider
	topic    string
}

func (po *providerOption) apply(o *options) {
	o.provider = po.provider
	o.topic = po.topic
}

func WithProvider(provider communication.Provider, topic string) Option {
	return &providerOption{provider: provider, topic: topic}
}

type WorkerOption struct {
	numWorkers int
}

func (wo *WorkerOption) apply(o *options) {
	o.numWorkers = wo.numWorkers
}

func WithNumWorkers(numWorkers int) Option {
	return &WorkerOption{numWorkers: numWorkers}
}

type queueSizeOption struct {
	queueSize int
}

func (qo *queueSizeOption) apply(o *options) {
	o.queueSize = qo.queueSize
}

func WithQueueSize(queueSize int) Option {
	return &queueSizeOption{queueSize: queueSize}
}

type retryPolicyOption struct {
	retryPolicy models.RetryPolicy
}

func (ro *retryPolicyOption) apply(o *options) {
	o.retryPolicy = ro.retryPolicy
}

func WithRetryPolicy(retryPolicy models.RetryPolicy) Option {
	return &retryPolicyOption{retryPolicy: retryPolicy}
}
