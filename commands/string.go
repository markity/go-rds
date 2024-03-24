package commands

import datastructure "go-rds/data_structure"

type SetCommand struct {
	Key      string
	Encoding datastructure.EncodingEnum
	ValueRaw string
	ValueInt int64
}
