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
	testString(t, "S'foobar'\np0\n.", "foobar")
	testString(t, "S'String with embedded\\nnewline.'\np0\n.", "String with embedded\nnewline.")
	testString(t,
		"\x53\x27\x53\x74\x72\x69\x6e\x67\x20\x77\x69\x74\x68\x20\x65\x6d\x62\x65\x64\x64\x65\x64\x5c\x6e\x6e\x65\x77\x6c\x69\x6e\x65\x20\x61\x6e\x64\x20\x65\x6d\x62\x65\x64\x64\x65\x64\x20\x71\x75\x6f\x74\x65\x20\x5c\x27\x20\x61\x6e\x64\x20\x65\x6d\x62\x65\x64\x64\x65\x64\x20\x64\x6f\x75\x62\x6c\x65\x71\x75\x6f\x74\x65\x20\x22\x2e\x27\x0a\x70\x30\x0a\x2e",
		"String with embedded\nnewline and embedded quote ' and embedded doublequote \".")
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

func TestProtocol1Dict(t *testing.T) {
	testDict(t, "}q\x00.", make(map[interface{}]interface{}))
	{
		expect := make(map[interface{}]interface{})
		expect["foo"] = "bar"
		expect["meow"] = "bar"
		expect[int64(5)] = "kitty"
		expect["num"] = 13.37
		expect["list"] = []interface{}{int64(1), int64(2), int64(3), int64(4)}
		input := "}q\x00(U\x04meowq\x01U\x03barq\x02U\x04listq\x03]q\x04(K\x01K\x02K\x03K\x04eU\x03fooq\x05h\x02U\x03numq\x06G@*\xbdp\xa3\xd7\n=K\x05U\x05kittyq\x07u."
		testDict(t, input, expect)
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

func TestProtocol1PopMark(t *testing.T) {
	var result int64
	/**
		This exapmle is ultra-contrived. I could not get anything to
		produce usage of POP_MARK using protocol 1. There are some
		comments in Lib/pickle.py about a recursive tuple generating
		this but I have no idea how that is even possible.

		The disassembly of this looks like
	    0: K    BININT1    1
	    2: (    MARK
	    3: K        BININT1    2
	    5: K        BININT1    3
	    7: 1        POP_MARK   (MARK at 2)
	    8: .    STOP

		There is just a mark placed on the stack with some numbers
		afterwards solely to test the correct behavior
		of the POP_MARK instruction.

		**/
	reader := strings.NewReader("K\x01(K\x02K\x031.")
	err := Unmarshal(reader, &result)
	const EXPECT = 1
	if err != nil {
		t.Fatalf("Got error %v", err)
	}

	if EXPECT != result {
		t.Fatalf("Got %d expected %d", result, EXPECT)
	}
}

func TestProtocol1Unicode(t *testing.T) {
	expect := "This is a slash \\. This is a newline \n. This is a character that is two embedded newlines: \u0a0a. This is a snowman: \u2603."

	if len([]rune(expect)) != 115 {
		t.Errorf("Expect shouldn't be :%v", expect)
		t.Fatalf("you messed up the escape sequence on the expecation, again. Length is %d", len(expect))
	}
	testString(t, "\x56\x54\x68\x69\x73\x20\x69\x73\x20\x61\x20\x73\x6c\x61\x73\x68\x20\x5c\x75\x30\x30\x35\x63\x2e\x20\x54\x68\x69\x73\x20\x69\x73\x20\x61\x20\x6e\x65\x77\x6c\x69\x6e\x65\x20\x5c\x75\x30\x30\x30\x61\x2e\x20\x54\x68\x69\x73\x20\x69\x73\x20\x61\x20\x63\x68\x61\x72\x61\x63\x74\x65\x72\x20\x74\x68\x61\x74\x20\x69\x73\x20\x74\x77\x6f\x20\x65\x6d\x62\x65\x64\x64\x65\x64\x20\x6e\x65\x77\x6c\x69\x6e\x65\x73\x3a\x20\x5c\x75\x30\x61\x30\x61\x2e\x20\x54\x68\x69\x73\x20\x69\x73\x20\x61\x20\x73\x6e\x6f\x77\x6d\x61\x6e\x3a\x20\x5c\x75\x32\x36\x30\x33\x2e\x0a\x70\x30\x0a\x2e",
		expect)

}
