package template

import "strings"

// ParseEnvVars returns a map from a list of env vars in the form
// key=value.
func ParseEnvVars(envVars []string) map[string]string {
	vars := map[string]string{}
	for _, str := range envVars {
		parts := strings.Split(str, "=")
		if len(parts) > 1 {
			vars[parts[0]] = parts[1]
		}
	}
	return vars
}
