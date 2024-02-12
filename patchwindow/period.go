package patchwindow

import (
	"encoding/json"
	"errors"
)

type Period int

const (
	MonthlyPatchPeriod Period = 0

	InvalidPatchPeriod Period = -1
)

func ParsePatchPeriod(period string) Period {
	switch period {
	case "monthly":
		return MonthlyPatchPeriod
	default:
		return InvalidPatchPeriod
	}
}

func (p *Period) UnmarshalJSON(b []byte) error {
	var period string
	if err := json.Unmarshal(b, &period); err != nil {
		return err
	}
	switch period {
	case "monthly":
		*p = MonthlyPatchPeriod
	default:
		*p = InvalidPatchPeriod
		return errors.New("invalid period type")
	}
	return nil
}
