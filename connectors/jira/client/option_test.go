package client

import (
	"testing"
	"time"
)

func TestWithName(t *testing.T) {
	opt := WithName("test-name")
	options := &options{}
	opt.apply(options)

	if options.name != "test-name" {
		t.Errorf("expected name to be 'test-name', got '%s'", options.name)
	}
}

func TestWithServer(t *testing.T) {
	opt := WithServer("test-server")
	options := &options{}
	opt.apply(options)

	if options.server != "test-server" {
		t.Errorf("expected server to be 'test-server', got '%s'", options.server)
	}
}

func TestWithUsername(t *testing.T) {
	opt := WithUsername("test-username")
	options := &options{}
	opt.apply(options)

	if options.username != "test-username" {
		t.Errorf("expected username to be 'test-username', got '%s'", options.username)
	}
}

func TestWithPassword(t *testing.T) {
	opt := WithPassword("test-password")
	options := &options{}
	opt.apply(options)

	if options.password != "test-password" {
		t.Errorf("expected password to be 'test-password', got '%s'", options.password)
	}
}

func TestWithToken(t *testing.T) {
	opt := WithToken("test-token")
	options := &options{}
	opt.apply(options)

	if options.token != "test-token" {
		t.Errorf("expected token to be 'test-token', got '%s'", options.token)
	}
}

func TestWithProjects(t *testing.T) {
	opt := WithProjects([]string{"test-projects"})
	options := &options{}
	opt.apply(options)

	if len(options.projects) != 1 || options.projects[0] != "test-projects" {
		t.Errorf("expected projects to contain ['test-projects'], got %v", options.projects)
	}
}

func TestWithStartDate(t *testing.T) {
	opt := WithStartDate("test-startDate")
	options := &options{}
	opt.apply(options)

	if options.startDate != "test-startDate" {
		t.Errorf("expected startDate to be 'test-startDate', got '%s'", options.startDate)
	}
}

func TestWithBatchSize(t *testing.T) {
	opt := WithBatchSize(100)
	options := &options{}
	opt.apply(options)

	if options.batchSize != 100 {
		t.Errorf("expected batchSize to be 100, got %d", options.batchSize)
	}
}

func TestWithBatchSizeBytes(t *testing.T) {
	opt := WithBatchSizeBytes(1024)
	options := &options{}
	opt.apply(options)

	if options.batchSizeBytes != 1024 {
		t.Errorf("expected batchSizeBytes to be 1024, got %d", options.batchSizeBytes)
	}
}

func TestWithBatchTimeout(t *testing.T) {
	timeout := 5 * time.Second
	opt := WithBatchTimeout(timeout)
	options := &options{}
	opt.apply(options)

	if options.batchTimeout != timeout {
		t.Errorf("expected batchTimeout to be %v, got %v", timeout, options.batchTimeout)
	}
}

func TestWithCreateIndex(t *testing.T) {
	opt := WithCreateIndex(true)
	options := &options{}
	opt.apply(options)

	if options.createIndex != true {
		t.Errorf("expected createIndex to be true, got %v", options.createIndex)
	}
}
