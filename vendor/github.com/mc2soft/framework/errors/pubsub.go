package errors

import errs "errors"

var (
	// ErrPubSubMessageTopicUndefined appears when trying to send a message without
	// defining topic in it.
	ErrPubSubMessageTopicUndefined = errs.New("pub-sub topic in passed message is undefined")

	// ErrPubSubTopicAlreadyExists appears when trying to create a topic in pubsub with already
	// used name.
	ErrPubSubTopicAlreadyExists = errs.New("pub-sub topic already exists")

	// ErrPubSubTopicNotFound appears when trying to send a message but topic for which
	// this message should belong isn't exists.
	ErrPubSubTopicNotFound = errs.New("pub-sub topic wasn't found")
)
