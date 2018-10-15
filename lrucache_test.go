package lrucache

import (
	"testing"
)

func TestBasicFunctions(t *testing.T) {
	testdata := map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
		"f": 6,
	}
	// Check Lrucache initialization with NewLrucache
	cache, err := NewLrucache(-1)
	if err == nil {
		t.Fatal("NewLrucache(-1) does not return error")
	}
	cache, err = NewLrucache(0)
	if err == nil {
		t.Fatal("NewLrucache(0) does not return error")
	}
	cache, err = NewLrucache(len(testdata) - 1)
	if err != nil {
		t.Fatalf("NewLrucache(%d) returns error: %v", len(testdata)-1, err)
	}
	for k, v := range testdata {
		cache.Set(k, v)
	}
	if cache.Len() != len(testdata)-1 {
		t.Error("Cache displacement does not work")
	}
	// Check Get()
	val, ok := cache.Get("d")
	if !ok {
		t.Error("Get method is not working")
	}
	if val != 4 {
		t.Error("Get method returns wrong vakue")
	}
	// Check Del()
	cache.Del("d")
	_, ok = cache.Get("d")
	if ok {
		t.Error("Del method is not working")
	}
	// Check Flush()
	cache.Flush()
	if cache.Len() != 0 {
		t.Error("Flush method is not working")
	}

}
