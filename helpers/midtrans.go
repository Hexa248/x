package helpers

import (
	"fmt"
	"time"
)

func CreatePaymentToken(bookingID uint) string {
	return fmt.Sprintf("PAY-%d-%d", bookingID, time.Now().Unix())
}
