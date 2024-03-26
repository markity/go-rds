package datastructure

import "time"

type TypeEnum int

const (
	TypeString TypeEnum = iota // string内部包含两个编码, int或str
	TypeHash                   // 也就是键值对
	TypeList                   // 双向链表
	TypeSet                    // 无序集合
)

type EncodingEnum int

const (
	EncodingStringInt EncodingEnum = iota
	EncodingStringRaw

	EncodingHash
	EncodiingSet
)

type RdsObject struct {
	Type     TypeEnum
	Encoding EncodingEnum
	TTL      *time.Time
	Data     interface{}
}
