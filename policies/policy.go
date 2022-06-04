package policies

type EvictionPolicy interface {
	// SetCapacity sets m-cache max capacity.
	SetCapacity(capacity int64)
	// Capacity returns m-cache capacity.
	Capacity() int64
	// Promote promotes the 'm-cache key', if 'm-cache key' doesn't exist, it will create the 'm-cache key'.
	Promote(key string)
	// PromoteIfExist promotes the 'm-cache key' if 'm-cache key' exists.
	PromoteIfExist(key string)
	// Evict evicts the 'm-cache key'.
	Evict(key string)
	// Ban evicts the 'm-cache key' and prevent subsequent 'promote'.
	Ban(key string)
	// NowEvict evict the 'm-cache key' and return it by eviction policy
	NowEvict() (key string)
}
