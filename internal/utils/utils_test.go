package utils

import (
	"github.com/fond-of-vertigo/amazon-sp-api/constants"
	"github.com/google/go-cmp/cmp"
	"net/url"
	"reflect"
	"testing"
)

func TestFirstNElementsOfSlice(t *testing.T) {
	type args[Element any] struct {
		slice []Element
		n     int
	}
	type testCase[Element any] struct {
		name string
		args args[Element]
		want []Element
	}
	tests := []testCase[string]{
		{
			name: "empty slice",
			args: args[string]{
				slice: []string{},
				n:     1,
			},
			want: []string{},
		},
		{
			name: "slice with 1 element",
			args: args[string]{
				slice: []string{"a"},
				n:     1,
			},
			want: []string{"a"},
		},
		{
			name: "slice with many elements",
			args: args[string]{
				slice: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
				n:     5,
			},
			want: []string{"a", "b", "c", "d", "e"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FirstNElementsOfSlice(tt.args.slice, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FirstNElementsOfSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertToCommaString(t *testing.T) {
	type testCase[t any] struct {
		name  string
		slice []constants.MarketplaceID
		want  string
	}
	tests := []testCase[constants.MarketplaceID]{
		{
			name:  "empty",
			slice: []constants.MarketplaceID{},
			want:  "",
		},
		{
			name:  "one",
			slice: []constants.MarketplaceID{constants.Germany},
			want:  "A1PA6795UKMFR9",
		},
		{
			name:  "two",
			slice: []constants.MarketplaceID{constants.Germany, constants.Australia},
			want:  "A1PA6795UKMFR9,A39IBJ37TRP1C6",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapToCommaString(tt.slice); got != tt.want {
				t.Errorf("convertToCommaString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddToQueryIfSet(t *testing.T) {
	type args struct {
		q     url.Values
		key   string
		value string
	}
	tests := []struct {
		name string
		args args
		want url.Values
	}{
		{
			name: "simple",
			args: args{
				q:     url.Values{},
				key:   "a",
				value: "b",
			},
			want: map[string][]string{"a": {"b"}},
		},
		{
			name: "overwrite",
			args: args{
				q:     map[string][]string{"a": {"b"}},
				key:   "a",
				value: "c",
			},
			want: map[string][]string{"a": {"c"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddToQueryIfSet(tt.args.q, tt.args.key, tt.args.value)
			if !reflect.DeepEqual(tt.args.q, tt.want) {
				t.Errorf("AddToQueryIfSet() = %v, want %v", tt.args.q, tt.want)
			}
		})
	}
}

func TestNewSet(t *testing.T) {
	type args[T comparable] struct {
		items []T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want map[T]struct{}
	}
	tests := []testCase[string]{
		{
			name: "empty",
			args: args[string]{
				items: []string{},
			},
			want: nil,
		},
		{
			name: "one",
			args: args[string]{
				items: []string{"a"},
			},
			want: map[string]struct{}{"a": {}},
		},
		{
			name: "two",
			args: args[string]{
				items: []string{"a", "b"},
			},
			want: map[string]struct{}{"a": {}, "b": {}},
		},
		{
			name: "two with duplicates",
			args: args[string]{
				items: []string{"a", "b", "a"},
			},
			want: map[string]struct{}{"a": {}, "b": {}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet[string](tt.args.items...)
			if diff := cmp.Diff(set.m, tt.want); diff != "" {
				t.Error(diff)
			}
		})
	}
}
