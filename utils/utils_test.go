package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestHash(t *testing.T) {
	hash := "e005c1d727f7776a57a661d61a182816d8953c0432780beeae35e337830b1746"
	s := struct{ Test string }{Test: "test"}
	// sub-test 1
	t.Run("Hash is always same", func(t *testing.T) {
		x := Hash(s)
		if x != hash {
			t.Errorf("Expected %s, got %s\n", hash, x)
		}
	})
	// sub-test 2
	t.Run("Hahsh is hex encoded", func(t *testing.T) {
		x := Hash(s)
		_, err := hex.DecodeString(x)
		if err != nil {
			t.Error("Hash should be hex encoded")
		}
	})
}

// godoc Example 표기 양식
func ExampleHash() {
	s := struct{ Test string }{Test: "test"}
	x := Hash(s)
	fmt.Println(x)
	// Output: e005c1d727f7776a57a661d61a182816d8953c0432780beeae35e337830b1746
}

func TestToBytes(t *testing.T) {
	bys := []byte{7, 12, 0, 4, 116, 101, 115, 116}
	s := "test"
	bs := ToBytes(s)
	// type check
	kind := reflect.TypeOf(bs).Kind()
	if kind != reflect.Slice {
		t.Errorf("ToBytes should return a slice of bytes got %s\n", kind)
	}
	for i, b := range bs {
		if b != bys[i] {
			t.Errorf("Expected %d, got %d", bys[i], b)
		}
	}
}

// table test
func TestSplitter(t *testing.T) {
	type test struct {
		input  string
		sep    string
		index  int
		output string
	}
	tests := []test{
		{"0:5:0", ":", 1, "5"},
		{"0:5:0", ":", 10, ""},
		{"0:5:0", "-", 0, "0:5:0"},
	}
	for _, tc := range tests {
		got := Splitter(tc.input, tc.sep, tc.index)
		if got != tc.output {
			t.Errorf("Expected %s and got %s", tc.output, got)
		}
	}
}

// 콘솔 출력, 로그 출력을 테스트 하는 방법.
// 출력 부분을 별도 함수로 빼고, 그 함수를 체크할 수 있는 함수로 잠시 바꿔서 테스트함
func TestHandleErr(t *testing.T) {
	oldLogFn := logFn
	defer func() {
		logFn = oldLogFn
	}()
	called := false
	logFn = func(v ...interface{}) {
		called = true
	}
	err := errors.New("test")
	HandleErr(err)
	if !called {
		t.Error("HandleError should call fn")
	}
}
