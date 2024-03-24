package resp

// 需要和simple string一样解析, 然后手动判断message是否合理
type parseMapStateEnum int

const (
	parseMapWantingLength parseMapStateEnum = iota
	parseMapWantingData
)

type parseMapState struct {
	state  parseMapStateEnum
	length int
	data   []Value
}

func (st *parseMapState) PushEntry(val Value) {
	st.data = append(st.data, val)
}
