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

var log *logging.Logger

func main() {
	flag.Parse()

	url := env("NEXUS_URL")
	username := env("NEXUS_USER")
	password := env("NEXUS_PASSWORD")
	cesAdminGroup := env("CES_ADMIN_GROUP")
	timeout := getTimeoutOrDefault("HTTP_REQUEST_TIMEOUT", 30)

	configuration, err := carp.InitializeAndReadConfiguration()
	if err != nil {
		log.Fatal("failed to read configuration:", err)
	}

	log = logging.MustGetLogger("nexus-carp")

	log.Infof("wait until nexus is ready")
	err = waitUntilNexusBecomesReady(url, username, password)
	if err != nil {
		log.Fatal("nexus does not become ready:", err)
	}

	log.Infof("start nexus-carp %s", Version)

	userReplicator := NewUserReplicator(url, username, password, timeout)
	err = userReplicator.CreateScript(cesAdminGroup)
	if err != nil {
		log.Fatal("failed to create user replication script:", err)
	}

	configuration.UserReplicator = userReplicator.Replicate
	configuration.ResponseModifier = getLogoutJSInjectResponseModifier(configuration.CasUrl + "/logout")

	server, err := carp.NewServer(configuration)
	if err != nil {
		panic(err)
	}

	server.ListenAndServe()
}

func getTimeoutOrDefault(variableName string, defaultValue int) int {
	value := os.Getenv(variableName)
	if value == "" {
		return defaultValue
	}
	timeoutFromEnv, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return timeoutFromEnv
}

func env(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("environment variable %s is not set", key)
	}
	return value
}

func waitUntilNexusBecomesReady(url string, username string, password string) error {
	checker := health.NewHTTPHealthCheckBuilder(url+"/service/metrics/healthcheck").
		WithMethod("GET").
		WithBasicAuth(username, password).
		Build()

	watcher := health.NewWatcher()
	watcher.RecheckLimit = 300
	watcher.ResultListener = func(counter int, err error) {
		log.Infof("nexus health check number %v failed, still waiting until nexus becomes ready", counter)
	}
	err := watcher.WaitUntilHealthy(checker)
	if err != nil {
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
	return func(resp *http.Response) error {
		if isJSInjectionRequired(resp) {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			injectedJSCode := fmt.Sprintf(injectedJSCodeTmpl, logoutUrl)
			b = bytes.Replace(b, []byte("</body>"), []byte(injectedJSCode+"</body>"), 1)
			body := ioutil.NopCloser(bytes.NewReader(b))
			resp.Body = body
			resp.ContentLength = int64(len(b))
			resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
		}
		return nil
	}
}

func isJSInjectionRequired(resp *http.Response) bool {
	injectPaths := []string{"/nexus/", "/"}
	return resp.Header.Get("Content-Type") == "text/html" && contains(injectPaths, resp.Request.URL.Path)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
