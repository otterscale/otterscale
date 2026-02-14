package pki

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
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

	// Each call produces a fresh key, so the certs must differ.
	if bytes.Equal(ca1.CertPEM(), ca2.CertPEM()) {
		t.Error("expected different CA certs from two NewCA calls")
	}
}

func TestLoadCA_Roundtrip(t *testing.T) {
	// Generate a CA, export its PEM material, and reload it.
	original, err := NewCA()
	if err != nil {
		t.Fatalf("NewCA: %v", err)
	}

	keyPEM, err := original.KeyPEM()
	if err != nil {
		t.Fatalf("KeyPEM: %v", err)
	}

	loaded, err := LoadCA(original.CertPEM(), keyPEM)
	if err != nil {
		t.Fatalf("LoadCA: %v", err)
	}

	// The loaded CA must produce the same cert PEM.
	if !bytes.Equal(original.CertPEM(), loaded.CertPEM()) {
		t.Error("loaded CA cert PEM differs from original")
	}

	// The loaded CA must be able to sign a CSR.
	key, _, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	csrPEM, err := GenerateCSR(key, "roundtrip-agent")
	if err != nil {
		t.Fatalf("GenerateCSR: %v", err)
	}
	certPEM, err := loaded.SignCSR(csrPEM)
	if err != nil {
		t.Fatalf("loaded CA SignCSR: %v", err)
	}

	// Verify the certificate was signed by the original CA.
	block, _ := pem.Decode(certPEM)
	if block == nil {
		t.Fatal("failed to decode signed cert PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("parse signed cert: %v", err)
	}
	pool := x509.NewCertPool()
	pool.AddCert(original.cert)
	if _, err := cert.Verify(x509.VerifyOptions{
		Roots:     pool,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}); err != nil {
		t.Errorf("certificate verification failed: %v", err)
	}
}

func TestLoadCA_InvalidCertPEM(t *testing.T) {
	_, err := LoadCA([]byte("not-a-cert"), []byte("not-a-key"))
	if err == nil {
		t.Error("expected error for invalid cert PEM, got nil")
	}
}

func TestLoadCA_InvalidKeyPEM(t *testing.T) {
	ca, err := NewCA()
	if err != nil {
		t.Fatalf("NewCA: %v", err)
	}
	_, err = LoadCA(ca.CertPEM(), []byte("not-a-key"))
	if err == nil {
		t.Error("expected error for invalid key PEM, got nil")
	}
}

func TestLoadCA_NotCA(t *testing.T) {
	// Generate a server cert (not a CA) and try to load it as a CA.
	ca, err := NewCA()
	if err != nil {
		t.Fatalf("NewCA: %v", err)
	}
	serverCert, serverKey, err := ca.GenerateServerCert("127.0.0.1")
	if err != nil {
		t.Fatalf("GenerateServerCert: %v", err)
	}
	_, err = LoadCA(serverCert, serverKey)
	if err == nil {
		t.Error("expected error when loading non-CA cert, got nil")
	}
}

func TestLoadCA_KeyMismatch(t *testing.T) {
	ca1, err := NewCA()
	if err != nil {
		t.Fatalf("NewCA 1: %v", err)
	}
	ca2, err := NewCA()
	if err != nil {
		t.Fatalf("NewCA 2: %v", err)
	}
	keyPEM2, err := ca2.KeyPEM()
	if err != nil {
		t.Fatalf("KeyPEM: %v", err)
	}
	// Cert from ca1, key from ca2.
	_, err = LoadCA(ca1.CertPEM(), keyPEM2)
	if err == nil {
		t.Error("expected error for key mismatch, got nil")
	}
}

func TestKeyPEM(t *testing.T) {
	ca, err := NewCA()
	if err != nil {
		t.Fatalf("NewCA: %v", err)
	}

	keyPEM, err := ca.KeyPEM()
	if err != nil {
		t.Fatalf("KeyPEM: %v", err)
	}

	if len(keyPEM) == 0 {
		t.Fatal("expected non-empty key PEM")
	}

	block, _ := pem.Decode(keyPEM)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		t.Fatal("expected EC PRIVATE KEY PEM block")
	}
}

