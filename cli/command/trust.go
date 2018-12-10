package command

import (
	"os"
	"strconv"

	"github.com/spf13/pflag"
)

// contentTrustEnabled determines whether Docker Content Trust is enabled in the environment.
func contentTrustEnabled() bool {
	if e := os.Getenv("DOCKER_CONTENT_TRUST"); e != "" {
		if t, err := strconv.ParseBool(e); t || err != nil {
			// treat any other value as true
			return true
		}
	}
	return false
}

// AddTrustVerificationFlags adds content trust flags to the provided flagset
func AddTrustVerificationFlags(fs *pflag.FlagSet, v *bool, trusted bool) {
	fs.BoolVar(v, "disable-content-trust", !trusted, "Skip image verification")
}

// AddTrustSigningFlags adds "signing" flags to the provided flagset
func AddTrustSigningFlags(fs *pflag.FlagSet, v *bool, trusted bool) {
	fs.BoolVar(v, "disable-content-trust", !trusted, "Skip image signing")
}
