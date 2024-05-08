package debug

import "flag"

type RunMode struct {
	devMode bool
}

func NewRunMode() *RunMode {
	devMode := flag.Bool("dev", false, "Run in dev mode")
	flag.Parse()
	return &RunMode{
		devMode: *devMode,
	}
}

func (r *RunMode) IsDevMode() bool {
	return r.devMode
}
