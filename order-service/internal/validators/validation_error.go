package validators

type ValidationErrors map[string]string

func (v ValidationErrors) Merge(other ValidationErrors) {
	for k, val := range other {
		v[k] = val
	}
}
