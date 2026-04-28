package runnermock

type Runner struct {
	RunCalled bool
	RunFn     func() error
}

func (r *Runner) Run(command map[string]any) error {
	r.RunCalled = true
	if r.RunFn != nil {
		return r.RunFn()
	}
	return nil
}
