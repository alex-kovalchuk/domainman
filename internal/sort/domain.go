package sort

type Domain []string

func (s Domain) Len() int {
	return len(s)
}

func (s Domain) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Domain) Less(i, j int) bool {
	if len(s[i]) == len(s[j]) {
		return s[i] < s[j]
	}
	return len(s[i]) < len(s[j])
}
