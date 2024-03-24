package resp

type parseBulkStringStateEnum int

const (
	parseBulkStringWantingLength parseBulkStringStateEnum = iota
	parseBulkStringWantingData
)

type parseBulkStringState struct {
	state  parseBulkStringStateEnum
	length int64
}
