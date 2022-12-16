package shorten

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShorten(t *testing.T) {
	type args struct {
		id uint32
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "zero arg",
			args: args{
				id: 0,
			},
			want: "",
		},
		{
			name: "id in uint32 124",
			args: args{
				id: 124,
			},
			want: "ca",
		},
		{
			name: "id in uint32 50",
			args: args{
				id: 50,
			},
			want: "Y",
		},
		{
			name: "",
			args: args{
				id: 2048,
			},
			want: "Hc",
		},
		{
			name: "Max is eQPpmd",
			args: args{
				id: 4294967295,
			},
			want: "eQPpmd",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, Shorten(tt.args.id))
		})
	}

	t.Run("same result for same input", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			assert.Equal(t, "Hc", Shorten(2048))
		}
	})
}

func Test_reverse(t *testing.T) {
	type args struct {
		nums []uint32
	}
	tests := []struct {
		name string
		args args
		want []uint32
	}{
		{
			name: "1",
			args: args{nums: []uint32{1, 2, 3}},
			want: []uint32{3, 2, 1},
		},
		{
			name: "2",
			args: args{nums: []uint32{}},
			want: []uint32{},
		},
		{
			name: "2",
			args: args{nums: []uint32{1, 2}},
			want: []uint32{2, 1},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reverse(tt.args.nums)
			if !reflect.DeepEqual(tt.args.nums, tt.want) {
				t.Errorf("Reverse() = %v, want %v", tt.args.nums, tt.want)
			}
		})
	}
}
