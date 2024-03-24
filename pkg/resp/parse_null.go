package resp

// 需要和simple string一样解析, 然后手动判断message是否合理
type parseNullStateEnum int

const (
	parseNullWantingMore parseNullStateEnum = iota
)

type parseNullState struct {
	state parseNullStateEnum
}
