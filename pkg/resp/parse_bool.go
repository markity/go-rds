package resp

type parseBoolStateEnum int

const (
	parseBoolWantingMore parseBoolStateEnum = iota
)

type parseBoolState struct {
	state parseBoolStateEnum
}
