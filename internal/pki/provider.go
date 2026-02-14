package pki

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

// ProvideCA is a Wire provider that loads the CA from the given
// directory. On first startup the directory is empty, so a new CA is
// generated (using crypto/rand backed by a FIPS-approved DRBG) and
// persisted. Subsequent restarts load the existing CA, keeping
// previously issued agent certificates valid.
func ProvideCA(dir string) (*CA, error) {
	certPath := filepath.Join(dir, "ca.pem")
	keyPath := filepath.Join(dir, "ca-key.pem")

	certPEM, errC := os.ReadFile(certPath)
	keyPEM, errK := os.ReadFile(keyPath)
	if errC == nil && errK == nil {
		slog.Info("loading existing CA", "dir", dir)
		return LoadCA(certPEM, keyPEM)
	}

	// First run: generate and persist.
	slog.Info("generating new CA", "dir", dir)
	ca, err := NewCA()
	if err != nil {
		return nil, fmt.Errorf("generate CA: %w", err)
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("create CA dir: %w", err)
	}

	keyPEM, err = ca.KeyPEM()
	if err != nil {
		return nil, fmt.Errorf("export CA key: %w", err)
	}

	// Write cert and key atomically (write to temp + rename) so
	// that a crash between the two writes does not leave a
	// half-written CA state on disk.
	if err := atomicWriteFile(certPath, ca.CertPEM(), 0600); err != nil {
		return nil, fmt.Errorf("write CA cert: %w", err)
	}
	if err := atomicWriteFile(keyPath, keyPEM, 0600); err != nil {
		return nil, fmt.Errorf("write CA key: %w", err)
	}

	return ca, nil
}

// atomicWriteFile writes data to a temporary file in the same
// directory as path, then renames it into place. This ensures that
// the target file is either fully written or not present — a crash
// mid-write cannot leave a partially written file at path.
func atomicWriteFile(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".tmp.*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmp.Name()

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmp.Chmod(perm); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("chmod temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("close temp file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}
