//go:generate go run scripts/generate.go replicator_scripts scripts
package main

import (
	"github.com/cloudogu/carp"
	"github.com/cloudogu/nexus-scripting/manager"
	"github.com/golang/glog"
	"github.com/pkg/errors"
  "strings"
)

const scriptName = "carp-user-replication"

func NewUserReplicator(url string, username string, password string, timeout int) *UserReplicator {
	manager := manager.New(url, username, password)
	manager.WithTimeout(timeout)
	return &UserReplicator{
		manager: manager,
	}
}

type UserReplicator struct {
	manager *manager.Manager
	script  *manager.Script
}

func (replicator *UserReplicator) CreateScript(cesAdminGroup string) error {
  userReplicationScript := CARP_USER_REPLICATION
  userReplicationScript = strings.Replace(userReplicationScript, "cesAdminGroup", cesAdminGroup, -1)
	script, err := replicator.manager.Create(scriptName, userReplicationScript)
	if err != nil {
		return errors.Wrap(err, "failed to create user replication script")
	}

	replicator.script = script
	return nil
}

func (replicator *UserReplicator) Replicate(username string, attributes carp.UserAttibutes) error {
	nexusUser := createNexusCarpUser(attributes)
	out, err := replicator.script.ExecuteWithJSONPayload(nexusUser)
	if err != nil {
		return errors.Wrapf(err, "user replication script failed for user %s", nexusUser.Username)
	}

	if out != "" && glog.V(2) {
		glog.Infof("user replication script returned %s for user %s", out, nexusUser.Username)
	}
	return nil
}

func createNexusCarpUser(attributes carp.UserAttibutes) *NexusCarpUser {
	return &NexusCarpUser{
		Username:  firstOrEmpty(attributes["username"]),
		FirstName: firstOrEmpty(attributes["givenName"]),
		LastName:  firstOrEmpty(attributes["surname"]),
		Email:     firstOrEmpty(attributes["mail"]),
		Groups:    attributes["groups"],
	}
}

func firstOrEmpty(values []string) string {
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

type NexusCarpUser struct {
	Username  string
	FirstName string
	LastName  string
	Email     string
	Groups    []string
}
