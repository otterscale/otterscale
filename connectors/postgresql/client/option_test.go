package client

import (
	"testing"
	"time"
)

func TestWithName(t *testing.T) {
	opt := options{}
	name := "test-name"
	WithName(name).apply(&opt)
	if opt.name != name {
		t.Errorf("WithName() = %v, want %v", opt.name, name)
	}
}

func TestWithConnString(t *testing.T) {
	opt := options{}
	connString := "postgres://user:password@localhost:5432/dbname"
	WithConnString(connString).apply(&opt)
	if opt.connString != connString {
		t.Errorf("WithConnString() = %v, want %v", opt.connString, connString)
	}
}

func TestWithNamespace(t *testing.T) {
	opt := options{}
	namespace := "test-namespace"
	WithNamespace(namespace).apply(&opt)
	if opt.namespace != namespace {
		t.Errorf("WithNamespace() = %v, want %v", opt.namespace, namespace)
	}
}

func TestWithBatchSize(t *testing.T) {
	opt := options{}
	batchSize := int64(100)
	WithBatchSize(batchSize).apply(&opt)
	if opt.batchSize != batchSize {
		t.Errorf("WithBatchSize() = %v, want %v", opt.batchSize, batchSize)
	}
}

func TestWithBatchSizeBytes(t *testing.T) {
	opt := options{}
	batchSizeBytes := int64(1024)
	WithBatchSizeBytes(batchSizeBytes).apply(&opt)
	if opt.batchSizeBytes != batchSizeBytes {
		t.Errorf("WithBatchSizeBytes() = %v, want %v", opt.batchSizeBytes, batchSizeBytes)
	}
}

func TestWithBatchTimeout(t *testing.T) {
	opt := options{}
	timeout := 5 * time.Second
	WithBatchTimeout(timeout).apply(&opt)
	if opt.batchTimeout != timeout {
		t.Errorf("WithBatchTimeout() = %v, want %v", opt.batchTimeout, timeout)
	}
}

func TestWithCreateIndex(t *testing.T) {
	opt := options{}
	createIndex := true
	WithCreateIndex(createIndex).apply(&opt)
	if opt.createIndex != createIndex {
		t.Errorf("WithCreateIndex() = %v, want %v", opt.createIndex, createIndex)
	}
}

func TestMultipleOptions(t *testing.T) {
	opt := options{}
	name := "test-name"
	connString := "postgres://localhost:5432"
	batchSize := int64(200)
	
	// Apply multiple options
	options := []Option{
		WithName(name),
		WithConnString(connString),
		WithBatchSize(batchSize),
	}
	
	for _, option := range options {
		option.apply(&opt)
	}
	
	// Verify all options were applied correctly
	if opt.name != name {
		t.Errorf("name = %v, want %v", opt.name, name)
	}
	if opt.connString != connString {
		t.Errorf("connString = %v, want %v", opt.connString, connString)
	}
	if opt.batchSize != batchSize {
		t.Errorf("batchSize = %v, want %v", opt.batchSize, batchSize)
	}
}