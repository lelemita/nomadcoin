package utils

import (
	"encoding/hex"
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
