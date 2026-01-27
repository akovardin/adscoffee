package health

import (
	"encoding/json"
	"net/http"
)

func (h *Health) HandlerExternal() http.HandlerFunc {
	return h.typeHandler(ComponentKindExternal)
}

func (h *Health) Handler() http.HandlerFunc {
	return h.typeHandler(ComponentKindApp | ComponentKindLocal)
}

func (h *Health) typeHandler(
	kind ComponentKind,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		h.Check(r.Context(), kind)
		state := NewState(h, kind)
		if state.Status != StatusUp {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		st, err := json.Marshal(state)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			st = []byte(err.Error())
		}

		if _, err := w.Write(st); err != nil {
			return
		}
	}
}
