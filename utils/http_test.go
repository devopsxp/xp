package utils

import (
	"testing"
)

func TestHttpGet(t *testing.T) {
	url := "https://www.baidu.com"
	_, err := Get(url)
	if err != nil {
		t.Errorf("Get() return an error: %v", err)
	}
}
