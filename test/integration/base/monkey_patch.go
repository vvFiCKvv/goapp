package integration

import (
	"goapp/internal/pkg/httpsrv"
	"goapp/internal/pkg/strgen"
	"time"
)

func monkeyPatchOrigin() func() {
	originalValidOrigins := httpsrv.ValidOrigins
	httpsrv.ValidOrigins = []string{"http://mytestsite.io", "http://localhost:8080"}

	return func() {
		httpsrv.ValidOrigins = originalValidOrigins
	}
}

func monkeyPatchStrGenTime() func() {
	StrGenTimeOrigins := strgen.StrGenTime

	strgen.StrGenTime = time.Millisecond * 100
	return func() {
		strgen.StrGenTime = StrGenTimeOrigins
	}
}

func monkeyPatches() []func() {
	return []func(){monkeyPatchOrigin(), monkeyPatchStrGenTime()}
}
