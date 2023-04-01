package toolkit

import "testing"

// TestTools_RandomString tests the RandomString function
func TestTools_RandomString(t *testing.T) {
	var testTools Tools
	s := testTools.RandomString(10)
	if len(s) != 10 {
		t.Errorf("Expected 10, got %d", len(s))
	}

}
