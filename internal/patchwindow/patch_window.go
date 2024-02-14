/*
Copyright Brendan Thompson

Licensed under the PolyForm Internal Use License, Version 1.0.0 (the "License");
you may not use this file except in compliance with the License.
A copy of the License may be obtained at

https://polyformproject.org/licenses/internal-use/1.0.0/
*/

package patchwindow

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func New(data []byte) (*PatchWindow, error) {
	var err error
	var window PatchWindow

	log.Debug().Msg("Creating new patch window from tag.")
	log.Debug().Msg(string(data))

	err = json.Unmarshal(data, &window)
	if err != nil {
		return nil, err
	}

	timeslice, err := NewTimesliceWithDuration(window.Time, window.Duration)
	if err != nil {
		return nil, err
	}

	window.Timeslice = *timeslice

	return &window, nil
}

type PatchWindow struct {
	Period    Period    `json:"period"`
	Week      int       `json:"week"`
	Day       string    `json:"day"`
	Time      string    `json:"time"`
	Duration  int       `json:"duration"`
	Timeslice Timeslice `json:"-"`
}

func (p *PatchWindow) IsToday() bool {
	if p == nil {
		return false
	}

	now := time.Now().Local()

	weekday, err := time.Parse("Monday", p.Day)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Failed to parse weekday")
	}

	if now.Weekday() != weekday.Weekday() {
		return false
	}

	weekNumber := getWeekOfMonth(now)

	if weekNumber != p.Week {
		return false
	}

	log.Debug().Msg("Today is patch day")

	return true
}

func (p *PatchWindow) CurrentTimeWithinRange() bool {
	if !p.IsToday() {
		return false
	}

	now := time.Now().Local()

	if now.After(p.Timeslice.Start.Add(time.Hour*-1)) && now.Before(p.Timeslice.End) {
		return true
	}

	return false
}

func (p *PatchWindow) NextWindowStart() (time.Time, error) {
	if p == nil {
		return time.Time{}, errors.New("Patch window is nil")
	}

	now := time.Now().Local()
	weekday := parseWeekday(p.Day)
	week := p.Week
	startTime, _ := time.Parse("15:04", p.Time)

	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, startTime.Hour(), startTime.Minute(), 0, 0, now.Location())
	firstDayOfWeek := firstDayOfMonth.AddDate(0, 0, int(weekday-firstDayOfMonth.Weekday()))
	day := firstDayOfWeek.AddDate(0, 0, (week-1)*7)

	return day, nil
}

func getWeekOfMonth(t time.Time) int {
	week := int(t.Day()/7) + 1
	if t.Weekday() < time.Monday && (t.Day()-int(t.Weekday()))%7 != 0 {
		week++
	}
	return week
}

func parseWeekday(day string) time.Weekday {
	switch strings.ToLower(day) {
	case "monday", "mon", "mo":
		return time.Monday
	case "tuesday", "tue", "tu":
		return time.Tuesday
	case "wednesday", "wed", "we":
		return time.Wednesday
	case "thursday", "thu", "th":
		return time.Thursday
	case "friday", "fri", "fr":
		return time.Friday
	case "saturday", "sat", "sa":
		return time.Saturday
	case "sunday", "sun", "su":
		return time.Sunday
	default:
		return time.Monday
	}
}
