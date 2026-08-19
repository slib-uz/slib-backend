package utils

func NormalizeTin(s string) string {
	out := make([]byte, 0, len(s))
	for _, r := range s {
		if r >= '0' && r <= '9' {
			out = append(out, byte(r))
		}
	}
	return string(out)
}

func IsValidOrganizationTin(s string) bool {
	return len(NormalizeTin(s)) == 9
}
