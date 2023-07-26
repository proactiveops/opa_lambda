// policyloader/error.go
package policyloader

import "fmt"

// FileNotFoundError is returned when a policy file cannot be found.
type FileNotFoundError struct {
	Key string
}

// Error returns the error message.
func (e *FileNotFoundError) Error() string {
	return fmt.Sprintf("unable to locate policy file: %s", e.Key)
}

// InvalidKeyNameError is returned when a policy key name contains a slash.
type InvalidKeyNameError struct {
	Key string
}

// Error returns the error message.
func (e *InvalidKeyNameError) Error() string {
	return fmt.Sprintf("policy key name contains slash: %s", e.Key)
}
