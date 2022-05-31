package dict

type ConcurrentMap interface {
	Put(key string, val interface{}) (result int)
	Get(key string) (val interface{}, exists bool)
	Len() int
	PutIfAbsent(key string, val interface{}) (result int)
	PutIfExists(key string, val interface{}) (result int)
	Remove(key string) (result int)
	ForEach(recallFunc RecallFunc)
}

type RecallFunc func(key string, val interface{}) bool
