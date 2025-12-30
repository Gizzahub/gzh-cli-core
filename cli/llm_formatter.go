package cli

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"
)

const (
	maxLLMDepth = 5
	indentSize  = 2
)

// llmFormatter formats data in a token-efficient format for LLM consumption.
type llmFormatter struct{}

// format converts any data to LLM-friendly string format.
func (l *llmFormatter) format(data interface{}, depth int) string {
	if data == nil {
		return ""
	}
	if depth > maxLLMDepth {
		return "..."
	}

	v := reflect.ValueOf(data)
	return l.formatValue(v, depth)
}

// formatValue handles reflect.Value formatting.
func (l *llmFormatter) formatValue(v reflect.Value, depth int) string {
	// Dereference pointers and interfaces
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	// Check for special types first
	if formatted, ok := l.formatSpecialType(v); ok {
		return formatted
	}

	switch v.Kind() {
	case reflect.Struct:
		return l.formatStruct(v, depth)
	case reflect.Slice, reflect.Array:
		return l.formatSlice(v, depth)
	case reflect.Map:
		return l.formatMap(v, depth)
	case reflect.String:
		return v.String()
	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() == 0 {
			return ""
		}
		return fmt.Sprintf("%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Uint() == 0 {
			return ""
		}
		return fmt.Sprintf("%d", v.Uint())
	case reflect.Float32, reflect.Float64:
		if v.Float() == 0 {
			return ""
		}
		return fmt.Sprintf("%g", v.Float())
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}

// formatStruct formats a struct with UPPER_CASE field labels.
func (l *llmFormatter) formatStruct(v reflect.Value, depth int) string {
	if depth > maxLLMDepth {
		return "..."
	}

	var sb strings.Builder
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		fieldValue := v.Field(i)

		// Skip empty/zero fields
		if l.shouldSkipField(fieldValue) {
			continue
		}

		fieldLabel := l.fieldNameToLabel(field.Name)
		formatted := l.formatValue(fieldValue, depth+1)

		if formatted == "" {
			continue
		}

		indent := l.indent(depth)

		// Multi-line values (nested structs/slices) get indented on next line
		if strings.Contains(formatted, "\n") {
			sb.WriteString(fmt.Sprintf("%s%s:\n%s", indent, fieldLabel, formatted))
		} else {
			sb.WriteString(fmt.Sprintf("%s%s: %s\n", indent, fieldLabel, formatted))
		}
	}

	return sb.String()
}

// formatSlice formats slices - primitives with pipe separator, structs with indices.
func (l *llmFormatter) formatSlice(v reflect.Value, depth int) string {
	if v.Len() == 0 {
		return ""
	}

	if depth > maxLLMDepth {
		return "..."
	}

	var sb strings.Builder

	// Check if it's a slice of primitives
	if v.Len() > 0 && !l.isComplexType(v.Index(0)) {
		var items []string
		for i := 0; i < v.Len(); i++ {
			formatted := l.formatValue(v.Index(i), depth)
			if formatted != "" {
				items = append(items, formatted)
			}
		}
		if len(items) == 0 {
			return ""
		}
		return strings.Join(items, " | ")
	}

	// For struct slices, use numbered items
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		formatted := l.formatValue(elem, depth+1)

		if formatted != "" {
			sb.WriteString(fmt.Sprintf("%s[%d]\n%s", l.indent(depth), i, formatted))
		}
	}

	return sb.String()
}

// formatMap formats maps with key-value pairs.
func (l *llmFormatter) formatMap(v reflect.Value, depth int) string {
	if v.Len() == 0 {
		return ""
	}

	if depth > maxLLMDepth {
		return "..."
	}

	var sb strings.Builder
	iter := v.MapRange()

	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		keyStr := l.formatValue(key, depth)
		valueStr := l.formatValue(value, depth+1)

		if valueStr == "" {
			continue
		}

		indent := l.indent(depth)

		if strings.Contains(valueStr, "\n") {
			sb.WriteString(fmt.Sprintf("%s%s:\n%s", indent, keyStr, valueStr))
		} else {
			sb.WriteString(fmt.Sprintf("%s%s: %s\n", indent, keyStr, valueStr))
		}
	}

	return sb.String()
}

// shouldSkipField returns true if the field should be omitted from output.
func (l *llmFormatter) shouldSkipField(v reflect.Value) bool {
	// Dereference pointers
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return true
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Slice, reflect.Array, reflect.Map:
		return v.Len() == 0
	case reflect.Struct:
		// Special handling for time.Time
		if t, ok := v.Interface().(time.Time); ok {
			return t.IsZero()
		}
		return false
	default:
		return false
	}
}

// fieldNameToLabel converts CamelCase to UPPER_CASE.
func (l *llmFormatter) fieldNameToLabel(name string) string {
	var result strings.Builder
	for i, r := range name {
		if i > 0 && unicode.IsUpper(r) {
			// Check if we should insert underscore
			prev := rune(name[i-1])
			if unicode.IsLower(prev) {
				result.WriteRune('_')
			} else if i+1 < len(name) {
				next := rune(name[i+1])
				if unicode.IsLower(next) {
					result.WriteRune('_')
				}
			}
		}
		result.WriteRune(unicode.ToUpper(r))
	}
	return result.String()
}

// formatSpecialType handles special types like time.Time, time.Duration.
func (l *llmFormatter) formatSpecialType(v reflect.Value) (string, bool) {
	if !v.IsValid() || !v.CanInterface() {
		return "", false
	}

	iface := v.Interface()

	// time.Time
	if t, ok := iface.(time.Time); ok {
		if t.IsZero() {
			return "", true
		}
		return t.Format(time.RFC3339), true
	}

	// time.Duration
	if d, ok := iface.(time.Duration); ok {
		if d == 0 {
			return "", true
		}
		return d.String(), true
	}

	// []byte - show hex
	if b, ok := iface.([]byte); ok {
		if len(b) == 0 {
			return "", true
		}
		if len(b) > 32 {
			return fmt.Sprintf("%x...(%d bytes)", b[:32], len(b)), true
		}
		return fmt.Sprintf("%x", b), true
	}

	// error interface
	if err, ok := iface.(error); ok {
		if err == nil {
			return "", true
		}
		return err.Error(), true
	}

	return "", false
}

// indent returns indentation string for given depth.
func (l *llmFormatter) indent(depth int) string {
	return strings.Repeat(" ", depth*indentSize)
}

// isComplexType checks if the value is a struct, slice, or map.
func (l *llmFormatter) isComplexType(v reflect.Value) bool {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return false
		}
		v = v.Elem()
	}
	k := v.Kind()
	return k == reflect.Struct || k == reflect.Slice || k == reflect.Array || k == reflect.Map
}
