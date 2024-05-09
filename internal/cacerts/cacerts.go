// Package cacerts provides an embedded CA root certificates bundle.
package cacerts

// To update the default certificates run the following command in this
// directory
//
//   curl -sfL https://curl.se/ca/cacert.pem -o ca-certificates.crt

import _ "embed"

//go:embed ca-certificates.crt
var caCertificates string

// CACertificates returns the default bundle of CA root certificates.
// The certificate bundle is under the MPL-2.0 licence from
// https://curl.se/ca/cacert.pem.
func CACertificates() string {
	return caCertificates
}
