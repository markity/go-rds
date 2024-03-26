package commands

import datastructure "go-rds/data_structure"

type SetCommand struct {
	Key      string
	Encoding datastructure.EncodingEnum
	ValueRaw string
	ValueInt int64
	Nx       bool
	Xx       bool
	Ex       *int
	Px       *int
}

type GetCommand struct {
	Key string
}

type TTLCommand struct {
	Key string
}

type PTTLCommand struct {
	Key string
}

type AppendCommand struct {
	Key     string
	Content string
}

type IncrCommand struct {
	Key string
}

type IncrByCommand struct {
	Key string
	By  int64
}

type DecrCommand struct {
	Key string
}

type DecrByCommand struct {
	Key string
	By  int64
}

type GetDelCommand struct {
	Key string
}

type StrLenCommand struct {
	Key string
}

type GetSetCommand struct {
	Key      string
	Encoding datastructure.EncodingEnum
	ValueRaw string
	ValueInt int64
}

type GetExCommand struct {
	Key     string
	Ex      *int
	Px      *int
	Persist bool
}

type SetNxCommand struct {
	Key      string
	Encoding datastructure.EncodingEnum
	ValueRaw string
	ValueInt int64
}
