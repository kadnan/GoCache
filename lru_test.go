package gCache

import (
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var (
	testCache = &Cache{
		capacity: 5,
		items: map[string]*cacheItem{
			"First": &cacheItem{
				value:   "value",
				lastUse: time.Now().UnixNano(),
			},
			"Second": &cacheItem{
				value:   "Value2",
				lastUse: time.Now().UnixNano(),
			},
		},
	}
)

// Diff two map[string]*cacheItem maps without testing the lastUse values.
func diffNoLU(A, B map[string]*cacheItem) bool {
	trans := cmp.Transformer("Sort", func(in []string) []string {
		out := append([]string(nil), in...) // Copy input to avoid mutating it
		sort.Strings(out)
		return out
	})

	// If the lengths don't match, fail.
	if len(A) != len(B) {
		return false
	}

	// If there are not the same keys fail.
	var aKeys, bKeys []string
	for k, _ := range A {
		aKeys = append(aKeys, k)
	}
	for k, _ := range B {
		bKeys = append(bKeys, k)
	}
	if !cmp.Equal(aKeys, bKeys, trans) {
		return false
	}

	// Compare key/value pairs between the 2 maps.
	for k, _ := range A {
		if A[k].value != B[k].value {
			return false
		}
	}
	return true
}

func TestSet(t *testing.T) {
	tests := []struct {
		desc string
		cap  int
		sets []string
		want *Cache
	}{{
		desc: "Success Set 1 k/v",
		cap:  2,
		sets: []string{"First|value"},
		want: &Cache{
			capacity: 5,
			items: map[string]*cacheItem{
				"First": &cacheItem{
					value:   "value",
					lastUse: time.Now().UnixNano(),
				},
			},
		},
	}, {
		desc: "Success Set 3 k/v leave 2",
		cap:  2,
		sets: []string{"First|value", "Second|Value2", "Third|Value3"},
		want: &Cache{
			capacity: 5,
			items: map[string]*cacheItem{
				"Second": &cacheItem{
					value:   "Value2",
					lastUse: time.Now().UnixNano(),
				},
				"Third": &cacheItem{
					value:   "Value3",
					lastUse: time.Now().UnixNano(),
				},
			},
		},
	}}

	for _, test := range tests {
		got := New(test.cap)
		for _, kv := range test.sets {
			splits := strings.Split(kv, "|")
			got.Set(splits[0], splits[1])
		}
		if !cmp.Equal(got.items, test.want.items, cmp.Comparer(diffNoLU)) {
			t.Errorf("[%v]: got/want mismatch (+got/-want):\n%v\n", test.desc, cmp.Diff(got.items, test.want.items, cmp.AllowUnexported(cacheItem{})))
		}
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		desc    string
		cache   *Cache
		get     string
		want    string
		wantErr bool
	}{{
		desc:  "Success Get First",
		cache: testCache,
		get:   "First",
		want:  "value",
	}, {
		desc:    "Error non-existent key",
		cache:   testCache,
		get:     "Third",
		want:    "zzz",
		wantErr: true,
	}}

	for _, test := range tests {
		got, err := test.cache.Get(test.get)
		switch {
		case err != nil && !test.wantErr:
			t.Errorf("[%v]: got error when not expecting one: %v", test.desc, err)
		case err == nil && test.wantErr:
			t.Errorf("[%v]: got no error, but expected one", test.desc)
		case err == nil:
			if got != test.want {
				t.Errorf("[%v]: got/want mismatch: %v / %v", test.desc, got, test.want)
			}
		}
	}
}
