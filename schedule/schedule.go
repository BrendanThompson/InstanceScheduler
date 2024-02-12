package schedule

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func NewSchedule(data []byte) (*Schedule, error) {
	var schedule Schedule

	err := json.Unmarshal(data, &schedule)
	if err != nil {
		return nil, err
	}

	return &schedule, nil
}

type Schedule struct {
	Default   string              `json:"default"`
	Overrides map[string][]string `json:"overrides"`
}

func (s *Schedule) Validate() bool {

	start, end, err := ParseWindow(s.Default)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Failed to parse the time window for 'Default' schedule.")
		return false
	}

	if end.After(start) {
		log.Debug().Msg("Valid time window for 'Default' schedule.")
		return true
	}

	log.Error().Msg("Invalid time window for 'Default' schedule.")

	return false
}

func (s *Schedule) ValidateOverrides() bool {
	if !s.HasOverrides() {
		return true
	}

	for day, windows := range s.Overrides {
		_, err := time.Parse("Monday", day)
		if err != nil {
			log.Error().Msgf("Invalid weekday provided to 'Override' schedule: '%s'", day)
			return false
		}

		for _, window := range windows {
			if window == "-" {
				return true
			}

			start, end, err := ParseWindow(window)
			if err != nil {
				log.Error().Stack().Err(err).Msg("Failed to parse the time window for 'Override' schedule.")
				return false
			}

			if end.After(start) {
				log.Debug().Msg("Valid time window for 'Override' schedule.")
				return true
			}
		}
	}

	log.Error().Msg("Invalid time window for 'Override' schedule.")

	return false
}

func (s *Schedule) ShouldShutdown() bool {
	var result bool
	var now time.Time = time.Now().Local()
	shouldOverride, overrideKey := s.UseOverride(now.Weekday())

	if shouldOverride {
		log.Debug().Msg("Using override schdeule")

		for _, t := range s.Overrides[overrideKey] {
			log.Info().Msg(t)

			if t == "-" {
				return true
			}

			start, end, err := ParseWindow(t)
			if err != nil {
				return false
			}

			if (now.After(start) || now.Equal(start)) && now.Before(end) {
				result = false
			} else {
				result = true
				return result
			}
		}
	} else {
		log.Debug().Msg("Using default schedule")

		start, end, err := ParseWindow(s.Default)
		if err != nil {
			return false
		}

		if (now.After(start) || now.Equal(start)) && now.Before(end) {
			result = false
		} else {
			result = true
			return result
		}

	}

	return result
}

func (s *Schedule) UseOverride(weekday time.Weekday) (bool, string) {
	for key := range s.Overrides {
		if strings.ToLower(key) == strings.ToLower(weekday.String()) {
			return true, key
		} else {
			return false, ""
		}
	}

	return false, ""
}

func (s *Schedule) HasOverrides() bool {
	if len(s.Overrides) > 0 {
		return true
	} else {
		return false
	}
}

// TODO: this should probably return an error
func (s *Schedule) IsWithinPatchWindow(start, end time.Time, isToday bool) bool {
	if !isToday {
		log.Debug().Msg("Not today patch window")
		return false
	}

	var now time.Time = time.Now().Local()
	shouldOverride, overrideKey := s.UseOverride(now.Weekday())

	if shouldOverride {
		log.Debug().Msg("Override Schedule - Patch Window Check")
		for _, t := range s.Overrides[overrideKey] {
			scheduleStart, scheduleEnd, err := ParseWindow(t)
			if err != nil {
				log.Error().Stack().Err(err).Msg("Failed parsing time window for `Override` schedule.")
				return false
			}

			if scheduleStart.Before(start) && scheduleEnd.After(end) {
				log.Debug().Msg("I return true")
				return true
			}
		}
	} else {
		log.Debug().Msg("Default Schedule - Patch Window Check")

		scheduleStart, scheduleEnd, err := ParseWindow(s.Default)
		if err != nil {
			log.Error().Stack().Err(err).Msg("Failed parsing time window for `Default` schedule.")
			return false
		}

		if start.After(scheduleEnd) || end.Before(scheduleStart) {
			log.Debug().Msg("Shutdown is within patch window")
			return true
		}

	}

	return false
}

func ParseWindow(data string) (time.Time, time.Time, error) {
	var now time.Time = time.Now().Local()
	var timeWindows []string = strings.Split(data, "-")

	start, err := time.Parse("15:04", timeWindows[0])
	if err != nil {
		return time.Now(), time.Now(), err
	}

	end, err := time.Parse("15:04", timeWindows[1])
	if err != nil {
		return time.Now(), time.Now(), err
	}

	return time.Date(now.Year(), now.Month(), now.Day(), start.Hour(), start.Minute(), 0, 0, now.Location()),
		time.Date(now.Year(), now.Month(), now.Day(), end.Hour(), end.Minute(), 0, 0, now.Location()),
		nil
}
