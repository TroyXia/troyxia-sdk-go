package clientcmd

import (
	restclient "github.com/TroyXia/troyxia-sdk-go/rest"
)

// ClientConfig is used to make it easy to get an api server client.
type ClientConfig interface {
	// ClientConfig returns a complete client config
	ClientConfig() (*restclient.Config, error)
}
