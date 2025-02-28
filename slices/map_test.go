package slices

import (
	"reflect"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	type args[T, U any] struct {
		items []T
		fn    func(T) U
	}
	type Test[T, U any] struct {
		name string
		args args[T, U]
		want []U
	}

	tests := []Test[int, string]{{
		name: "main",
		args: args[int, string]{
			items: []int{1, 2, 3},
			fn: func(i int) string {
				return strconv.Itoa(i)
			},
		},
		want: []string{"1", "2", "3"},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.args.items, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
