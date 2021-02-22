package groups

import (
	"chores-suck/web/messages"
	"regexp"
	"strings"
)

func validateName(name string, msg *messages.CreateGroup) bool {
	valid := false
	regName := regexp.MustCompile(`^(?:[0-9a-zA-Z]+-)*[0-9a-zA-Z]+$`)

	if strings.TrimSpace(name) == "" {
		msg.Name = "Group name cannot be empty!"

	} else if !regName.MatchString(name) {
		msg.Name = "Group name must only consist of alphanumeric characters and hyphens and cannot start or end with a hyphen"
	} else {
		valid = true
	}

	return valid
}
