package integration

import (
	"goapp/internal/pkg/httpsrv"
	"goapp/internal/pkg/strgen"
	"log"
	"net"
	"strconv"
	"time"
)

func GetFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

func monkeyPatchOrigin() func() {
	originalValidOrigins := httpsrv.ValidOrigins
	httpsrv.ValidOrigins = []string{"http://mytestsite.io", "http://localhost:8080", "http://localhost:" + strconv.Itoa(httpsrv.Port)}

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

func monkeyPatchServerPort() func() {
	PortOrigins := httpsrv.Port

	freePort, err := GetFreePort()
	if err != nil {
		log.Fatalf("fatal: %+v\n", err)
	}
	httpsrv.Port = freePort
	return func() {
		httpsrv.Port = PortOrigins
	}
}

func monkeyPatches() []func() {
	return []func(){monkeyPatchServerPort(), monkeyPatchOrigin(), monkeyPatchStrGenTime()}
}
