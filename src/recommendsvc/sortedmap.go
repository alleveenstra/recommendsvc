package recommendsvc

type SortedMap struct {
	M map[int]float64
	S []int
}

func NewSortedMap(m map[int]float64) *SortedMap {
	sm := new(SortedMap)
	sm.M = m
	sm.S = make([]int, len(m))
	i := 0
	for key, _ := range m {
		sm.S[i] = key
		i++
	}
	return sm
}

func (sm *SortedMap) Len() int {
	return len(sm.M)
}

func (sm *SortedMap) Less(i, j int) bool {
	return sm.M[sm.S[i]] < sm.M[sm.S[j]]
}

func (sm *SortedMap) Swap(i, j int) {
	sm.S[i], sm.S[j] = sm.S[j], sm.S[i]
}
