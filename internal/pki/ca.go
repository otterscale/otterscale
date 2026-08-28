// Package pki provides a minimal Certificate Authority for issuing
// short-lived TLS certificates used in agent-to-server mTLS tunnels.
//
// The CA key pair is generated using crypto/rand (NIST SP 800-90A DRBG
// in FIPS 140-3 mode). It is ephemeral: a fresh CA is created on every
// server start and is never persisted, so previously issued agent
// certificates stop being trusted after a restart. Agents recover by
// re-registering through the link API, which signs a new CSR.
package pki

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"
)

// certValidity is the default validity period for agent certificates
// signed by the CA. Short-lived certificates limit the blast radius
// of a compromised key and avoid the need for explicit revocation.
const certValidity = 24 * time.Hour

// CA holds a self-signed certificate authority key pair and provides
// methods for signing CSRs and generating server certificates.
type CA struct {
	cert    *x509.Certificate
	key     *ecdsa.PrivateKey
	certPEM []byte
}

// NewCA generates a new ECDSA P-256 CA key pair and self-signed
// certificate using crypto/rand.Reader. In FIPS 140-3 mode the
// reader is backed by a NIST SP 800-90A DRBG.
func NewCA() (*CA, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("pki: generate CA key: %w", err)
	}

	serial, err := randomSerial()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			Organization: []string{"otterscale"},
			CommonName:   "otterscale-ca",
		},
		NotBefore:             now.Add(-5 * time.Minute),
		NotAfter:              now.Add(10 * 365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            0,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		return nil, fmt.Errorf("pki: create CA cert: %w", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, fmt.Errorf("pki: parse CA cert: %w", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	return &CA{cert: cert, key: key, certPEM: certPEM}, nil
}

// CertPEM returns the PEM-encoded CA certificate. Agents use this to
// verify the tunnel server's identity and to be verified themselves.
func (ca *CA) CertPEM() []byte {
	return ca.certPEM
}

// SignCSR validates a PEM-encoded PKCS#10 certificate signing request
// and returns a PEM-encoded X.509 certificate signed by the CA. The
// certificate is valid for the default certValidity period.
func (ca *CA) SignCSR(csrPEM []byte) ([]byte, error) {
	block, _ := pem.Decode(csrPEM)
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		return nil, fmt.Errorf("pki: invalid CSR PEM")
	}

	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("pki: parse CSR: %w", err)
	}

	if err := csr.CheckSignature(); err != nil {
		return nil, fmt.Errorf("pki: CSR signature invalid: %w", err)
	}

	serial, err := randomSerial()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject:      csr.Subject,
		NotBefore:    now.Add(-5 * time.Minute),
		NotAfter:     now.Add(certValidity),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, tmpl, ca.cert, csr.PublicKey, ca.key)
	if err != nil {
		return nil, fmt.Errorf("pki: sign certificate: %w", err)
	}

	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER}), nil
}

// GenerateServerCert creates a TLS server certificate signed by the
// CA. The hosts parameter accepts IP addresses and DNS names that are
// added as Subject Alternative Names.
func (ca *CA) GenerateServerCert(hosts ...string) (certPEM, keyPEM []byte, err error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("pki: generate server key: %w", err)
	}

	serial, err := randomSerial()
	if err != nil {
		return nil, nil, err
	}

	now := time.Now()
	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			Organization: []string{"otterscale"},
			CommonName:   "otterscale-tunnel",
		},
		NotBefore:   now.Add(-5 * time.Minute),
		NotAfter:    now.Add(365 * 24 * time.Hour), // 1 year; regenerated on every server start
		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			tmpl.IPAddresses = append(tmpl.IPAddresses, ip)
		} else {
			tmpl.DNSNames = append(tmpl.DNSNames, h)
		}
	}

	certDER, err := x509.CreateCertificate(rand.Reader, tmpl, ca.cert, &key.PublicKey, ca.key)
	if err != nil {
		return nil, nil, fmt.Errorf("pki: create server cert: %w", err)
	}

	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, nil, fmt.Errorf("pki: marshal server key: %w", err)
	}

	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	return certPEM, keyPEM, nil
}

// GenerateKey creates a new ECDSA P-256 private key suitable for use
// in a CSR. It returns the key and its PEM encoding.
func GenerateKey() (*ecdsa.PrivateKey, []byte, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("pki: generate key: %w", err)
	}

	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, nil, fmt.Errorf("pki: marshal key: %w", err)
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	return key, keyPEM, nil
}

// GenerateCSR creates a PEM-encoded PKCS#10 certificate signing
// request with the given common name.
func GenerateCSR(key *ecdsa.PrivateKey, cn string) ([]byte, error) {
	tmpl := &x509.CertificateRequest{
		Subject: pkix.Name{
			Organization: []string{"otterscale"},
			CommonName:   cn,
		},
	}

	csrDER, err := x509.CreateCertificateRequest(rand.Reader, tmpl, key)
	if err != nil {
		return nil, fmt.Errorf("pki: create CSR: %w", err)
	}

	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrDER}), nil
}

// tunnelPasswordBytes is the entropy of a generated tunnel password.
const tunnelPasswordBytes = 24

// NewTunnelPassword returns a fresh random password for a tunnel user.
//
// It replaces an earlier scheme that derived the password from the
// signed certificate so that both sides could compute it without
// exchanging it. That coupled the two sides to an identical derivation
// and, more importantly, made the secret a function of a value the
// agent already holds and presents on the wire. The server now issues
// the password in the registration response instead, so it can be
// random and need not be reproducible.
func NewTunnelPassword() (string, error) {
	buf := make([]byte, tunnelPasswordBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("pki: generate tunnel password: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// randomSerial generates a cryptographically random serial number.
func randomSerial() (*big.Int, error) {
	const serialBits = 128 // 128-bit random serial number
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), serialBits))
	if err != nil {
		return nil, fmt.Errorf("pki: generate serial: %w", err)
	}
	return serial, nil
}
