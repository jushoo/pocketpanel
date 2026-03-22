package sync

import "regexp"

var semverRegex = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

func IsValidSemver(version string) bool {
	return semverRegex.MatchString(version)
}
