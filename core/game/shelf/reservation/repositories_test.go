package reservation

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/test"
	"testing"
	"time"
)

func TestReservationRowsToUserIdArray(t *testing.T) {
	type testCase struct {
		rows         []*ReservationRow
		expectedUser []core.UserId
	}

	testCases := []testCase{
		{
			expectedUser: []core.UserId{
				"1", "2",
			},
			rows: []*ReservationRow{
				{
					Id:            "s1",
					UserId:        "1",
					Index:         0,
					ScheduledTime: time.Time{},
					PurchaseNum:   1,
				},
				{
					Id:            "s2",
					UserId:        "1",
					Index:         0,
					ScheduledTime: time.Time{},
					PurchaseNum:   1,
				},
				{
					Id:            "s3",
					UserId:        "2",
					Index:         0,
					ScheduledTime: time.Time{},
					PurchaseNum:   1,
				},
			},
		},
	}

	for _, tc := range testCases {
		actual := ReservationRowsToUserIdArray(tc.rows)
		if !test.DeepEqual(actual, tc.expectedUser) {
			t.Errorf("got: %v, want: %v", actual, tc.expectedUser)
		}
	}
}
