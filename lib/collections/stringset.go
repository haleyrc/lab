package collections

func StringSet() Set {
	return make(stringSet)
}

type stringSet map[string]struct{}

func (set stringSet) Add(s string) {
	if _, found := set[s]; !found {
		set[s] = struct{}{}
	}
}

func (set stringSet) AddN(ss ...string) {
	for _, s := range ss {
		set.Add(s)
	}
}

func (set stringSet) Slice() []string {
	ss := make([]string, 0, len(set))
	for k := range set {
		ss = append(ss, k)
	}
	return ss
}
