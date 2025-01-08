package utils

import "os"

func EnvOrDefault(key, def string) string {
	e, ok := os.LookupEnv(key)
	if !ok || e == "" {
		return def
	}
	return e
}
