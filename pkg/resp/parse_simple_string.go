package resp

type parseSimpleStringStateEnum int

const (
	parseSimpleStringStateWantingMore parseSimpleStringStateEnum = iota
)

type parseSimpleStringState struct {
	state parseSimpleStringStateEnum
}
