package resp

type parseBigNumStateEnum int

const (
	stateBigNumWantingMore parseBigNumStateEnum = iota
)

type parseBigNumState struct {
	state parseBigNumStateEnum
}
