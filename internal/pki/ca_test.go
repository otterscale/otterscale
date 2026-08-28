package pki

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
	"testing"
)

func TestNewCA(t *testing.T) {
	ca, err := NewCA()
	if err != nil {
		t.Fatalf("NewCA: %v", err)
	}

	if len(ca.CertPEM()) == 0 {
		t.Error("expected non-empty cert PEM")
	}

	block, _ := pem.Decode(ca.CertPEM())
	if block == nil {
		t.Fatal("failed to decode CA cert PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("parse cert: %v", err)
	}

	if !cert.IsCA {
		t.Error("expected IsCA to be true")
	}
	if cert.Subject.CommonName != "otterscale-ca" {
		t.Errorf("expected CN=otterscale-ca, got %s", cert.Subject.CommonName)
	}
	// MaxPathLen should be 0 or -1 (Go represents "0 but set" as
	// MaxPathLen=0 + MaxPathLenZero=true; when parsed it may show
	// as -1 in some Go versions).
	if cert.MaxPathLen > 0 {
		t.Errorf("expected MaxPathLen<=0, got %d", cert.MaxPathLen)
	}
}

func TestNewCA_UniquePerCall(t *testing.T) {
	ca1, err := NewCA()
	if err != nil {
		t.Fatalf("NewCA: %v", err)
	}
	ca2, err := NewCA()
	if err != nil {
		t.Fatalf("NewCA: %v", err)
	}

	// Fresh key per call, so the certs must differ.
	if bytes.Equal(ca1.CertPEM(), ca2.CertPEM()) {
		t.Error("expected different CA certs from two NewCA calls")
	}
}

func TestSignCSR(t *testing.T) {
	ca, err := NewCA()
	if err != nil {
		t.Fatalf("NewCA: %v", err)
	}

	key, _, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}

	csrPEM, err := GenerateCSR(key, "test-agent")
	if err != nil {
		t.Fatalf("GenerateCSR: %v", err)
	}

	certPEM, err := ca.SignCSR(csrPEM)
	if err != nil {
		t.Fatalf("SignCSR: %v", err)
	}

	block, _ := pem.Decode(certPEM)
	if block == nil {
		t.Fatal("failed to decode signed cert PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("parse cert: %v", err)
	}

	if cert.Subject.CommonName != "test-agent" {
		t.Errorf("expected CN=test-agent, got %s", cert.Subject.CommonName)
	}

	pool := x509.NewCertPool()
	pool.AddCert(ca.cert)
	if _, err := cert.Verify(x509.VerifyOptions{
		Roots:     pool,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}); err != nil {
		t.Errorf("certificate verification failed: %v", err)
	}
}

func TestSignCSR_InvalidPEM(t *testing.T) {
	ca, err := NewCA()
	if err != nil {
		t.Fatalf("NewCA: %v", err)
	}

	_, err = ca.SignCSR([]byte("not-a-pem"))
	if err == nil {
		t.Error("expected error for invalid PEM, got nil")
	}
}

func TestGenerateServerCert(t *testing.T) {
	ca, err := NewCA()
	if err != nil {
		t.Fatalf("NewCA: %v", err)
	}

	certPEM, keyPEM, err := ca.GenerateServerCert("127.0.0.1", "example.com")
	if err != nil {
		t.Fatalf("GenerateServerCert: %v", err)
	}

	if len(certPEM) == 0 {
		t.Error("expected non-empty cert PEM")
	}
	if len(keyPEM) == 0 {
		t.Error("expected non-empty key PEM")
	}

	block, _ := pem.Decode(certPEM)
	if block == nil {
		t.Fatal("failed to decode server cert PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("parse cert: %v", err)
	}

	if len(cert.IPAddresses) != 1 || cert.IPAddresses[0].String() != "127.0.0.1" {
		t.Errorf("expected IP SAN 127.0.0.1, got %v", cert.IPAddresses)
	}
	if len(cert.DNSNames) != 1 || cert.DNSNames[0] != "example.com" {
		t.Errorf("expected DNS SAN example.com, got %v", cert.DNSNames)
	}

	pool := x509.NewCertPool()
	pool.AddCert(ca.cert)
	if _, err := cert.Verify(x509.VerifyOptions{
		Roots:     pool,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}); err != nil {
		t.Errorf("certificate verification failed: %v", err)
	}
}

func TestNewTunnelPassword(t *testing.T) {
	pass, err := NewTunnelPassword()
	if err != nil {
		t.Fatalf("NewTunnelPassword: %v", err)
	}

	// 24 random bytes in unpadded base64url.
	if want := 32; len(pass) != want {
		t.Errorf("password length = %d, want %d (%q)", len(pass), want, pass)
	}

	// The password is the tunnel secret and is concatenated into a
	// "user:password" auth string, so a colon would split in the wrong
	// place. base64url never produces one, but the property is what
	// callers rely on.
	if strings.ContainsRune(pass, ':') {
		t.Errorf("password contains a colon: %q", pass)
	}
}

// TestNewTunnelPasswordIsUnpredictable guards the property that
// replaced the old certificate-derived scheme: the password must not be
// reproducible from anything the agent already holds.
func TestNewTunnelPasswordIsUnpredictable(t *testing.T) {
	const samples = 64

	seen := make(map[string]struct{}, samples)
	for range samples {
		pass, err := NewTunnelPassword()
		if err != nil {
			t.Fatalf("NewTunnelPassword: %v", err)
		}
		if _, dup := seen[pass]; dup {
			t.Fatalf("NewTunnelPassword repeated a value: %q", pass)
		}
		seen[pass] = struct{}{}
	}
}

func TestGenerateKey_And_CSR(t *testing.T) {
	key, keyPEM, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}

	if key == nil {
		t.Fatal("expected non-nil key")
	}
	if len(keyPEM) == 0 {
		t.Fatal("expected non-empty key PEM")
	}

	csrPEM, err := GenerateCSR(key, "test-cn")
	if err != nil {
		t.Fatalf("GenerateCSR: %v", err)
	}

	block, _ := pem.Decode(csrPEM)
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		t.Fatal("expected CERTIFICATE REQUEST PEM block")
	}

	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		t.Fatalf("parse CSR: %v", err)
	}

	if csr.Subject.CommonName != "test-cn" {
		t.Errorf("expected CN=test-cn, got %s", csr.Subject.CommonName)
	}
}

