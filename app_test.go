package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"net/url"

	"github.com/stretchr/testify/assert"
)

func TestGetLogoutJSInjectResponseModifier(t *testing.T) {
	logoutUrl := "/cas/logout"
	responseModifier := getLogoutJSInjectResponseModifier(logoutUrl)
	b := bytes.NewBufferString("<html><body>stuff</body></html>")
	res := http.Response{
		Body:    ioutil.NopCloser(b),
		Header:  http.Header{},
		Request: &http.Request{URL: &url.URL{Path: "/"}},
	}
	res.Header.Set("Content-Type", "text/html")

	_ = responseModifier(&res)

	newBytes, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	injectedJSCode := fmt.Sprintf(injectedJSCodeTmpl, logoutUrl)
	assert.Equal(t, string(newBytes), fmt.Sprintf("<html><body>stuff%s</body></html>", injectedJSCode))
}

func TestIsJSInjectionRequiredRootPath(t *testing.T) {
	res := http.Response{
		Request: &http.Request{URL: &url.URL{Path: "/"}},
		Header:  http.Header{},
	}
	res.Header.Set("Content-Type", "text/html")

	injectRequired := isJSInjectionRequired(&res)

	assert.True(t, injectRequired)
}

func TestIsJSInjectionRequiredNexusPath(t *testing.T) {
	res := http.Response{
		Request: &http.Request{URL: &url.URL{Path: "/nexus/"}},
		Header:  http.Header{},
	}
	res.Header.Set("Content-Type", "text/html")

	injectRequired := isJSInjectionRequired(&res)

	assert.True(t, injectRequired)
}

func TestIsJSInjectionRequiredWrongPath(t *testing.T) {
	res := http.Response{
		Request: &http.Request{URL: &url.URL{Path: "/wrongPath"}},
		Header:  http.Header{},
	}
	res.Header.Set("Content-Type", "text/html")

	injectRequired := isJSInjectionRequired(&res)

	assert.False(t, injectRequired)
}

func TestIsJSInjectionRequiredMissingContentType(t *testing.T) {
	res := http.Response{
		Request: &http.Request{URL: &url.URL{Path: "/"}},
	}

	injectRequired := isJSInjectionRequired(&res)

	assert.False(t, injectRequired)
}

const succeed = "\u2713"
const failed = "\u2717"

func TestMapCesappLoglevelToGlogLoglevel(t *testing.T) {
	tables := []struct {
		name                string
		givenCesappLogLevel string
		expectedGlogLevel   int
	}{
		{"TestForError", "ERROR", 0},
		{"TestForWarn", "WARN", 1},
		{"TestForInfo", "INFO", 2},
		{"TestForDebug", "DEBUG", 5},
	}

	for _, table := range tables {
		receivedGlogLogLevel := mapCesappLoglevelToGlogLoglevel(table.givenCesappLogLevel)
		if table.expectedGlogLevel != receivedGlogLogLevel {
			log.Errorf("\t%s: should return logLevel: %v, got: %v", failed, table.expectedGlogLevel, receivedGlogLogLevel)
		} else {
			t.Logf("\t%s: logLevel should be accurate.", succeed)
		}
	}

}
