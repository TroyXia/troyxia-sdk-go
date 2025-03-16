package rest

import (
	"time"

	"github.com/marmotedu/component-base/pkg/runtime"
	"github.com/marmotedu/component-base/pkg/scheme"
)

// Config holds the common attributes that can be passed to a IAM client on
// initialization.
type Config struct {
	Host    string
	APIPath string
	ContentConfig

	// Server requires Basic authentication
	Username string
	Password string

	SecretID  string
	SecretKey string

	// Server requires Bearer authentication. This client will not attempt to use
	// refresh tokens for an OAuth2 flow.
	// TODO: demonstrate an OAuth2 compatible client.
	BearerToken string

	// Path to a file containing a BearerToken.
	// If set, the contents are periodically read.
	// The last successfully read value takes precedence over BearerToken.
	BearerTokenFile string

	// TLSClientConfig contains settings to enable transport layer security
	TLSClientConfig

	// UserAgent is an optional field that specifies the caller of this request.
	UserAgent string
	// The maximum length of time to wait before giving up on a server request. A value of zero means no timeout.
	Timeout       time.Duration
	MaxRetries    int
	RetryInterval time.Duration
}

// ContentConfig defines config for content.
type ContentConfig struct {
	ServiceName        string
	AcceptContentTypes string
	ContentType        string
	GroupVersion       *scheme.GroupVersion
	Negotiator         runtime.ClientNegotiator
}

// TLSClientConfig contains settings to enable transport layer security.
type TLSClientConfig struct {
	// Server should be accessed without verifying the TLS certificate. For testing only.
	Insecure bool
	// ServerName is passed to the server for SNI and is used in the client to check server
	// ceritificates against. If ServerName is empty, the hostname used to contact the
	// server is used.
	ServerName string

	// Server requires TLS client certificate authentication
	CertFile string
	// Server requires TLS client certificate authentication
	KeyFile string
	// Trusted root certificates for server
	CAFile string

	// CertData holds PEM-encoded bytes (typically read from a client certificate file).
	// CertData takes precedence over CertFile
	CertData []byte
	// KeyData holds PEM-encoded bytes (typically read from a client certificate key file).
	// KeyData takes precedence over KeyFile
	KeyData []byte
	// CAData holds PEM-encoded bytes (typically read from a root certificates bundle).
	// CAData takes precedence over CAFile
	CAData []byte

	// NextProtos is a list of supported application level protocols, in order of preference.
	// Used to populate tls.Config.NextProtos.
	// To indicate to the server http/1.1 is preferred over http/2, set to ["http/1.1", "h2"] (though the server is free
	// to ignore that preference).
	// To use only http/1.1, set to ["http/1.1"].
	NextProtos []string
}
