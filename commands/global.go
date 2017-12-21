package commands

import (
	"github.com/analogrepublic/kongctl/kong"
)

// This file defines methods & variables that should be
// globally available to all commands.

var kongApi *kong.Kong

func SetKongApi(k *kong.Kong) {
	kongApi = k
}
