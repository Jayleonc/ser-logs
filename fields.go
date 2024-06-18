package serlogs

// Field represents a key-value pair for log content.
type Field struct {
	Key string
	Val any
}

// Error creates a Field with an error value.
func Error(err error) Field {
	return Field{"error", err}
}

// Int64 creates a Field with an int64 value.
func Int64(key string, val int64) Field {
	return Field{key, val}
}

// Uint8 creates a Field with a uint8 value.
func Uint8(key string, val uint8) Field {
	return Field{key, val}
}

// String creates a Field with a string value.
func String(key, val string) Field {
	return Field{key, val}
}

// Slice creates a Field with a slice of values.
func Slice(key string, val ...any) Field {
	return Field{key, val}
}

// Bool creates a Field with a bool value.
func Bool(key string, val bool) Field {
	return Field{key, val}
}

// Any creates a Field with any type of value.
func Any(key string, val any) Field {
	return Field{key, val}
}
