package main

import "os"

func Run(code string) result {
	dir, err := os.MkdirTemp("", "jail-*")
	if err != nil {
		return result{Stderr: err.Error(), ExitCode: -1}
	}
	_ = dir

	return result{}
}
