package form

import "net/http"

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

//nolint:errcheck
func Values(r *http.Request, key string) []string {
	if r.Form == nil {
		r.ParseMultipartForm(defaultMaxMemory)
	}
	if vs := r.Form[key]; len(vs) > 0 {
		return vs
	}
	return nil
}
