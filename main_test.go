package main

import (
	"os"
	"testing"
)

const debug_key string = "DEBUG"

func TestConverMessageJson(t *testing.T) {
	const msg string = `{"region":"ap-east-1","deploymentId":"d-X9Y85B59X","instanceId":"i-007d56b93274560f6"}`

	if v := convertMessage(msg); v != "" {
		t.Log("PASS")
	} else {
		t.Error("Fail, convert message error")
	}
}

func TestSetDebugEnv(t *testing.T) {
	os.Setenv(debug_key, "true")

	if v := getBoolEnv(debug_key); v {
		t.Log("PASS")
	} else {
		t.Error("Fail, debug not true")
	}
}

func TestNullDebugEnv(t *testing.T) {
	os.Unsetenv(debug_key)

	v := getBoolEnv(debug_key)

	if v == false {
		t.Log("PASS")
	} else {
		t.Error("Fail, debug env not false")
	}
}

func TestNullSendInfo(t *testing.T) {
	info := LineInfo{}

	err := send(info)

	if err != nil {
		t.Log("PASS")
	} else {
		t.Error("Fail, error not found")
	}
}

func TestSend(t *testing.T) {
	info := LineInfo{
		Token:   "1234567890",
		Message: "Test",
	}

	err := send(info)

	if err == nil {
		t.Log("PASS")
	} else {
		t.Error("Fail, send error")
	}
}