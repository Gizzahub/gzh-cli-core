package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestLLMFormatter_BasicStruct(t *testing.T) {
	type BasicStruct struct {
		Name   string
		Count  int
		Active bool
	}

	data := BasicStruct{
		Name:   "test",
		Count:  42,
		Active: true,
	}

	var buf bytes.Buffer
	out := NewOutput().SetWriter(&buf).SetFormat("llm")
	if err := out.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "NAME: test") {
		t.Errorf("Expected NAME: test, got: %s", output)
	}
	if !strings.Contains(output, "COUNT: 42") {
		t.Errorf("Expected COUNT: 42, got: %s", output)
	}
	if !strings.Contains(output, "ACTIVE: true") {
		t.Errorf("Expected ACTIVE: true, got: %s", output)
	}
}

func TestLLMFormatter_SkipEmptyFields(t *testing.T) {
	type DataWithEmpty struct {
		Name       string
		EmptyStr   string
		ZeroInt    int
		FalseBool  bool
		EmptySlice []string
	}

	data := DataWithEmpty{
		Name:       "test",
		EmptyStr:   "",
		ZeroInt:    0,
		FalseBool:  false,
		EmptySlice: []string{},
	}

	var buf bytes.Buffer
	out := NewOutput().SetWriter(&buf).SetFormat("llm")
	if err := out.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "NAME: test") {
		t.Errorf("Expected NAME: test, got: %s", output)
	}
	if strings.Contains(output, "EMPTY_STR") {
		t.Errorf("Should not contain EMPTY_STR, got: %s", output)
	}
	if strings.Contains(output, "ZERO_INT") {
		t.Errorf("Should not contain ZERO_INT, got: %s", output)
	}
	if strings.Contains(output, "FALSE_BOOL") {
		t.Errorf("Should not contain FALSE_BOOL, got: %s", output)
	}
	if strings.Contains(output, "EMPTY_SLICE") {
		t.Errorf("Should not contain EMPTY_SLICE, got: %s", output)
	}
}

func TestLLMFormatter_NestedStruct(t *testing.T) {
	type Inner struct {
		Value string
	}
	type Outer struct {
		Name  string
		Inner Inner
	}

	data := Outer{
		Name:  "outer",
		Inner: Inner{Value: "inner-value"},
	}

	var buf bytes.Buffer
	out := NewOutput().SetWriter(&buf).SetFormat("llm")
	if err := out.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "NAME: outer") {
		t.Errorf("Expected NAME: outer, got: %s", output)
	}
	if !strings.Contains(output, "INNER:") {
		t.Errorf("Expected INNER:, got: %s", output)
	}
	if !strings.Contains(output, "VALUE: inner-value") {
		t.Errorf("Expected VALUE: inner-value, got: %s", output)
	}
}

func TestLLMFormatter_SliceOfStructs(t *testing.T) {
	type Item struct {
		Path    string
		Changes int
	}
	type Status struct {
		Files []Item
	}

	data := Status{
		Files: []Item{
			{Path: "main.go", Changes: 10},
			{Path: "README.md", Changes: 5},
		},
	}

	var buf bytes.Buffer
	out := NewOutput().SetWriter(&buf).SetFormat("llm")
	if err := out.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "FILES:") {
		t.Errorf("Expected FILES:, got: %s", output)
	}
	if !strings.Contains(output, "[0]") {
		t.Errorf("Expected [0], got: %s", output)
	}
	if !strings.Contains(output, "[1]") {
		t.Errorf("Expected [1], got: %s", output)
	}
	if !strings.Contains(output, "PATH: main.go") {
		t.Errorf("Expected PATH: main.go, got: %s", output)
	}
	if !strings.Contains(output, "PATH: README.md") {
		t.Errorf("Expected PATH: README.md, got: %s", output)
	}
}

func TestLLMFormatter_PrimitiveSlice(t *testing.T) {
	type Tags struct {
		Items []string
	}

	data := Tags{
		Items: []string{"go", "cli", "tool"},
	}

	var buf bytes.Buffer
	out := NewOutput().SetWriter(&buf).SetFormat("llm")
	if err := out.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	output := buf.String()

	// Primitive slices use pipe separator
	if !strings.Contains(output, "ITEMS: go | cli | tool") {
		t.Errorf("Expected pipe-separated items, got: %s", output)
	}
}

func TestLLMFormatter_Map(t *testing.T) {
	type Config struct {
		Settings map[string]string
	}

	data := Config{
		Settings: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	var buf bytes.Buffer
	out := NewOutput().SetWriter(&buf).SetFormat("llm")
	if err := out.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "SETTINGS:") {
		t.Errorf("Expected SETTINGS:, got: %s", output)
	}
	if !strings.Contains(output, "key1: value1") {
		t.Errorf("Expected key1: value1, got: %s", output)
	}
	if !strings.Contains(output, "key2: value2") {
		t.Errorf("Expected key2: value2, got: %s", output)
	}
}

