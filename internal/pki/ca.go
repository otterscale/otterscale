// Package pki provides a minimal Certificate Authority for issuing
// short-lived TLS certificates used in agent-to-server mTLS tunnels.
//
// The CA key pair is generated using crypto/rand (NIST SP 800-90A DRBG
// in FIPS 140-3 mode) and must be persisted externally so that
// restarts reload the same CA, keeping previously issued agent
// certificates valid until they expire.
package pki

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hkdf"
	"crypto/rand"
	"crypto/sha256"
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
//
// The caller is responsible for persisting CertPEM() and KeyPEM()
// so that subsequent restarts can reload the same CA via LoadCA.
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

// LoadCA reconstructs a CA from PEM-encoded certificate and private
// key material. It validates that the certificate is a CA and that the
// private key matches the certificate's public key.
func LoadCA(certPEM, keyPEM []byte) (*CA, error) {
	certBlock, _ := pem.Decode(certPEM)
	if certBlock == nil {
		return nil, fmt.Errorf("pki: failed to decode CA certificate PEM")
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("pki: parse CA cert: %w", err)
	}

	if !cert.IsCA {
		return nil, fmt.Errorf("pki: certificate is not a CA")
	}

	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return nil, fmt.Errorf("pki: failed to decode CA private key PEM")
	}

	key, err := x509.ParseECPrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("pki: parse CA key: %w", err)
	}

	// Verify the private key corresponds to the certificate's public
	// key by comparing the public key parameters.
	certPub, ok := cert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("pki: CA certificate does not contain an ECDSA public key")
	}
	if !key.PublicKey.Equal(certPub) {
		return nil, fmt.Errorf("pki: CA private key does not match certificate public key")
	}

	return &CA{cert: cert, key: key, certPEM: certPEM}, nil
}

// CertPEM returns the PEM-encoded CA certificate. Agents use this to
// verify the tunnel server's identity and to be verified themselves.
func (ca *CA) CertPEM() []byte {
	return ca.certPEM
}

// KeyPEM returns the PEM-encoded CA private key for external
// persistence. The caller should store this securely (e.g. with
// 0600 permissions) so the CA can be reloaded via LoadCA.
func (ca *CA) KeyPEM() ([]byte, error) {
	keyDER, err := x509.MarshalECPrivateKey(ca.key)
	if err != nil {
		return nil, fmt.Errorf("pki: marshal CA key: %w", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}), nil
}

// DeriveHMACKey deterministically derives a 32-byte HMAC key from the
// CA's private key and a label using HKDF (RFC 5869). The derivation
// is deterministic for the same CA key and label, so the result
// survives server restarts as long as the same CA is loaded.
//
// This uses crypto/hkdf which is inside the Go Cryptographic Module's
// FIPS 140-3 boundary.
func (ca *CA) DeriveHMACKey(label string) ([]byte, error) {
	keyDER, err := x509.MarshalECPrivateKey(ca.key)
	if err != nil {
		return nil, fmt.Errorf("pki: marshal key for HKDF: %w", err)
	}
	return hkdf.Key(sha256.New, keyDER, nil, label, 32)
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

// DeriveAuth deterministically computes a chisel auth string
// ("user:password") from the agent ID and a signed certificate.
// Both the server (which signed the cert) and the agent (which
// received the cert) can independently compute this value.
func DeriveAuth(agentID string, certPEM []byte) (string, error) {
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return "", fmt.Errorf("pki: failed to decode certificate PEM")
	}
	h := sha256.Sum256(block.Bytes)
	pass := base64.RawURLEncoding.EncodeToString(h[:24])
	return agentID + ":" + pass, nil
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// randomSerial generates a cryptographically random serial number.
func randomSerial() (*big.Int, error) {
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, fmt.Errorf("pki: generate serial: %w", err)
	}
	return serial, nil
}
