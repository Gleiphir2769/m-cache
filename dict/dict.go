package dict

type ConcurrentMap interface {
	// Put k, v anyway, result 1 means a new key added.
	Put(key string, val interface{}) (result int)
	// Get v by k
	Get(key string) (val interface{}, exists bool)
	// Len returns the length of map
	Len() int
	// PutIfAbsent return 1 when key isn't existed (a new key added).
	PutIfAbsent(key string, val interface{}) (result int)
	// PutIfExists return 1 when key is existed (a new key added).
	PutIfExists(key string, val interface{}) (result int)
	// Remove return 1 when an existed key removed.
	Remove(key string) (result int)
	// ForEach calls the recallFunc on all elements.
	ForEach(recallFunc RecallFunc)
}

type RecallFunc func(key string, val interface{}) bool
