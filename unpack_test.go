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
	B int64
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
	D *big.Int
	E bool
}

func TestUnpackStructB(t *testing.T) {
	dst := &testStructB{}
	expect := &testStructB{
		A: 42,
		B: 13.37,
		C: "foobar",
		D: big.NewInt(1),
		E: true,
	}

	err := UnpackInto(dst).From(Unpickle(strings.NewReader(inputB)))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dst, expect) {
		t.Fatalf("Got %v expected %v", *dst, *expect)
	}
}
