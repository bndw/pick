package utils

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	versionSeparator = "."
	versionLength    = 3
)

func ParseVersion(v string) ([]int, error) {
	parts := make([]int, versionLength)
	split := strings.SplitN(v, versionSeparator, versionLength)
	if len(split) != versionLength {
		return nil, fmt.Errorf("could not parse version '%s'", v)
	}
	// Remove 'v' from major version part
	split[0] = strings.Replace(split[0], "v", "", 1)
	for i, part := range split {
		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		parts[i] = num
	}
	return parts, nil
}
