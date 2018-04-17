package main

import (
	"flag"
	"os"

	"github.com/cloudogu/carp"
	"github.com/golang/glog"
)

var Version = "x.y.z-dev"

func main() {
	flag.Parse()

	url := env("NEXUS_URL")
	username := env("NEXUS_USER")
	password := env("NEXUS_PASSWORD")

	configuration, err := carp.ReadConfiguration()
	if err != nil {
		glog.Fatal("failed to read configuration:", err)
	}

	glog.Infof("start nexus-carp %s", Version)

	userReplicator := NewUserReplicator(url, username, password)
	err = userReplicator.CreateScript()
	if err != nil {
		glog.Fatal("failed to create user replication script:", err)
	}

	configuration.UserReplicator = userReplicator.Replicate

	server, err := carp.NewServer(configuration)
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
