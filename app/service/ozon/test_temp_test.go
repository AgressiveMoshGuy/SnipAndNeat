package ozon

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

type OurVault struct {
	A *sync.Map
}

func NewOurVault() OurVault {
	var m sync.Map
	return OurVault{
		A: &m,
	}
}

func (n *OurVault) Push[K comparable, V int64|string](key K , value V) {
	n.A.Store(key, value)
}

func (n *OurVault) Pop(k string) (interface{}, error) {
	var v interface{}
	v, ok := n.A.Load(k)
	if !ok {
		return nil, errors.New("some error")
	}
	return v, nil
}

func main() {
	n := NewOurVault()
	n.Push("key", "value")
	n.Push("1", 2)
	res, err := n.Pop("1")
	fmt.Println(res)
	fmt.Println("run me", err)
}

func TestMain(t *testing.T) {
	n := NewOurVault()
	n.Push("key", "value")
	n.Push("1", 2)
	res, err := n.Pop("1")
	switch res.(type) {
	case int:
		fmt.Println("int")
	case string:
		fmt.Println("string")
	}
	if err != nil {
		t.Fatal(err)
	}
	if res != 2 {
		t.Fatalf("expected 2, got %v", res)
	}
}
