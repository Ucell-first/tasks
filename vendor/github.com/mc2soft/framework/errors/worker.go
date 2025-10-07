package errors

import errs "errors"

// ErrWorkerIsStopping appears when application is shutting down and
// some part of it still trying to add item for asynchronous processing
// with worker.
var ErrWorkerIsStopping = errs.New("worker is stopping, cannot accept new items for processing")
