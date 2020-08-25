package common

import "goCache"

// 实体
type Entry struct {
	Key   string
	Value interface{}
}

func (e *Entry) Len() int {
	return goCache.CalcLen(e.Value)
}

// 回调
type OnEvicted func(key string, value interface{})
