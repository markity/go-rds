package global

import (
	"sync/atomic"
)

var ConnIDCounter atomic.Int64
