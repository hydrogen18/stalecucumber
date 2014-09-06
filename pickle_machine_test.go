package stalecucumber

import "testing"
import "strings"
import "math/big"

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

func TestProtocol0Long(t *testing.T) {
	result := new(big.Int)
	reader := strings.NewReader("L5L\n.")
	expect := big.NewInt(5)
	err := Unmarshal(reader, &result)
	if err != nil {
		t.Fatalf("Got error %v", err)
	}

	if result.Cmp(expect) != 0 {
		t.Fatalf("Got value %q expected %q", result, expect)
	}
}

func TestProtocol0Float(t *testing.T) {
	var result float64
	reader := strings.NewReader("F3.14\n.")
	const EXPECT = 3.14

	err := Unmarshal(reader, &result)
	if err != nil {
		t.Fatalf("Got error %v", err)
	}

	if EXPECT != result {
		t.Fatalf("Got value %q expected %q", result, EXPECT)
	}
}
func TestProtocol0Dict(t *testing.T) {
	reader := strings.NewReader("(dp0\nS'a'\np1\nI1\nsS'b'\np2\nI5\ns.")
	var result map[interface{}]interface{}
	expect := make(map[string]int64)
	expect["a"] = 1
	expect["b"] = 5

	err := Unmarshal(reader, &result)
	if err != nil {
		t.Fatalf("Got error %v", err)
	}
	if len(result) != len(expect) {
		t.Errorf("result has wrong length %d", len(result))
	}

	for k, v := range result {

		var expectedV int64
		if kstr, ok := k.(string); ok {
			expectedV, ok = expect[kstr]
			if !ok {
				t.Errorf("key %q not found in expectation", kstr)
				continue
			}
		} else {
			t.Errorf("key %v has unexpected type %T, not %T", k, k, kstr)
			continue
		}

		if vint, ok := v.(int64); ok {
			if vint != expectedV {
				t.Errorf("result[%q] has unexpected value %d not %d", k, vint, expectedV)
			}
		} else {
			t.Errorf("result[%q] has unexpected type %T not %T", k, v, expectedV)
		}
	}
}

func testList(t *testing.T, input string, expect []int64) {
	var result []interface{}
	reader := strings.NewReader(input)

	err := Unmarshal(reader, &result)
	if err != nil {
		t.Fatalf("Got error %v", err)
	}
	if len(result) != len(expect) {
		t.Errorf("Result has wrong length %d", len(result))
	}
	for i, v := range result {
		var vi int64
		var ok bool
		vi, ok = v.(int64)
		if !ok {
			t.Errorf("result[%d]=%v not type %T", i, v, vi)
			continue
		}

		if vi != expect[i] {
			t.Errorf("result[%d] != expect[%d]", i, i)
		}
	}

}

func TestProtocol0List(t *testing.T) {
	testList(t, "(lp0\nI1\naI2\naI3\na.", []int64{1, 2, 3})
}

func TestProtocol1List(t *testing.T) {
	testList(t, "]q\x00.", []int64{})
	testList(t, "]q\x00(M9\x05M9\x05M9\x05e.", []int64{1337, 1337, 1337})
	testList(t, "]q\x00(M9\x05I3735928559\nM\xb1\"e.", []int64{1337, 0xdeadbeef, 8881})
}

func testInt(t *testing.T, input string, expect int64) {
	var result int64
	reader := strings.NewReader(input)

	err := Unmarshal(reader, &result)
	if err != nil {
		t.Fatalf("Got error %v", err)
	}
	if result != expect {
		t.Fatalf("Got %d(%T) expected %d(%T)", result, result, expect, expect)
	}

}

func TestProtocol1Binint(t *testing.T) {
	testInt(t, "J\xff\xff\xff\x00.", 0xffffff)
	testInt(t, "K*.", 42)
	testInt(t, "M\xff\xab.", 0xabff)
}

func testString(t *testing.T, input string, expect string) {
	var result string
	reader := strings.NewReader(input)

	err := Unmarshal(reader, &result)

	if err != nil {
		t.Fatalf("Got error %v", err)
	}

	if result != expect {
		t.Fatalf("Got %q(%T) expected %q(%T)", result, result, expect, expect)
	}
}

func TestProtocol1String(t *testing.T) {
	testString(t,
		"T\x04\x01\x00\x00abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZq\x00.",
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	testString(t, "U\x13queen of the castleq\x00.", "queen of the castle")
}

func TestProtocol1Float(t *testing.T) {
	var result float64
	reader := strings.NewReader("G?\xc1\x1d\x14\xe3\xbc\xd3[.")

	err := Unmarshal(reader, &result)

	if err != nil {
		t.Fatalf("Got error %v", err)
	}
	var expect float64
	expect = 0.1337
	if result != expect {
		t.Fatalf("Got %f expected %f", result, expect)
	}
}
