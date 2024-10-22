// Package cgroup provides an simple abstraction over the cgroup v1 interface.
// The rationale to implement was driven by the fact that on my system, both
// the v1 and v2 interfaces are mounted and the v1 interface is in use.
//
// I've versioned the package structure so that we could easily build a v2
// implementation.  We could also add some code to the parent package to
// enable us to examine the system and automatically pick the most suitable
// implementation.
package cgroup
