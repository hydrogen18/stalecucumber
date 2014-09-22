package stalecucumber

import "testing"
import "strings"
import "reflect"
import "math/big"
import "github.com/hydrogen18/stalecucumber/struct_export_test"

func TestUnpackMapIntoStructHiddenField(t *testing.T) {

	m := make(map[interface{}]interface{})
	m["shouldntMessWithMe"] = 2

	m = make(map[interface{}]interface{})
	m["likewise"] = 3
	s := struct_export_test.Struct1{}
	err := UnpackInto(&s).From(m, nil)
	if err == nil {
		t.Fatal("should have failed")
	}

	m = make(map[interface{}]interface{})
	m["likewise"] = 3

	err = UnpackInto(&s).From(m, nil)
	if err == nil {
		t.Fatal("should have failed")
	}
}

func TestUnpackIntIntoStruct(t *testing.T) {
	s := struct{}{}

	err := UnpackInto(&s).From(Unpickle(strings.NewReader("\x80\x02K\x00.")))
	if err == nil {
		t.Fatal("Should have failed!")
	}

	upe, ok := err.(UnpackingError)
	if !ok {
		t.Fatalf("Should have failed with type %T but got %T:%v", upe, err, err)
	}
}

const input0AsListOfDicts = "(lp0\n(dp1\nS'a'\np2\nL1L\nsS'c'\np3\nI3\nsS'b'\np4\nI2\nsa(dp5\ng2\nL1L\nsg3\nI3\nsg4\nI4\nsa(dp6\ng2\nL1L\nsg3\nI5\nsg4\nI2\nsa."

func TestUnpackListOfDictsIntoSliceOfStructs(t *testing.T) {
	dst := make([]testStruct, 0)
	expect := make([]testStruct, 3)
	expect[0] = testStruct{
		A: 1,
		B: 2,
		C: 3,
	}
	expect[1] = testStruct{
		A: 1,
		B: 4,
		C: 3,
	}
	expect[2] = testStruct{
		A: 1,
		B: 2,
		C: 5,
	}

	err := UnpackInto(&dst).From(Unpickle(strings.NewReader(input0AsListOfDicts)))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dst, expect) {
		t.Fatalf("Got %v expected %v", dst, expect)
	}
}

func TestUnpackListOfDictsIntoSliceOfPointersToStructs(t *testing.T) {
	dst := make([]*testStruct, 0)
	expect := make([]*testStruct, 3)
	expect[0] = &testStruct{
		A: 1,
		B: 2,
		C: 3,
	}
	expect[1] = &testStruct{
		A: 1,
		B: 4,
		C: 3,
	}
	expect[2] = &testStruct{
		A: 1,
		B: 2,
		C: 5,
	}

	err := UnpackInto(&dst).From(Unpickle(strings.NewReader(input0AsListOfDicts)))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dst, expect) {
		t.Fatalf("Got %v expected %v", dst, expect)
	}
}

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
const input0WithLong = "(dp0\nS'a'\np1\nL1L\nsS'c'\np2\nI3\nsS'b'\np3\nI2\ns."

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

	//Test with python long type in input. Generates *big.Int
	//with value 1
	dst = &testStruct{}

	err = UnpackInto(dst).From(Unpickle(strings.NewReader(input0WithLong)))
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

	//Test again w/ source having {"C": None }
	dst.A = 0
	dst.B = 0
	err = UnpackInto(dst).From(
		Unpickle(strings.NewReader("(dp0\nS'A'\np1\nI1\nsS'C'\np2\nNsS'B'\np3\nI2\ns.")))
	if err != nil {
		t.Fatal(err)
	}
	expect.C = nil
	expect.A = 1
	expect.B = 2

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
const inputDWithUnicode = "\x80\x02}q\x00(X\x08\x00\x00\x00Aardvarkq\x01K\x01U\x05Bolusq\x02G@\x08\x00\x00\x00\x00\x00\x00U\x03Catq\x03}q\x04(U\x05appleq\x05K\x02X\x06\x00\x00\x00bananaq\x06K\x03uu."

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

