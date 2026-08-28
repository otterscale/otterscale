package core

// Version is the build-time binary version (e.g. "v1.2.3"). It is a distinct
// type so Wire can tell it apart from plain strings.
type Version string
