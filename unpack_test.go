package stalecucumber

import "testing"
import "strings"
import "reflect"
import "math/big"

type testStruct struct {
	A int64
	B int64
	C int64
}

type testStructWithPointer struct {
	A int64
	B uint64
	C *int64
}

const input0 = "\x80\x02}q\x00(U\x01aq\x01K\x01U\x01cq\x02K\x03U\x01bq\x03K\x02u."

func TestUnpackIntoStruct(t *testing.T) {
	dst := &testStruct{}
	expect := &testStruct{
		A: 1,
		B: 2,
		C: 3,
	}

	err := UnpackInto(dst).From(Unpickle(strings.NewReader(input0)))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dst, expect) {
		t.Fatalf("Got %v expected %v", *dst, *expect)
	}
}

func TestUnpackIntoStructWithPointer(t *testing.T) {
	dst := &testStructWithPointer{}
	expect := &testStructWithPointer{
		A: 1,
		B: 2,
		C: new(int64),
	}
	*expect.C = 3

	err := UnpackInto(dst).From(Unpickle(strings.NewReader(input0)))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dst, expect) {
		t.Fatalf("Got %v expected %v", *dst, *expect)
	}

	//Test again w/ dst.C non-nil
	dst.A = 0
	dst.B = 0
	dst.C = new(int64)
	*dst.C = 1337

	err = UnpackInto(dst).From(Unpickle(strings.NewReader(input0)))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dst, expect) {
		t.Fatalf("Got %v expected %v", *dst, *expect)
	}

}

const inputB = "\x80\x02}q\x00(U\x01aq\x01K*U\x01cq\x02U\x06foobarq\x03U\x01bq\x04G@*\xbdp\xa3\xd7\n=U\x01eq\x05\x88U\x01dq\x06\x8a\x01\x01u."

type testStructB struct {
	A int
	B float32
	C string
	D big.Int
	E bool
}

type testStructBWithPointers struct {
	A *int
	B *float32
	C *string
	D *big.Int
	E *bool
}

func TestUnpackStructB(t *testing.T) {
	dst := &testStructB{}
	expect := &testStructB{
		A: 42,
		B: 13.37,
		C: "foobar",
		D: *big.NewInt(1),
		E: true,
	}

	err := UnpackInto(dst).From(Unpickle(strings.NewReader(inputB)))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dst, expect) {
		t.Fatalf("Got %v expected %v", *dst, *expect)
	}

	dstP := &testStructBWithPointers{}

	err = UnpackInto(dstP).From(Unpickle(strings.NewReader(inputB)))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dst, expect) {
		t.Fatalf("Got %v expected %v", *dst, *expect)
	}
}

const inputC = "\x80\x02}q\x00(U\x03dogq\x01U\x01aq\x02U\x01bq\x03U\x01cq\x04\x87q\x05U\x05appleq\x06K\x01K\x02K\x03\x87q\x07U\ncanteloupeq\x08h\x05U\x06bananaq\th\x07u."

type testStructC struct {
	Apple      []interface{}
	Banana     []interface{}
	Canteloupe []interface{}
	Dog        []interface{}
}

func TestUnpackStructC(t *testing.T) {
	dst := &testStructC{}
	expect := &testStructC{
		Apple:      []interface{}{int64(1), int64(2), int64(3)},
		Banana:     []interface{}{int64(1), int64(2), int64(3)},
		Canteloupe: []interface{}{"a", "b", "c"},
		Dog:        []interface{}{"a", "b", "c"},
	}

	err := UnpackInto(dst).From(Unpickle(strings.NewReader(inputC)))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dst, expect) {
		t.Fatalf("Got %v expected %v", *dst, *expect)
	}
}

const inputD = "\x80\x02}q\x00(U\x08Aardvarkq\x01K\x01U\x05Bolusq\x02G@\x08\x00\x00\x00\x00\x00\x00U\x03Catq\x03}q\x04(U\x05appleq\x05K\x02U\x06bananaq\x06K\x03uu."

type testStructDWithMap struct {
	Aardvark uint
	Bolus    float32
	Cat      map[interface{}]interface{}
}

type testStructDWithStruct struct {
	Aardvark uint
	Bolus    float32
	Cat      struct {
		Apple  int
		Banana uint
	}
}

func TestUnpackStructDWithStruct(t *testing.T) {
	dst := &testStructDWithStruct{}
	expect := &testStructDWithStruct{
		Aardvark: 1,
		Bolus:    3.0,
	}
	expect.Cat.Apple = 2
	expect.Cat.Banana = 3

	err := UnpackInto(dst).From(Unpickle(strings.NewReader(inputD)))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dst, expect) {
		t.Fatalf("Got %v expected %v", *dst, *expect)
	}
}

func TestUnpackStructDWithMap(t *testing.T) {
	dst := &testStructDWithMap{}
	expect := &testStructDWithMap{
		Aardvark: 1,
		Bolus:    3.0,
		Cat:      make(map[interface{}]interface{}),
	}
	expect.Cat["apple"] = int64(2)
	expect.Cat["banana"] = int64(3)

	err := UnpackInto(dst).From(Unpickle(strings.NewReader(inputD)))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dst, expect) {
		t.Fatalf("Got %v expected %v", *dst, *expect)
	}
}
