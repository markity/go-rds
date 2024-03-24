package resp

// 需要和simple string一样解析, 然后手动判断message是否合理
type parseIntegerStateEnum int

const (
	parseIntegerWantingMore parseIntegerStateEnum = iota
)

type parseIntegerState struct {
	state parseIntegerStateEnum
}
