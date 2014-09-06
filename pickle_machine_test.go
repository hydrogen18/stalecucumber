package stalecucumber

import "testing"
import "strings"
import "math/big"
import "reflect"

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

func testDict(t *testing.T, input string, expect map[interface{}]interface{}) {
	reader := strings.NewReader(input)
	var result map[interface{}]interface{}

	err := Unmarshal(reader, &result)
	if err != nil {
		t.Fatalf("Got error %v", err)
	}
	if len(result) != len(expect) {
		t.Errorf("result has wrong length %d", len(result))
	}

	for k, v := range result {
		var expectedV interface{}

		expectedV, ok := expect[k]
		if !ok {
			t.Errorf("Result has key %v(%T) which is not in expectation", k, k)
			continue
		}

		if reflect.TypeOf(v) != reflect.TypeOf(expectedV) {
			t.Errorf("At key %v result has type %T where expectation has type %T", k, v, expectedV)
			continue
		}

		if !reflect.DeepEqual(expectedV, v) {
			t.Errorf("At key %v result %v != expectation %v", k, v, expectedV)
		}

	}
}

func TestProtocol0Get(t *testing.T) {
	testList(t, "(lp0\nS'hydrogen18'\np1\nag1\na.", []interface{}{"hydrogen18", "hydrogen18"})
}

func TestProtocol1Get(t *testing.T) {
	testList(t, "]q\x00(U\nhydrogen18q\x01h\x01e.", []interface{}{"hydrogen18", "hydrogen18"})
}

func TestProtocol0Dict(t *testing.T) {

	{
		input := "(dp0\nS'a'\np1\nI1\nsS'b'\np2\nI5\ns."
		expect := make(map[interface{}]interface{})
		expect["a"] = int64(1)
		expect["b"] = int64(5)
		testDict(t, input, expect)
	}

	{
		expect := make(map[interface{}]interface{})
		expect["foo"] = "bar"
		expect[int64(5)] = "kitty"
		expect["num"] = 13.37
		expect["list"] = []interface{}{int64(1), int64(2), int64(3), int64(4)}
		testDict(t, "(dp0\nS'list'\np1\n(lp2\nI1\naI2\naI3\naI4\nasS'foo'\np3\nS'bar'\np4\nsS'num'\np5\nF13.37\nsI5\nS'kitty'\np6\ns.", expect)
	}

}

func testList(t *testing.T, input string, expect []interface{}) {
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

		vexpect := expect[i]

		if !reflect.DeepEqual(v, vexpect) {
			t.Errorf("result[%v(%T)] != expect[%v(%T)]", i, v, i, vexpect)
		}
	}

}

func TestProtocol0List(t *testing.T) {
	testList(t, "(lp0\nI1\naI2\naI3\na.", []interface{}{int64(1), int64(2), int64(3)})
}

func TestProtocol1List(t *testing.T) {
	testList(t, "]q\x00.", []interface{}{})
	testList(t, "]q\x00(M9\x05M9\x05M9\x05e.", []interface{}{int64(1337), int64(1337), int64(1337)})
	testList(t, "]q\x00(M9\x05I3735928559\nM\xb1\"e.", []interface{}{int64(1337), int64(0xdeadbeef), int64(8881)})
}

func TestProtocol1Tuple(t *testing.T) {
	testList(t, ").", []interface{}{})
	testList(t, "(K*K\x18K*K\x1cKRK\x1ctq\x00.", []interface{}{int64(42), int64(24), int64(42), int64(28), int64(82), int64(28)})
}

func TestProtocol1Dict(t *testing.T) {

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
