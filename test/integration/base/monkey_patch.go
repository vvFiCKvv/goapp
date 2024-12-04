package integration

import "goapp/internal/pkg/httpsrv"

func monkeyPatchOrigin() func() {
	originalValidOrigins := httpsrv.ValidOrigins
	httpsrv.ValidOrigins = []string{"http://mytestsite.io", "http://localhost:8080"}

	return func() {
		httpsrv.ValidOrigins = originalValidOrigins
	}
}

func monkeyPatches() []func() {
	return []func(){monkeyPatchOrigin()}
}
