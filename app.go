package main

import (
	"flag"
	"os"

	"github.com/cloudogu/carp"
	"github.com/cloudogu/go-health"
	"github.com/golang/glog"
)

var Version = "x.y.z-dev"

func main() {
	flag.Parse()

	url := env("NEXUS_URL")
	username := env("NEXUS_USER")
	password := env("NEXUS_PASSWORD")
	cesAdminGroup := env("CES_ADMIN_GROUP")

	configuration, err := carp.ReadConfiguration()
	if err != nil {
		glog.Fatal("failed to read configuration:", err)
	}

	glog.Infof("wait until nexus is ready")
	err = waitUntilNexusBecomesReady(url, username, password)
	if err != nil {
		glog.Fatal("nexus does not become ready:", err)
	}

	glog.Infof("start nexus-carp %s", Version)

	userReplicator := NewUserReplicator(url, username, password)
	err = userReplicator.CreateScript(cesAdminGroup)
	if err != nil {
		glog.Fatal("failed to create user replication script:", err)
	}

	configuration.UserReplicator = userReplicator.Replicate

	server, err := carp.NewServer(configuration, true)
	if err != nil {
		panic(err)
	}

	server.ListenAndServe()
}

func env(key string) string {
	value := os.Getenv(key)
	if value == "" {
		glog.Fatalf("environment variable %s is not set", key)
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
		glog.Infof("nexus health check number %v failed, still waiting until nexus becomes ready", counter)
	}
	err := watcher.WaitUntilHealthy(checker)
	if err != nil {
		return err
	}
	return nil
}
