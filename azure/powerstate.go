package azure

import "strings"

type PowerState string

const (
	PowerStateDeallocated  PowerState = "PowerState/deallocated"
	PowerStateDeallocating PowerState = "PowerState/deallocating"
	PowerStateRunning      PowerState = "PowerState/running"
	PowerStateStarting     PowerState = "PowerState/starting"
	PowerStateStopped      PowerState = "PowerState/stopped"
	PowerStateStopping     PowerState = "PowerState/stopping"
	PowerStateUnknown      PowerState = "PowerState/unknown"
)

func (p PowerState) String() string {
	switch p {
	case PowerStateDeallocated:
		return "PowerState/deallocated"
	case PowerStateDeallocating:
		return "PowerState/deallocating"
	case PowerStateRunning:
		return "PowerState/running"
	case PowerStateStarting:
		return "PowerState/starting"
	case PowerStateStopped:
		return "PowerState/stopped"
	case PowerStateStopping:
		return "PowerState/stopping"
	case PowerStateUnknown:
		return "PowerState/unknown"
	default:
		return "unknown"
	}
}

func ParsePowerState(data string) PowerState {
	if strings.Contains(data, PowerStateDeallocated.String()) {
		return PowerStateDeallocated
	} else if strings.Contains(data, PowerStateDeallocating.String()) {
		return PowerStateDeallocating
	} else if strings.Contains(data, PowerStateRunning.String()) {
		return PowerStateRunning
	} else if strings.Contains(data, PowerStateStarting.String()) {
		return PowerStateStarting
	} else if strings.Contains(data, PowerStateStopped.String()) {
		return PowerStateStopped
	} else if strings.Contains(data, PowerStateStopping.String()) {
		return PowerStateStopping
	} else {
		return PowerStateUnknown
	}
}
