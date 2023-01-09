package execution

import (
	"fmt"
	"strings"
	"time"

	"github.com/alexallah/ethereum-healthmon/internal/common"
)

func StartUpdater(state *common.State, addr string, timeout int64, jwtPath string) {
	var secret []byte
	if jwtPath != "" {
		secret = loadJwt(jwtPath)
	}
	// add http: prefix if necessary
	if !strings.HasPrefix(addr, "http") {
		addr = fmt.Sprintf("http://%s", addr)
	}
	go update(state, addr, timeout, secret)
}

func update(state *common.State, addr string, timeout int64, secret []byte) {
	for {
		time.Sleep(time.Second)

		var token string
		if secret != nil {
			token = genToken(secret)
		}
		err := isReady(addr, token, timeout)

		if err != nil {
			state.Error(err)
		} else {
			state.SetHealthy()
		}
	}
}
