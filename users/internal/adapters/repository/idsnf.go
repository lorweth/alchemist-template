package repository

import (
	"github.com/sony/sonyflake"
)

// We're using snoyflake for generate unique id instead of postgres running number, it helps for data migration
var (
	userIDGenerator *sonyflake.Sonyflake
)

func initIDGenerator() {
	if userIDGenerator == nil {
		userIDGenerator = sonyflake.NewSonyflake(sonyflake.Settings{})
	}
}
