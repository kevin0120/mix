package storage

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeUTC(t *testing.T) {
	dateString := "2019-10-16T11:20:30+08:00"

	loc, _ := time.LoadLocation("Local")
	dt, _ := time.ParseInLocation(time.RFC3339, dateString, loc)

	fmt.Println(dt.Format(time.RFC3339))
}
