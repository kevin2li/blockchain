package pkg

import (
	"fmt"

	"github.com/pkg/errors"
)

/* HashSet */
type HashSet map[string]bool

func (h HashSet) Len() int {
	return len(h)
}

func (h HashSet) Add(s string) {
	if _, ok := h[s]; !ok {
		h[s] = true
	}
}

func (h HashSet) Remove(s string) error {
	if _, ok := h[s]; !ok {
		return errors.New(fmt.Sprintf("KeyError: Key `%s` does not exist!", s))
	}
	delete(h, s)
	return nil
}

func (h HashSet) GetData() []string {
	var result []string
	for k := range h {
		result = append(result, k)
	}
	return result
}

func (h HashSet) Init(list []string) {
	for _, v := range list {
		h.Add(v)
	}
}

// remove duplicate addrs
func Unique(addrs []string) []string {
	set := make(HashSet)
	for _, addr := range addrs {
		set.Add(addr)
	}
	result := set.GetData()
	return result
}
