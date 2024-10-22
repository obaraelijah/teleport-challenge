package syscall

import gosyscall "syscall"

// Adapter serves as a shim between between callers of standard syscall.* APIs
// and the functions themselves.  The default behavior is simply to dispatch
// the call to the corresponding syscall.* function.
//
// If the field associated with that function is non-nil, then the adapter will
// dispatch to that function instead.
type Adapter struct {
	ExecFn func(argv0 string, argv []string, envv []string) (err error)
}

func (a *Adapter) Exec(argv0 string, argv []string, envv []string) (err error) {
	fn := gosyscall.Exec
	if a != nil && a.ExecFn != nil {
		fn = a.ExecFn
	}

	return fn(argv0, argv, envv)
}
