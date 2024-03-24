package resp

// 需要和simple string一样解析, 然后手动判断message是否合理
type parseSimpleErrorStateEnum int

const (
	parseSimpleErrorStateWantingMore parseSimpleErrorStateEnum = iota
)

type parseSimpleErrorState struct {
	state parseSimpleErrorStateEnum
}
