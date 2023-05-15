package syncutil_test

import (
	"testing"

	"github.com/teamwork/utils/v2/syncutil"
)

func TestMap_Store(t *testing.T) {
	t.Parallel()
	tests := map[string]int{
		"hello": 4,
		"world": 7,
		"ohwow": 10,
	}

	m := syncutil.NewMap[string, int]()
	for name, value := range tests {
		name := name
		value := value
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m.Store(name, value)

			got, ok := m.Get(name)
			if !ok {
				t.Fatalf("key '%s' expected 'true' got 'false'", name)
			}
			if value != got {
				t.Fatalf("key '%s' expected '%d' got '%d'", name, value, got)
			}
		})
	}
}

func TestMap_Store_Concurrent(t *testing.T) {
	t.Parallel()
	m := syncutil.NewMap[string, int]()
	for i := 0; i < 10000; i++ {
		i := i
		go func() {
			m.Store("hello", i)
		}()
	}
}

func TestMap_Range(t *testing.T) {
	t.Parallel()
	base := map[string]int{
		"hello":    1,
		"world":    2,
		"wow":      3,
		"holyhell": 4,
	}
	test := map[string]bool{
		"hello":    false,
		"world":    false,
		"wow":      false,
		"holyhell": false,
	}

	m := syncutil.NewMap[string, int]()
	for k, v := range base {
		m.Store(k, v)
	}

	m.Range(func(k string, v int) bool {
		seen, ok := test[k]
		if !ok {
			t.Fatalf("unexpected key '%s'", k)
		}

		if seen {
			t.Fatalf("duplicate range '%s'", k)
		}

		test[k] = true

		return true
	})

	for k, v := range test {
		if !v {
			t.Fatalf("missed key '%s'", k)
		}
	}
}

func TestMap_Range_EarlyReturn(t *testing.T) {
	t.Parallel()
	m := syncutil.NewMap[string, int]()
	m.Store("test", 1)
	m.Store("wow", 2)

	var count int

	m.Range(func(k string, v int) bool {
		count++
		return false
	})

	if count != 1 {
		t.Fatalf("unexpected count '%d'", count)
	}
}
