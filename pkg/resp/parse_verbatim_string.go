package resp

type parseVerbatimStringStateEnum int

const (
	parseVerbatimStringStateWantingLength parseVerbatimStringStateEnum = iota
	parseVerbatimStringStateWantingData
)

type parseVerbatimStringState struct {
	state  parseVerbatimStringStateEnum
	length int64
}
