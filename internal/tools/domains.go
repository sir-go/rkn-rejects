package tools

// GetUpperDomains returns a list of all upper domains of the given hostname
func GetUpperDomains(d string) (res []string) {
	var dotIndexes []int
	for idx, r := range d {
		if r == '.' {
			dotIndexes = append(dotIndexes, idx)
		}
	}
	for i := len(dotIndexes) - 2; i > -1; i-- {
		res = append(res, d[dotIndexes[i]+1:])
	}
	return append(res, d)
}
