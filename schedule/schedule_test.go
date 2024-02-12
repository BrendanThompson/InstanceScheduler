package schedule

import (
	// "github.com/stretchr/testify/assert"
	"testing"
)

func TestSchedules(t *testing.T) {
	testCases := []struct {
		name string
		data string
		want bool
	}{
		{
			name: "default_single_window",
			data: `{"default":["09:00-17:00"]}`,
			want: true,
		},
		{
			name: "default_invalid_time",
			data: `{"default":["18:00-17:00"]}`,
			want: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			s, err := NewSchedule([]byte(test.data))
			if err != nil {
				t.Error(err)
			}

			got := s.Validate()

			if got != test.want {
				t.Errorf("got: %t, want: %t", got, test.want)
			}
		})
	}
}
