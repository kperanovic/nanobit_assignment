package common

import (
	"os"
	"os/signal"
)

// GetEnv simply returns an evironment variable value.
// If it is not set or empty, it returns the default value
func GetEnv(name string, def string) string {
	env := os.Getenv(name)
	if env == "" {
		env = def
	}
	return env
}

// TrapSignals hangs unti SIGINT is sent to the process
func TrapSignals() chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	return sig
}
