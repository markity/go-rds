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
