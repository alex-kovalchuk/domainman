package sort

type DomainLengthSort []string

func (s DomainLengthSort) Len() int {
	return len(s)
}

func (s DomainLengthSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s DomainLengthSort) Less(i, j int) bool {
	if len(s[i]) == len(s[j]) {
		return s[i] < s[j]
	}
	return len(s[i]) < len(s[j])
}