func TestCA_DeriveHMACKey_Deterministic(t *testing.T) {
	ca, err := NewCA()
	if err != nil {
		t.Fatalf("NewCA: %v", err)
	}

	k1, err := ca.DeriveHMACKey("label-1")
	if err != nil {
		t.Fatalf("DeriveHMACKey: %v", err)
	}
	k2, err := ca.DeriveHMACKey("label-1")
	if err != nil {
		t.Fatalf("DeriveHMACKey: %v", err)
	}

	if len(k1) != 32 {
		t.Fatalf("expected 32-byte key, got %d", len(k1))
	}

	if !bytes.Equal(k1, k2) {
		t.Error("expected identical keys for same CA and label")
	}
}

func TestCA_DeriveHMACKey_DifferentLabels(t *testing.T) {
	ca, err := NewCA()
	if err != nil {
		t.Fatalf("NewCA: %v", err)
	}

	k1, err := ca.DeriveHMACKey("label-a")
	if err != nil {
		t.Fatalf("DeriveHMACKey: %v", err)
	}
	k2, err := ca.DeriveHMACKey("label-b")
	if err != nil {
		t.Fatalf("DeriveHMACKey: %v", err)
	}

	if bytes.Equal(k1, k2) {
		t.Error("expected different keys for different labels")
	}
}

func TestCA_DeriveHMACKey_DifferentCAs(t *testing.T) {
	ca1, err := NewCA()
	if err != nil {
		t.Fatalf("NewCA 1: %v", err)
	}
	ca2, err := NewCA()
	if err != nil {
		t.Fatalf("NewCA 2: %v", err)
	}

	k1, err := ca1.DeriveHMACKey("same-label")
	if err != nil {
		t.Fatalf("DeriveHMACKey 1: %v", err)
	}
	k2, err := ca2.DeriveHMACKey("same-label")
	if err != nil {
		t.Fatalf("DeriveHMACKey 2: %v", err)
	}

	if bytes.Equal(k1, k2) {
		t.Error("expected different keys for different CAs")
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

	// Verify the certificate was signed by the CA.
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

	// Verify signed by CA.
	pool := x509.NewCertPool()
	pool.AddCert(ca.cert)
	if _, err := cert.Verify(x509.VerifyOptions{
		Roots:     pool,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}); err != nil {
		t.Errorf("certificate verification failed: %v", err)
	}
}

func TestDeriveAuth(t *testing.T) {
	ca, err := NewCA()
	if err != nil {
		t.Fatalf("NewCA: %v", err)
	}

	key, _, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}

	csrPEM, err := GenerateCSR(key, "agent-1")
	if err != nil {
		t.Fatalf("GenerateCSR: %v", err)
	}

	certPEM, err := ca.SignCSR(csrPEM)
	if err != nil {
		t.Fatalf("SignCSR: %v", err)
	}

	auth1, err := DeriveAuth("agent-1", certPEM)
	if err != nil {
		t.Fatalf("DeriveAuth: %v", err)
	}

	// Same inputs should produce the same auth.
	auth2, err := DeriveAuth("agent-1", certPEM)
	if err != nil {
		t.Fatalf("DeriveAuth: %v", err)
	}

	if auth1 != auth2 {
		t.Error("expected deterministic auth string, got different results")
	}

	// Should contain the agent ID.
	if len(auth1) < len("agent-1:")+1 {
		t.Errorf("auth string too short: %s", auth1)
	}
	if auth1[:len("agent-1:")] != "agent-1:" {
		t.Errorf("expected auth to start with agent-1:, got %s", auth1)
	}
}

func TestDeriveAuth_InvalidPEM(t *testing.T) {
	_, err := DeriveAuth("agent", []byte("not-a-pem"))
	if err == nil {
		t.Error("expected error for invalid PEM, got nil")
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
