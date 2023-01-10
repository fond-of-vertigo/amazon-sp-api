package apis

import (
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
