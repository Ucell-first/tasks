package errors

import "errors"

var (
	// ErrConfigurationIsNotForLivecache appears when passed pointer to configuration
	// structure is not for Livecache provider.
	ErrConfigurationIsNotForLivecache = errors.New("passed configuration is not for livecache provider")

	// ErrExpiredIntervalNotSet appears when ExpiredInterval is not set.
	ErrExpiredIntervalNotSet = errors.New("expiration interval not set")

	// ErrClearCycleMechanismNotSet appears when clear cycle mechanism is not set
	// (ClearInterval == 0 and MaxElements == 0)
	ErrClearCycleMechanismNotSet = errors.New("clear cycle mechanism not set")

	// ErrKeyNotFound appears when key is not exists in cache.
	ErrKeyNotFound = errors.New("key not found")

	// ErrExpiredCacheItem appears when key has expired.
	ErrExpiredCacheItem = errors.New("cache item has expired")
)
