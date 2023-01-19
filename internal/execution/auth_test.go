package execution

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"os"
	"testing"
)

func newJWT() []byte {
	buf := make([]byte, 32)
	rand.Read(buf)
	return buf
}

func capturePanic(f func()) (ret any) {
	defer func() {
		ret = recover()
	}()
	f()
	return ret
}

func checkPanic(t *testing.T, f func(), expected any) {
	result := capturePanic(f)
	if result != expected {
		t.Fatalf("expected panic with \"%v\", got %+v", expected, result)
	}
}

func writeJwt(filepath string, value string) {
	err := os.WriteFile(filepath, []byte(value), 0600)
	if err != nil {
		panic("can not write temporary jwt file")
	}
}

func TestJwt(t *testing.T) {
	// generate a jwt file
	dir := t.TempDir()
	keypath := dir + "/jwt.hex"
	jwt := newJWT()
	jwtHex := fmt.Sprintf("%x", jwt)

	// read
	// wrong file
	if capturePanic(func() {
		readJwt(dir + "/nofile")
	}) == nil {
		t.Fatal("expected to panic on a wrong file")
	}
	// good
	writeJwt(keypath, jwtHex)
	readJwtValue := readJwt(keypath)
	if readJwtValue != jwtHex {
		t.Fatalf("jwt hex values don't match %s != %s", jwtHex, readJwtValue)
	}

	// load
	writeJwt(keypath, jwtHex)
	if !bytes.Equal(loadJwt(keypath), jwt) {
		t.Fatal("jwt byte values are different")
	}
	// incorrect data
	incorrectkeypath := dir + "/broken.jwt.hex"
	writeJwt(incorrectkeypath, "wrong data !@#")
	if capturePanic(func() { loadJwt(incorrectkeypath) }) == nil {
		t.Fatal("expected to panic with wrong data")
	}

	// gen token
	// all good
	if len(genToken(jwt)) == 0 {
		t.Fatal("empty token")
	}
	// nill token
	checkPanic(t, func() {
		genToken(nil)
	}, "no secret")
}
