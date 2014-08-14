package stalecucumber

import "testing"
import "strings"

func TestProtocol0Integer(t *testing.T) {
	var result int64
	reader := strings.NewReader("I42\n.")
	err := Unmarshal(reader, &result)
	if err != nil {
		t.Fatalf("Got error %v", err)
	}
	const EXPECT = 42
	if result != EXPECT {
		t.Fatalf("Got value %d expected %d", result, EXPECT)
	}
}

func TestProtocol0Bool(t *testing.T) {
	var result bool

	reader := strings.NewReader("I00\n.")
	err := Unmarshal(reader, &result)
	if err != nil {
		t.Fatalf("Got error %v", err)
	}

	if result != false {
		t.Fatalf("Got value %v expected %v", result, false)
	}

	reader = strings.NewReader("I01\n.")
	err = Unmarshal(reader, &result)
	if err != nil {
		t.Fatalf("Got error %v", err)
	}

	if result != true {
		t.Fatalf("Got value %v expected %v", result, true)
	}

}

func TestProtocol0String(t *testing.T) {
	var result string
	reader := strings.NewReader("S'foobar'\np0\n.")
	const EXPECT = "foobar"

	err := Unmarshal(reader, &result)
	if err != nil {
		t.Fatalf("Got error %v", err)
	}

	if result != EXPECT {
		t.Fatalf("Got value %q expected %q", result, EXPECT)
	}
}
