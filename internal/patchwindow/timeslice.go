/*
Copyright Brendan Thompson

Licensed under the PolyForm Internal Use License, Version 1.0.0 (the "License");
you may not use this file except in compliance with the License.
A copy of the License may be obtained at

https://polyformproject.org/licenses/internal-use/1.0.0/
*/

package patchwindow

import (
	"errors"
	"time"
)

func NewTimesliceWithDuration(start string, duration int) (*Timeslice, error) {
	var timeslice Timeslice

	now := time.Now().Local()

	parsedStart, err := time.Parse("15:04", start)
	if err != nil {
		return nil, errors.New("Unable to parse start time")
	}

	parsedEnd := parsedStart.Add(time.Hour * time.Duration(duration))

	timeslice.Start = time.Date(now.Year(), now.Month(), now.Day(), parsedStart.Hour(), parsedStart.Minute(), 0, 0, now.Location())
	timeslice.End = time.Date(now.Year(), now.Month(), now.Day(), parsedEnd.Hour(), parsedEnd.Minute(), 0, 0, now.Location())

	return &timeslice, nil
}

type Timeslice struct {
	Start time.Time
	End   time.Time
}
