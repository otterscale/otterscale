package core

// Version is the build-time binary version (e.g. "v1.2.3").
// It is a distinct type so that Wire can distinguish it from plain
// strings when injecting dependencies.
type Version string
