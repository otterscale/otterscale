// Package client provides tests for various option configurations used in the CSV connector.
//
// The tests included in this file are:
// - TestWithName: Verifies that the WithName option correctly sets the name field.
// - TestWithFilePath: Verifies that the WithFilePath option correctly sets the filePath field.
// - TestWithTableName: Verifies that the WithTableName option correctly sets the tableName field.
// - TestWithInferring: Verifies that the WithInferring option correctly sets the inferring field.
// - TestWithBatchSize: Verifies that the WithBatchSize option correctly sets the batchSize field.
// - TestWithBatchSizeBytes: Verifies that the WithBatchSizeBytes option correctly sets the batchSizeBytes field.
// - TestWithBatchTimeout: Verifies that the WithBatchTimeout option correctly sets the batchTimeout field.
// - TestWithCreateIndex: Verifies that the WithCreateIndex option correctly sets the createIndex field.
//
// Each test function initializes an options struct, applies the respective option function, and checks if the field is set correctly.
package client

import (
	"testing"
	"time"
)

func TestWithName(t *testing.T) {
	opt := options{}
	WithName("testName").apply(&opt)
	if opt.name != "testName" {
		t.Errorf("expected name to be 'testName', got %s", opt.name)
	}
}

func TestWithFilePath(t *testing.T) {
	opt := options{}
	WithFilePath("/path/to/file").apply(&opt)
	if opt.filePath != "/path/to/file" {
		t.Errorf("expected filePath to be '/path/to/file', got %s", opt.filePath)
	}
}

func TestWithTableName(t *testing.T) {
	opt := options{}
	WithTableName("testTable").apply(&opt)
	if opt.tableName != "testTable" {
		t.Errorf("expected tableName to be 'testTable', got %s", opt.tableName)
	}
}

func TestWithInferring(t *testing.T) {
	opt := options{}
	WithInferring(true).apply(&opt)
	if !opt.inferring {
		t.Errorf("expected inferring to be true, got %v", opt.inferring)
	}
}

func TestWithBatchSize(t *testing.T) {
	opt := options{}
	WithBatchSize(100).apply(&opt)
	if opt.batchSize != 100 {
		t.Errorf("expected batchSize to be 100, got %d", opt.batchSize)
	}
}

func TestWithBatchSizeBytes(t *testing.T) {
	opt := options{}
	WithBatchSizeBytes(1024).apply(&opt)
	if opt.batchSizeBytes != 1024 {
		t.Errorf("expected batchSizeBytes to be 1024, got %d", opt.batchSizeBytes)
	}
}

func TestWithBatchTimeout(t *testing.T) {
	opt := options{}
	timeout := 5 * time.Second
	WithBatchTimeout(timeout).apply(&opt)
	if opt.batchTimeout != timeout {
		t.Errorf("expected batchTimeout to be %v, got %v", timeout, opt.batchTimeout)
	}
}

func TestWithCreateIndex(t *testing.T) {
	opt := options{}
	WithCreateIndex(true).apply(&opt)
	if !opt.createIndex {
		t.Errorf("expected createIndex to be true, got %v", opt.createIndex)
	}
}


