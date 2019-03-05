package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLogoutJSInjectResponseModifier(t *testing.T) {
	logoutUrl := "/cas/logout"
	responseModifier := getLogoutJSInjectResponseModifier(logoutUrl)
	b := bytes.NewBufferString("<html><body>stuff</body></html>")
	res := http.Response{
		Body:   ioutil.NopCloser(b),
		Header: http.Header{},
	}
	res.Header.Set("Content-Type", "text/html")

	_ = responseModifier(&res)

	newBytes, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	injectedJSCode := fmt.Sprintf(injectedJSCodeTmpl, logoutUrl)
	assert.Equal(t, string(newBytes), fmt.Sprintf("<html><body>stuff%s</body></html>", injectedJSCode))
}

func TestGetLogoutJSInjectResponseModifierWithWrongContent(t *testing.T) {

	responseModifier := getLogoutJSInjectResponseModifier("test")
	originalBodyContent := "<html><body>stuff</body></html>"
	b := bytes.NewBufferString(originalBodyContent)
	res := http.Response{
		Body: ioutil.NopCloser(b),
	}

	err := responseModifier(&res)

	assert.Nil(t, err)
	newBytes, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Len(t, newBytes, len(originalBodyContent))
}
