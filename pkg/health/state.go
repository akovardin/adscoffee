package health

import (
	"math"
)

type Detail string

const (
	HostnameDetail Detail = "hostname"
	VersionDetail  Detail = "version"
	TimeDetail     Detail = "time"
	MessageDetail  Detail = "message"
)

type ComponentState struct {
	Status  Status         `json:"status"`
	Details map[Detail]any `json:"details,omitempty"`
}

type State struct {
	Status Status `json:"status"`

	Components map[string]ComponentState `json:"components"`
}

func NewState(h *Health, requestedKind ComponentKind) State {
	state := State{
		Status:     StatusUp,
		Components: make(map[string]ComponentState, len(h.components)),
	}

	h.Iter(requestedKind, func(c *Component) {
		var status Status

		switch c.CheckErr {
		case nil:
			status = StatusUp
		default:
			status = StatusDown
		}

		details := make(map[Detail]any)

		if c.Kind == ComponentKindApp {
			cfg := h.Config()

			details[HostnameDetail] = cfg.Hostname
			details[VersionDetail] = cfg.Version
		}

		elapsedSeconds := c.CheckDuration.Seconds()

		// round to 4 decimal placess
		details[TimeDetail] = math.Round(elapsedSeconds*10000) / 10000

		if c.CheckErr != nil {
			details[MessageDetail] = c.CheckErr.Error()
		}

		componentState := ComponentState{
			Status:  status,
			Details: details,
		}

		if componentState.Status == StatusDown {
			state.Status = StatusDown
		}

		state.Components[c.Name] = componentState
	})

	return state
}
