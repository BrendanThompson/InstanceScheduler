/*
Copyright Brendan Thompson

Licensed under the PolyForm Internal Use License, Version 1.0.0 (the "License");
you may not use this file except in compliance with the License.
A copy of the License may be obtained at

https://polyformproject.org/licenses/internal-use/1.0.0/
*/

package schedule

import (
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
