package bloom_filter

type Filter struct {
	keys map[string]bool
}

func NewFilter() *Filter {
	return &Filter{keys: map[string]bool{}}
}

func (f *Filter) Add(key string) {
	f.keys[key] = true
}

func (f *Filter) Get(key string) bool {
	if _, ok := f.keys[key]; ok {
		return true
	}
	return false
}
