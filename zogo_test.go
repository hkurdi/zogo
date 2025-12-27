package zogo

import "testing"

func TestParseResultSuccess(t *testing.T) {
	result := Success("test value")

	if !result.Ok {
		t.Error("Expected Ok to be true")
	}

	if result.Value != "test value" {
		t.Errorf("Expected value 'test value', got %v", result.Value)
	}

	if len(result.Errors) != 0 {
		t.Error("Expected no errors")
	}
}

func TestParseResultFailure(t *testing.T) {
	result := FailureMessage("test error")

	if result.Ok {
		t.Error("Expected Ok to be false")
	}

	if len(result.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(result.Errors))
	}

	if result.Errors[0].Message != "test error" {
		t.Errorf("Expected message 'test error', got %s", result.Errors[0].Message)
	}
}