func TestLLMFormatter_TimeFields(t *testing.T) {
	type Event struct {
		Name      string
		CreatedAt time.Time
		ZeroTime  time.Time
	}

	now := time.Date(2025, 12, 30, 10, 30, 0, 0, time.UTC)
	data := Event{
		Name:      "test-event",
		CreatedAt: now,
		ZeroTime:  time.Time{},
	}

	var buf bytes.Buffer
	out := NewOutput().SetWriter(&buf).SetFormat("llm")
	if err := out.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "CREATED_AT: 2025-12-30T10:30:00Z") {
		t.Errorf("Expected RFC3339 time, got: %s", output)
	}
	if strings.Contains(output, "ZERO_TIME") {
		t.Errorf("Should not contain ZERO_TIME, got: %s", output)
	}
}

func TestLLMFormatter_Duration(t *testing.T) {
	type Task struct {
		Name     string
		Duration time.Duration
	}

	data := Task{
		Name:     "build",
		Duration: 2*time.Hour + 30*time.Minute,
	}

	var buf bytes.Buffer
	out := NewOutput().SetWriter(&buf).SetFormat("llm")
	if err := out.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "DURATION: 2h30m0s") {
		t.Errorf("Expected duration string, got: %s", output)
	}
}

func TestLLMFormatter_Pointer(t *testing.T) {
	type Data struct {
		Name   string
		PtrVal *string
		NilPtr *string
		PtrInt *int
	}

	val := "pointer-value"
	num := 42
	data := Data{
		Name:   "test",
		PtrVal: &val,
		NilPtr: nil,
		PtrInt: &num,
	}

	var buf bytes.Buffer
	out := NewOutput().SetWriter(&buf).SetFormat("llm")
	if err := out.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "PTR_VAL: pointer-value") {
		t.Errorf("Expected PTR_VAL: pointer-value, got: %s", output)
	}
	if strings.Contains(output, "NIL_PTR") {
		t.Errorf("Should not contain NIL_PTR, got: %s", output)
	}
	if !strings.Contains(output, "PTR_INT: 42") {
		t.Errorf("Expected PTR_INT: 42, got: %s", output)
	}
}

func TestLLMFormatter_FieldNameConversion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Name", "NAME"},
		{"FieldName", "FIELD_NAME"},
		{"HTTPStatus", "HTTP_STATUS"},
		{"ID", "ID"},
		{"URLPath", "URL_PATH"},
		{"XMLParser", "XML_PARSER"},
	}

	f := &llmFormatter{}
	for _, tt := range tests {
		result := f.fieldNameToLabel(tt.input)
		if result != tt.expected {
			t.Errorf("fieldNameToLabel(%s) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func TestLLMFormatter_MaxDepth(t *testing.T) {
	type Deep struct {
		Level int
		Next  *Deep
	}

	// Create deeply nested structure
	data := &Deep{Level: 1}
	current := data
	for i := 2; i <= 10; i++ {
		current.Next = &Deep{Level: i}
		current = current.Next
	}

	var buf bytes.Buffer
	out := NewOutput().SetWriter(&buf).SetFormat("llm")
	if err := out.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	output := buf.String()

	// Should contain "..." for truncated deep levels
	if !strings.Contains(output, "...") {
		t.Errorf("Expected truncation marker ..., got: %s", output)
	}
}

func TestLLMFormatter_NilData(t *testing.T) {
	var buf bytes.Buffer
	out := NewOutput().SetWriter(&buf).SetFormat("llm")
	if err := out.Print(nil); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	if buf.Len() != 0 {
		t.Errorf("Expected empty output for nil, got: %s", buf.String())
	}
}

func TestLLMFormatter_EmptyStruct(t *testing.T) {
	type Empty struct {
		ZeroInt   int
		EmptyStr  string
		FalseBool bool
	}

	data := Empty{}

	var buf bytes.Buffer
	out := NewOutput().SetWriter(&buf).SetFormat("llm")
	if err := out.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	// All fields are zero values, should be empty output
	if buf.Len() != 0 {
		t.Errorf("Expected empty output for zero struct, got: %s", buf.String())
	}
}

func TestOutput_SetFormat_LLM(t *testing.T) {
	out := NewOutput().SetFormat("LLM") // uppercase

	type Data struct {
		Name string
	}

	var buf bytes.Buffer
	out.SetWriter(&buf)
	if err := out.Print(Data{Name: "test"}); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	// SetFormat converts to lowercase, so "LLM" should work
	if !strings.Contains(buf.String(), "NAME: test") {
		t.Errorf("LLM format should work with uppercase, got: %s", buf.String())
	}
}
