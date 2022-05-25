package main

import (
	"reflect"
	"testing"
)

func TestNewCache(t *testing.T) {
	tests := []struct {
		name string
		want Cache
	}{
		{
			name: "create new cache",
			want: Cache{memoryCache: make(map[string]*Data)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCache(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Put(t *testing.T) {
	type args struct {
		key   string
		value Data
	}
	tests := []struct {
		name string
		c    *Cache
		args args
	}{
		{
			name: "put",
			c:    &Cache{memoryCache: make(map[string]*Data)},
			args: args{key: "abc", value: Data{Payment: Payment{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Put(tt.args.key, tt.args.value)
		})
	}
}

func TestCache_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		c       *Cache
		args    args
		want    Data
		wantErr bool
	}{
		{
			name:    "get",
			c:       &Cache{memoryCache: make(map[string]*Data)},
			args:    args{key: "abc"},
			want:    Data{Payment: Payment{Transaction: "qwer"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Put("abc", Data{Payment: Payment{Transaction: "qwer"}})
			got, err := tt.c.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cache.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cache.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
