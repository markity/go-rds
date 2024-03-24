package resp

type parseDoubleStateEnum int

const (
	stateDoubleWantingMore parseDoubleStateEnum = iota
)

type parseDoubleState struct {
	state parseDoubleStateEnum
}
