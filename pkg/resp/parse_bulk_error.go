package resp

type parseBulkErrorStateEnum int

const (
	parseBulkErrorWantingLength parseBulkErrorStateEnum = iota
	parseBulkErrorWantingData
)

type parseBulkErrorState struct {
	state  parseBulkErrorStateEnum
	length int64
}
