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
