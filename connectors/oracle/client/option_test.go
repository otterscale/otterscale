package client

import (
	"testing"
	"time"
)

func TestWithName(t *testing.T) {
	opt := options{}
	o := WithName("testname")
	o.apply(&opt)
	if opt.name != "testname" {
		t.Errorf("expected name to be 'testname', got '%s'", opt.name)
	}
}

func TestWithConnString(t *testing.T) {
	opt := options{}
	o := WithConnString("oracle://user:pass@localhost:1521/orclpdb1")
	o.apply(&opt)
	if opt.connString != "oracle://user:pass@localhost:1521/orclpdb1" {
		t.Errorf("expected connString to be set, got '%s'", opt.connString)
	}
}

func TestWithNamespace(t *testing.T) {
	opt := options{}
	o := WithNamespace("namespace1")
	o.apply(&opt)
	if opt.namespace != "namespace1" {
		t.Errorf("expected namespace to be 'namespace1', got '%s'", opt.namespace)
	}
}

func TestWithBatchSize(t *testing.T) {
	opt := options{}
	o := WithBatchSize(100)
	o.apply(&opt)
	if opt.batchSize != 100 {
		t.Errorf("expected batchSize to be 100, got %d", opt.batchSize)
	}
}

func TestWithBatchSizeBytes(t *testing.T) {
	opt := options{}
	o := WithBatchSizeBytes(2048)
	o.apply(&opt)
	if opt.batchSizeBytes != 2048 {
		t.Errorf("expected batchSizeBytes to be 2048, got %d", opt.batchSizeBytes)
	}
}

func TestWithBatchTimeout(t *testing.T) {
	opt := options{}
	timeout := 5 * time.Second
	o := WithBatchTimeout(timeout)
	o.apply(&opt)
	if opt.batchTimeout != timeout {
		t.Errorf("expected batchTimeout to be %v, got %v", timeout, opt.batchTimeout)
	}
}

func TestWithCreateIndex(t *testing.T) {
	opt := options{}
	o := WithCreateIndex(true)
	o.apply(&opt)
	if !opt.createIndex {
		t.Errorf("expected createIndex to be true, got false")
	}
}

func TestMultipleOptions(t *testing.T) {
	opt := options{}
	opts := []Option{
		WithName("multi"),
		WithConnString("conn"),
		WithNamespace("ns"),
		WithBatchSize(10),
		WithBatchSizeBytes(1000),
		WithBatchTimeout(2 * time.Second),
		WithCreateIndex(true),
	}
	for _, o := range opts {
		o.apply(&opt)
	}
	if opt.name != "multi" ||
		opt.connString != "conn" ||
		opt.namespace != "ns" ||
		opt.batchSize != 10 ||
		opt.batchSizeBytes != 1000 ||
		opt.batchTimeout != 2*time.Second ||
		!opt.createIndex {
		t.Errorf("multiple options did not apply correctly: %+v", opt)
	}
}