package testutil

import "os"

func SetEnvs(sets ...map[string]string) {
	for _, set := range sets {
		for k, v := range set {
			_ = os.Setenv(k, v)
		}
	}
}
