package errors

import "fmt"

// InvalidPath returns a standardized invalid path error.
func InvalidPath(pathType string, err error) error {
	if err != nil {
		return fmt.Errorf("invalid %s path: %w", pathType, err)
	}
	return fmt.Errorf("invalid %s path", pathType)
}

// FileNotFound returns a standardized file not found error.
func FileNotFound(path string) error {
	return fmt.Errorf("file not found: %s", path)
}

// DirNotFound returns a standardized directory not found error.
func DirNotFound(path string) error {
	return fmt.Errorf("directory does not exist: %s", path)
}

// ValidationError returns a standardized validation error.
func ValidationError(message string) error {
	return fmt.Errorf("validation error: %s", message)
}

// ValidationErrorf returns a formatted validation error.
func ValidationErrorf(format string, args ...interface{}) error {
	return fmt.Errorf("validation error: %s", fmt.Sprintf(format, args...))
}

// RequiredFlag returns a standardized required flag error with optional examples.
func RequiredFlag(flagName string, examples ...string) error {
	msg := fmt.Sprintf("--%s flag is required", flagName)
	if len(examples) > 0 {
		msg += "\n\nExamples:"
		for _, ex := range examples {
			msg += "\n  " + ex
		}
	}
	return fmt.Errorf("%s", msg)
}

// MutuallyExclusive returns an error for mutually exclusive flags.
func MutuallyExclusive(flag1, flag2 string) error {
	return fmt.Errorf("--%s and --%s cannot be used together", flag1, flag2)
}

// MinValue returns an error for values below minimum.
func MinValue(name string, minValue int) error {
	return fmt.Errorf("%s must be at least %d", name, minValue)
}

// MaxValue returns an error for values above maximum.
func MaxValue(name string, maxValue int) error {
	return fmt.Errorf("%s must be at most %d", name, maxValue)
}

// Range returns an error for values outside a range.
func Range(name string, min, max int) error {
	return fmt.Errorf("%s must be between %d and %d", name, min, max)
}

// EmptyValue returns an error for empty values.
func EmptyValue(name string) error {
	return fmt.Errorf("%s cannot be empty", name)
}

// InvalidValue returns an error for invalid values.
func InvalidValue(name string, value interface{}, reason string) error {
	if reason != "" {
		return fmt.Errorf("invalid %s %q: %s", name, value, reason)
	}
	return fmt.Errorf("invalid %s: %v", name, value)
}
