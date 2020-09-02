package collectors_test

import (
	"testing"
	"time"

	"github.com/fabiodmferreira/crypto-trading/collectors"
)

func TestGetPreviousIntervalDates(t *testing.T) {

	gotSD, gotED := collectors.GetPreviousIntervalDates(time.Date(2019, 3, 1, 9, 23, 23, 0, time.UTC))
	wantSD, wantED := time.Date(2019, 3, 1, 9, 22, 0, 0, time.UTC), time.Date(2019, 3, 1, 9, 23, 0, 0, time.UTC)

	if gotSD != wantSD {
		t.Errorf("start date: got %v want %v", gotSD, wantSD)
	}

	if gotED != wantED {
		t.Errorf("end date: got %v want %v", gotED, wantED)
	}

}
