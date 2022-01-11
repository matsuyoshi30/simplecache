package lru_test

import (
	"testing"

	"github.com/matsuyoshi30/simplecache/lru"
)

func TestLRU(t *testing.T) {
	tests := []struct {
		desc string
		cap  int
	}{
		{
			desc: "cap 1",
			cap:  1,
		},
		{
			desc: "cap 10",
			cap:  10,
		},
		{
			desc: "cap 100",
			cap:  100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			l := lru.NewLRU(tt.cap)
			for i := 0; i < tt.cap+1; i++ {
				err := l.Put(i, i) // test LRU.Put
				if err != nil {
					t.Fatal(err)
				}
			}

			_, err := l.Get(0) // test capacity
			if err != lru.ErrNotFound {
				t.Errorf("want %v but got %v\n", lru.ErrNotFound, err)
			}

			v, err := l.Get(1) // test get
			if err != nil {
				t.Errorf("want no error but got %v\n", err)
			}
			if v != 1 {
				t.Errorf("want 1 but got %v\n", v)
			}
		})
	}
}