// TestSignCSRClassifiesBadRequests pins the distinction callers rely on
// to pick an error code: a CSR the caller sent badly is the caller's
// fault, while anything else means this CA could not sign.
func TestSignCSRClassifiesBadRequests(t *testing.T) {
	ca, err := NewCA()
	if err != nil {
		t.Fatalf("NewCA: %v", err)
	}

	key, _, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	validCSR, err := GenerateCSR(key, "agent-1")
	if err != nil {
		t.Fatalf("GenerateCSR: %v", err)
	}

	// A body that parses as PEM but not as a certificate request.
	notACSR := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: []byte("not a pkcs#10 body"),
	})

	tests := []struct {
		name string
		csr  []byte
	}{
		{"not PEM at all", []byte("plainly not pem")},
		{"empty", nil},
		{"wrong PEM type", pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: []byte("x")})},
		{"unparseable request", notACSR},
		{"truncated request", validCSR[:len(validCSR)/2]},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ca.SignCSR(tt.csr)
			if err == nil {
				t.Fatal("expected an error")
			}
			if !errors.Is(err, ErrInvalidCSR) {
				t.Errorf("error %v does not wrap ErrInvalidCSR, so it would be reported as a server fault", err)
			}
		})
	}

	// A well-formed request must still sign cleanly.
	if _, err := ca.SignCSR(validCSR); err != nil {
		t.Fatalf("SignCSR on a valid request: %v", err)
	}
}
