package client

import (
	"testing"
	"time"
)

func TestWithName(t *testing.T) {
	opt := WithName("test-name")
	opts := &options{}
	opt.apply(opts)

	if opts.name != "test-name" {
		t.Errorf("expected name to be 'test-name', got '%s'", opts.name)
	}
}

func TestWithOwner(t *testing.T) {
	opt := WithOwner("test-owner")
	opts := &options{}
	opt.apply(opts)

	if opts.owner != "test-owner" {
		t.Errorf("expected owner to be 'test-owner', got '%s'", opts.owner)
	}
}

func TestWithRepo(t *testing.T) {
	opt := WithRepo("test-repo")
	opts := &options{}
	opt.apply(opts)

	if opts.repo != "test-repo" {
		t.Errorf("expected repo to be 'test-repo', got '%s'", opts.repo)
	}
}

func TestWithOpts(t *testing.T) {
	opt := WithOpts("test-opts")
	opts := &options{}
	opt.apply(opts)

	if opts.opts != "test-opts" {
		t.Errorf("expected opts to be 'test-opts', got '%s'", opts.opts)
	}
}

func TestWithBatchSize(t *testing.T) {
	opt := WithBatchSize(100)
	opts := &options{}
	opt.apply(opts)

	if opts.batchSize != 100 {
		t.Errorf("expected batchSize to be 100, got %d", opts.batchSize)
	}
}

func TestWithBatchSizeBytes(t *testing.T) {
	opt := WithBatchSizeBytes(1024)
	opts := &options{}
	opt.apply(opts)

	if opts.batchSizeBytes != 1024 {
		t.Errorf("expected batchSizeBytes to be 1024, got %d", opts.batchSizeBytes)
	}
}

func TestWithBatchTimeout(t *testing.T) {
	timeout := 5 * time.Second
	opt := WithBatchTimeout(timeout)
	opts := &options{}
	opt.apply(opts)

	if opts.batchTimeout != timeout {
		t.Errorf("expected batchTimeout to be %v, got %v", timeout, opts.batchTimeout)
	}
}

func TestWithCreateIndex(t *testing.T) {
	opt := WithCreateIndex(true)
	opts := &options{}
	opt.apply(opts)

	if opts.createIndex != true {
		t.Errorf("expected createIndex to be true, got %v", opts.createIndex)
	}
}