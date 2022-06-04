package policies

import "math"

type Non struct {
	maxCap int
}

func NewNon() *Non {
	return &Non{}
}

func (n *Non) SetCapacity(capacity int64) {
}

func (n *Non) Capacity() int64 {
	return math.MaxInt64
}

func (n *Non) Promote(key string) {
}

func (n *Non) PromoteIfExist(key string) {
}

func (n *Non) Evict(key string) {
}

func (n *Non) Ban(key string) {
}

func (n *Non) NowEvict() (key string) {
	return ""
}