type testStructDWithTags struct {
	One   uint    `pickle:"Aardvark"`
	Two   float32 `pickle:"Bolus"`
	Three struct {
		Four int  `pickle:"apple"`
		Five uint `pickle:"banana"`
	} `pickle:"Cat"`
}

func TestStructDWithPickleNames(t *testing.T) {
	dst := &testStructDWithTags{}
	expect := &testStructDWithTags{
		One: 1,
		Two: 3.0,
	}
	expect.Three.Four = 2
	expect.Three.Five = 3

	err := UnpackInto(dst).From(Unpickle(strings.NewReader(inputD)))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dst, expect) {
		t.Fatalf("Got %v expected %v", *dst, *expect)
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

	dst = &testStructDWithStruct{}
	err = UnpackInto(dst).From(Unpickle(strings.NewReader(inputDWithUnicode)))
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

type testStructDWithBadStruct struct {
	Aardvark uint
	Bolus    float32
	Cat      struct {
		Apple  string
		Banana uint
	}
}

func TestUnpackStructDWithBadStruct(t *testing.T) {
	dst := &testStructDWithBadStruct{}

	err := UnpackInto(dst).From(Unpickle(strings.NewReader(inputD)))
	if err == nil {
		t.Fatalf("Should not have unpacked:%v", dst)
	}
}

const inputE = "(dp0\nS'ds'\np1\n(lp2\n(dp3\nS'a'\np4\nL1L\nsS'c'\np5\nI3\nsS'b'\np6\nI2\nsa(dp7\ng4\nL1L\nsg5\nI3\nsg6\nI4\nsa(dp8\ng4\nL1L\nsg5\nI5\nsg6\nI2\nsas."

type testStructureE struct {
	Ds []testStruct
}

func TestUnpackDictWithListOfDictsIntoStructWithListOfDicts(t *testing.T) {
	dst := testStructureE{}
	e := testStructureE{}
	e.Ds = make([]testStruct, 3)
	e.Ds[0] = testStruct{
		A: 1,
		B: 2,
		C: 3,
	}
	e.Ds[1] = testStruct{
		A: 1,
		B: 4,
		C: 3,
	}
	e.Ds[2] = testStruct{
		A: 1,
		B: 2,
		C: 5,
	}

	err := UnpackInto(&dst).From(Unpickle(strings.NewReader(inputE)))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dst, e) {
		t.Fatalf("Got %v expected %v", dst, e)
	}

}

const inputF = "(lp0\nI0\naI1\naI2\naI3\naI4\na."

func TestUnpackSliceOfInts(t *testing.T) {
	dst := make([]int64, 0)
	expect := []int64{0, 1, 2, 3, 4}

	err := UnpackInto(&dst).From(Unpickle(strings.NewReader(inputF)))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dst, expect) {
		t.Fatalf("Got %v expected %v", dst, expect)
	}
}

const inputG = "(lp0\nS'foo'\np1\naVbar\np2\naS'qux'\np3\na."

func TestUnpackSliceOfStrings(t *testing.T) {
	dst := []string{"disappears"}
	expect := []string{"foo", "bar", "qux"}

	err := UnpackInto(&dst).From(Unpickle(strings.NewReader(inputG)))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dst, expect) {
		t.Fatalf("Got %v expected %v", dst, expect)
	}
}

const inputH = "(lp0\nS'meow'\np1\naI42\naS'awesome'\np2\na."

func TestUnpackHeterogeneousList(t *testing.T) {
	dst := []interface{}{}
	expect := []interface{}{"meow", int64(42), "awesome"}

	err := UnpackInto(&dst).From(Unpickle(strings.NewReader(inputH)))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dst, expect) {
		t.Fatalf("Got %v expected %v", dst, expect)
	}

	dst2 := []string{}

	err = UnpackInto(&dst2).From(Unpickle(strings.NewReader(inputH)))
	if err == nil {
		t.Fatal(err)
	}

	upe, ok := err.(UnpackingError)
	if !ok {
		t.Fatalf("Got wrong error type %T:%v", err, err)
	}

	i, ok := upe.Source.(int64)
	if !ok && i == 42 {
		t.Fatalf("Failed on wrong value %v(%T)", upe.Source, upe.Source)
	}

}
