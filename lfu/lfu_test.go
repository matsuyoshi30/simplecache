package lfu_test

import (
	"testing"

	"github.com/matsuyoshi30/simplecache/lfu"
)

func TestLFU(t *testing.T) {
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
			l := lfu.NewLFU(tt.cap)
			for i := 0; i < tt.cap+1; i++ {
				err := l.Put(i, i) // test LFU.Put
				if err != nil {
					t.Fatal(err)
				}
			}

			_, err := l.Get(0) // test capacity
			if err != lfu.ErrNotFound {
				t.Errorf("want %v but got %v\n", lfu.ErrNotFound, err)
			}

			v, err := l.Get(1) // test get
			if err != nil {
				t.Errorf("want no error but got %v\n", err)
			}
			if v != 1 {
				t.Errorf("want 1 but got %v\n", v)
			}

			if tt.cap > 1 {
				// update odd entry
				for i := 0; i < tt.cap; i++ {
					if i%2 == 1 {
						_, err := l.Get(i)
						if err != nil {
							t.Errorf("want no error but got %v\n", err)
						}
					}
				}
				err := l.Put(tt.cap+1, tt.cap+1)
				if err != nil {
					t.Fatal(err)
				}

				_, err = l.Get(2) // expect not found oldest and minimum references node
				if err != lfu.ErrNotFound {
					t.Errorf("want %v but got %v\n", lfu.ErrNotFound, err)
				}
			}
		})
	}
}

func TestLFU_PutExistKey(t *testing.T) {
	l := lfu.NewLFU(2)
	for i := 0; i < 2; i++ {
		err := l.Put(i, i)
		if err != nil {
			t.Fatal(err)
		}
	}

	err := l.Put(1, 10)
	if err != nil {
		t.Errorf("want no error but got %v\n", err)
	}

	v, err := l.Get(1)
	if err != nil {
		t.Errorf("want no error but got %v\n", err)
	}
	if v != 10 {
		t.Errorf("want 10 but got %v\n", v)
	}
}
