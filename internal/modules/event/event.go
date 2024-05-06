package event

// TODO: move this to GEM or something
//import (
//	"fmt"
//	"reflect"
//)
//
//// implementsInterface checks if the provided struct (structPtr) implements the
//// specified interface (interfaceType). It returns true if the struct
//// implements the interface, otherwise false.
//func ImplementsInterface(structPtr interface{}, iface interface{}) bool {
//	structType := reflect.TypeOf(structPtr)
//	if structType.Kind() != reflect.Ptr {
//		fmt.Println("Error: First argument must be a pointer to a struct")
//		return false
//	}
//
//	// Get the interface type from the empty interface value passed.
//	interfaceType := reflect.TypeOf(iface).Elem()
//
//	// Check if the type of the struct implements the interface.
//	return structType.Implements(interfaceType)
//}
//
//// ImplementsInterfaceGeneric checks if the provided struct pointer implements the interface T.
//// T should be an interface type, and structPtr should be a pointer to a struct.
//func ImplementsInterfaceGeneric[T any](structPtr any) bool {
//	structType := reflect.TypeOf(structPtr)
//	if structType.Kind() != reflect.Ptr {
//		fmt.Println("Error: structPtr must be a pointer to a struct")
//		return false
//	}
//	// Get the type of the interface from the generic type T.
//	interfaceType := reflect.TypeOf((*T)(nil)).Elem()
//	// Check if the struct type implements the interface type.
//	return structType.Implements(interfaceType)
//}
