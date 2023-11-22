package repository

import (
	"time"
)

var timeNowWrapper = func() time.Time {
	return time.Now()
}
