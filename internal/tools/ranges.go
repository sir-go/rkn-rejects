package tools

import (
	"strconv"
	"strings"
)

// ParseRange unwraps numbers range contained string to the given slice
// 100-103 -> [100 101 102 103]
func ParseRange(s string, a *[]uint16) error {
	var (
		v0, v1 uint64
		err    error
	)
	if !strings.ContainsRune(s, '-') {
		v0, err = strconv.ParseUint(s, 10, 16)
		if err != nil {
			return err
		}
		*a = []uint16{uint16(v0)}
		return nil
	}

	p := strings.Split(s, "-")
	if v0, err = strconv.ParseUint(p[0], 10, 16); err != nil {
		return err
	}
	if v1, err = strconv.ParseUint(p[1], 10, 16); err != nil {
		return err
	}
	for v0 <= v1 {
		*a = append(*a, uint16(v0))
		v0++
	}
	return nil
}
