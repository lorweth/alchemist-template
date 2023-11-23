package generator

import (
	"github.com/sony/sonyflake"
)

// We're using sonyflake for generate unique id instead of postgres running number, it helps for data migration
var (
	UserIDGenerator *sonyflake.Sonyflake
)

func InitIDGenerator() {
	if UserIDGenerator == nil {
		UserIDGenerator = sonyflake.NewSonyflake(sonyflake.Settings{})
	}
}
