package clientcmd

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	utilerrors "github.com/TroyXia/errors"
)

var (
	// ErrNoContext defines no context chosen error.
	ErrNoContext = errors.New("no context chosen")

	// ErrEmptyConfig defines no configuration has been provided error.
	ErrEmptyConfig = NewEmptyConfigError(
		"no configuration has been provided, try setting IAM_SERVER_ADDRESS environment variable",
	)

	// ErrEmptyServer defines a no server defined error.
	ErrEmptyServer = errors.New("server has no server defined")
)

// errConfigurationInvalid is a set of errors indicating the configuration is invalid.
type errConfigurationInvalid []error

// Error implements the error interface.
func (e errConfigurationInvalid) Error() string {
	return fmt.Sprintf("invalid configuration: %v", utilerrors.NewAggregate(e).Error())
}

// NewEmptyConfigError returns an error wrapping the given message which IsEmptyConfig()
// will recognize as an empty config error.
func NewEmptyConfigError(message string) error {
	return &errEmptyConfig{message}
}

type errEmptyConfig struct {
	message string
}

func (e *errEmptyConfig) Error() string {
	return e.message
}

// validateAuthInfo looks for conflicts and errors in the auth info.
func validateAuthInfo(authInfo AuthInfo) []error {
	validationErrors := make([]error, 0)

	usingAuthPath := false

	methods := make([]string, 0, 3)
	if len(authInfo.Token) != 0 {
		methods = append(methods, "token")
	}

	if len(authInfo.Username) != 0 || len(authInfo.Password) != 0 {
		methods = append(methods, "basicAuth")
	}

	if len(authInfo.SecretID) != 0 || len(authInfo.SecretKey) != 0 {
		methods = append(methods, "secretAuth")
	}

	// authPath also provides information for the client to identify the server,
	// so allow multiple auth methods in that case
	if (len(methods) > 1) && (!usingAuthPath) {
		validationErrors = append(validationErrors,
			fmt.Errorf("more than one authentication method found; found %v, only one is allowed", methods))
	}

	if len(authInfo.ClientCertificate) == 0 || len(authInfo.ClientCertificateData) == 0 {
		return validationErrors
	}

	// Make sure cert data and file aren't both specified
	if len(authInfo.ClientCertificate) != 0 && len(authInfo.ClientCertificateData) != 0 {
		validationErrors = append(validationErrors,
			fmt.Errorf("client-cert-data and client-cert are both specified. client-cert-data will override"))
	}
	// Make sure key data and file aren't both specified
	if len(authInfo.ClientKey) != 0 && len(authInfo.ClientKeyData) != 0 {
		validationErrors = append(validationErrors,
			fmt.Errorf("client-key-data and client-key are both specified; client-key-data will override"))
	}
	// Make sure a key is specified
	if len(authInfo.ClientKey) == 0 && len(authInfo.ClientKeyData) == 0 {
		validationErrors = append(validationErrors,
			fmt.Errorf("client-key-data or client-key must be specified to use the clientCert authentication method"))
	}

	if len(authInfo.ClientCertificate) != 0 {
		clientCertFile, err := os.Open(authInfo.ClientCertificate)
		if err != nil {
			validationErrors = append(validationErrors,
				fmt.Errorf("unable to read client-cert %v due to %v", authInfo.ClientCertificate, err))
		} else {
			defer clientCertFile.Close()
		}
	}

	if len(authInfo.ClientKey) != 0 {
		clientKeyFile, err := os.Open(authInfo.ClientKey)
		if err != nil {
			validationErrors = append(validationErrors,
				fmt.Errorf("unable to read client-key %v due to %v", authInfo.ClientKey, err))
		} else {
			defer clientKeyFile.Close()
		}
	}

	return validationErrors
}

// validateServerInfo looks for conflicts and errors in the server info.
func validateServerInfo(serverInfo Server) []error {
	validationErrors := make([]error, 0)

	emptyServer := &Server{}
	if reflect.DeepEqual(*emptyServer, serverInfo) {
		return []error{ErrEmptyServer}
	}

	/*
		if len(serverInfo.Address) == 0 {
			validationErrors = append(validationErrors, fmt.Errorf("no server found"))
		}
	*/
	// Make sure CA data and CA file aren't both specified
	if len(serverInfo.CertificateAuthority) != 0 && len(serverInfo.CertificateAuthorityData) != 0 {
		validationErrors = append(
			validationErrors,
			fmt.Errorf(
				"certificate-authority-data and certificate-authority are both specified. certificate-authority-data will override",
			),
		)
	}

	if len(serverInfo.CertificateAuthority) != 0 {
		clientCertCA, err := os.Open(serverInfo.CertificateAuthority)
		if err != nil {
			validationErrors = append(validationErrors,
				fmt.Errorf("unable to read certificate-authority %v due to %v", serverInfo.CertificateAuthority, err))
		} else {
			defer clientCertCA.Close()
		}
	}

	return validationErrors
}

func newErrConfigurationInvalid(errs []error) error {
	switch len(errs) {
	case 0:
		return nil
	default:
		return errConfigurationInvalid(errs)
	}
}
