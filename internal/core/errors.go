package core

import (
	"errors"
	"fmt"
)

// ErrorCode represents a domain-level error category that abstracts
// away infrastructure-specific error codes (e.g. K8s StatusReason).
// The handler layer maps these codes to transport-level codes
// (e.g. ConnectRPC codes).
type ErrorCode int

const (
	ErrorCodeInternal           ErrorCode = iota // catch-all
	ErrorCodeInvalidArgument                     // bad input
	ErrorCodeNotFound                            // resource missing
	ErrorCodeAlreadyExists                       // duplicate
	ErrorCodeUnauthenticated                     // no/invalid creds
	ErrorCodePermissionDenied                    // forbidden
	ErrorCodeFailedPrecondition                  // conflict / precondition
	ErrorCodeDeadlineExceeded                    // timeout
	ErrorCodeResourceExhausted                   // rate-limit / quota
	ErrorCodeUnimplemented                       // method not allowed
	ErrorCodeUnavailable                         // service unavailable
	ErrorCodeCanceled                            // caller gave up
)

// DomainError is a generic domain error carrying an ErrorCode and an
// optional cause. Infrastructure adapters wrap external errors into
// DomainErrors so that the handler layer only needs to understand
// domain-level codes, not infrastructure-specific error types.
type DomainError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *DomainError) Unwrap() error { return e.Cause }

// DomainErrorCode extracts the ErrorCode from an error if it is a
// *DomainError. Returns ErrorCodeInternal and false for non-domain errors.
func DomainErrorCode(err error) (ErrorCode, bool) {
	var de *DomainError
	if errors.As(err, &de) {
		return de.Code, true
	}
	return ErrorCodeInternal, false
}

// ErrClusterNotFound indicates that the requested cluster is not
// registered with the tunnel provider.
type ErrClusterNotFound struct {
	Cluster string
}

func (e *ErrClusterNotFound) Error() string {
	return fmt.Sprintf("cluster %s not registered", e.Cluster)
}

// ErrNotReady indicates that a required subsystem (e.g. the tunnel
// server) has not been initialized yet.
type ErrNotReady struct {
	Subsystem string
}

func (e *ErrNotReady) Error() string {
	return fmt.Sprintf("%s not initialized", e.Subsystem)
}

// ErrInvalidInput indicates a domain-level input validation failure.
// It replaces the use of k8s apierrors.NewBadRequest in the domain
// layer, keeping the core package free of infrastructure error types.
type ErrInvalidInput struct {
	Field   string
	Message string
}

func (e *ErrInvalidInput) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("invalid %s: %s", e.Field, e.Message)
	}
	return e.Message
}

// ErrSessionNotFound indicates that a requested session (exec or
// port-forward) does not exist in the session store.
type ErrSessionNotFound struct {
	Resource string
	ID       string
}

func (e *ErrSessionNotFound) Error() string {
	return fmt.Sprintf("%s %q not found", e.Resource, e.ID)
}

// ErrCommandExited reports that an exec'd command ran to completion and
// exited with a non-zero status.
//
// This is deliberately distinct from every other error an exec session
// can end with. Failing to *start* a session — no such pod, no such
// container, RBAC denial, a failed protocol upgrade — is a failure of
// the RPC. A command that started, produced its output, and exited 1 is
// not: the caller got exactly what it asked for, and reporting it as a
// transport failure would turn every unsuccessful shell command into a
// server error.
type ErrCommandExited struct {
	Code   int
	Reason string
}

func (e *ErrCommandExited) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("command exited with status %d: %s", e.Code, e.Reason)
	}
	return fmt.Sprintf("command exited with status %d", e.Code)
}
