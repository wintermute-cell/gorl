package store

import (
	"reflect"
)

// store holds values keyed by their type.
type store struct {
	data map[reflect.Type]interface{}
}

// globalStore is the default store instance.
var globalStore = newStore()

// newStore creates a new instance of Store.
func newStore() *store {
	return &store{
		data: make(map[reflect.Type]interface{}),
	}
}

// Add adds or replaces a value in the store keyed by its type.
func Add[T any](value T) {
	typ := reflect.TypeOf(value)
	globalStore.data[typ] = value
}

// Get retrieves a value from the store by its type. It returns the value and a
// boolean indicating if it was found.
func Get[T any]() (T, bool) {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	v, ok := globalStore.data[typ]
	if !ok {
		return *new(T), false
	}
	return v.(T), true
}
