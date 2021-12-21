package main

import (
	"bytes"
	"flag"
	"github.com/cloudogu/carp"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"fmt"

	"github.com/cloudogu/go-health"
	logging "github.com/op/go-logging"
)

var Version = "x.y.z-dev"

var log = logging.MustGetLogger("nexus-carp")

func main() {
	log.Debug("Entering Method 'main'")
	defer func() {
		log.Debug("End of Function 'main'")
	}()

	flag.Parse()

	url := env("NEXUS_URL")
	log.Debugf("Variable: %s", url)
	username := env("NEXUS_USER")
	log.Debugf("Variable: %s", username)
	password := env("NEXUS_PASSWORD")
	log.Debugf("Variable: %s", password)
	cesAdminGroup := env("CES_ADMIN_GROUP")
	log.Debugf("Variable: %s", cesAdminGroup)
	timeout := getTimeoutOrDefault("HTTP_REQUEST_TIMEOUT", 30)

	configuration, err := carp.InitializeAndReadConfiguration()
	log.Debugf("Variable: %s", configuration)
	if err != nil {
		log.Debugf("Error: %s", err.Error())
		log.Fatal("failed to read configuration:", err)
	}

	log = logging.MustGetLogger("nexus-carp")

	log.Info("wait until nexus is ready")
	err = waitUntilNexusBecomesReady()
	if err != nil {
		log.Fatal("nexus did not become ready:", err)
	}

	log.Infof("start nexus-carp %s", Version)

	userReplicator := NewUserReplicator(url, username, password, timeout)
	log.Debugf("Variable: %s", userReplicator)
	err = userReplicator.CreateScript(cesAdminGroup)
	if err != nil {
		log.Debugf("Error: %s", err.Error())
		log.Fatal("failed to create user replication script:", err)
	}

	configuration.UserReplicator = userReplicator.Replicate
	log.Debugf("Variable: %s", configuration.UserReplicator)
	configuration.ResponseModifier = getLogoutJSInjectResponseModifier(configuration.CasUrl + "/logout")
	log.Debugf("Variable: %s", configuration.ResponseModifier)

	server, err := carp.NewServer(configuration)
	log.Debugf("Variable: %s", server)
	if err != nil {
		log.Debugf("Error: %s", err.Error())
		panic(err)
	}

	server.ListenAndServe()
}

func getTimeoutOrDefault(variableName string, defaultValue int) int {
	log.Debug("Entering Method 'getTimeoutOrDefault'")
	defer func() {
		log.Debug("End of Function 'getTimeoutOrDefault'")
	}()

	log.Debugf("Param '%s'", variableName)
	log.Debugf("Param '%s'", defaultValue)
	value := os.Getenv(variableName)
	log.Debugf("Variable: %s", value)
	if value == "" {
		log.Debugf("Condition true: 'value == \"\"'")
		return defaultValue
	}
	timeoutFromEnv, err := strconv.Atoi(value)
	log.Debugf("Variable: %s", timeoutFromEnv)
	if err != nil {
		log.Debugf("Error: %s", err.Error())
		return defaultValue
	}
	return timeoutFromEnv
}

func env(key string) string {
	log.Debug("Entering Method 'env'")
	defer func() {
		log.Debug("End of Function 'env'")
	}()

	log.Debugf("Param '%s'", key)
	value := os.Getenv(key)
	if value == "" {
		log.Debugf("Condition true: 'value == \"\"'")
		log.Fatalf("environment variable %s is not set", key)
	}
	return value
}

func waitUntilNexusBecomesReady() error {
	log.Debug("Entering Method 'waitUntilNexusBecomesReady'")
	defer func() {
		log.Debug("End of Function 'waitUntilNexusBecomesReady'")
	}()

	checker := health.NewTCPHealthCheckBuilder(8081).Build()
	log.Debugf("Variable: %s", checker)

	watcher := health.NewWatcher()
	log.Debugf("Variable: %s", watcher)
	watcher.RecheckLimit = 300
	log.Debugf("Variable: %s", watcher.RecheckLimit)
	watcher.ResultListener = func(counter int, err error) {
		log.Infof("nexus health check number %d failed, still waiting until nexus becomes ready", counter)
	}
	err := watcher.WaitUntilHealthy(checker)
	if err != nil {
		log.Debugf("Error: %s", err.Error())
		return err
	}
	return nil
}

const injectedJSCodeTmpl = "<script>" +
	"var timer = setInterval(function () {" +
	"var signoutElements = document.querySelectorAll(\"span[id^='nx-header-signout-']\");" +
	"if (signoutElements.length == 0) { return; }" +
	"signoutElements[0].addEventListener('click', function () { " +
	"window.location.href = '%s';" +
	"return true;" +
	"});" +
	"clearInterval(timer);" +
	"}, 500);" +
	"</script>"

func getLogoutJSInjectResponseModifier(logoutUrl string) func(resp *http.Response) error {
	log.Debug("Entering Method 'getLogoutJSInjectResponseModifier'")
	defer func() {
		log.Debug("End of Function 'getLogoutJSInjectResponseModifier'")
	}()

	log.Debugf("Param '%s'", logoutUrl)
	return func(resp *http.Response) error {
		log.Debug("Entering Function 'func(resp *http.Response) error'")
		if isJSInjectionRequired(resp) {
			log.Debugf("Condition true: 'isJSInjectionRequired(resp)'")
			b, err := ioutil.ReadAll(resp.Body)
			log.Debugf("Variable: %s", b)
			if err != nil {
				log.Debugf("Error: %s", err.Error())
				return err
			}
			injectedJSCode := fmt.Sprintf(injectedJSCodeTmpl, logoutUrl)
			log.Debugf("Variable: %s", injectedJSCode)
			b = bytes.Replace(b, []byte("</body>"), []byte(injectedJSCode+"</body>"), 1)
			log.Debugf("Variable: %s", b)
			body := ioutil.NopCloser(bytes.NewReader(b))
			log.Debugf("Variable: %s", body)
			resp.Body = body
			log.Debugf("Variable: %s", resp.Body)
			resp.ContentLength = int64(len(b))
			log.Debugf("Variable: %s", resp.ContentLength)
			resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
			log.Debugf("Variable: %s", resp.Header)
		}

		return nil
	}
}

func isJSInjectionRequired(resp *http.Response) bool {
	log.Debug("Entering Method 'isJSInjectionRequired'")
	defer func() {
		log.Debug("End of Function 'isJSInjectionRequired'")
	}()

	log.Debugf("Param '%s'", resp)
	injectPaths := []string{"/nexus/", "/"}
	log.Debugf("Variable: %s", injectPaths)
	b := resp.Header.Get("Content-Type") == "text/html" && contains(injectPaths, resp.Request.URL.Path)
	log.Debugf("Variable: %s", b)
	return b
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
