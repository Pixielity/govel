package reflector

import (
	"fmt"
	"reflect"
	"sync"
)

// Reflector provides Laravel-style reflection utilities for Go
type Reflector struct {
	// Cache for reflection results to improve performance
	cache sync.Map
}

// Global reflector instance
var defaultReflector = &Reflector{}

// ReflectionResult caches reflection information
type ReflectionResult struct {
	Type        reflect.Type
	Value       reflect.Value
	Methods     map[string]reflect.Method
	Fields      map[string]reflect.StructField
	IsValid     bool
	IsCallable  bool
	IsInterface bool
	IsStruct    bool
	IsPointer   bool
}

// CallableInfo represents information about a callable
type CallableInfo struct {
	IsCallable  bool
	IsMethod    bool
	IsFunction  bool
	HasReceiver bool
	Object      interface{}
	Method      string
	Function    interface{}
	Type        reflect.Type
	NumIn       int
	NumOut      int
	IsVariadic  bool
}

// ParameterInfo represents information about a parameter
type ParameterInfo struct {
	Name      string
	Type      reflect.Type
	Kind      reflect.Kind
	IsBuiltin bool
	IsPointer bool
	IsStruct  bool
	Package   string
	TypeName  string
}

// MethodInfo represents detailed method information
type MethodInfo struct {
	Name       string
	Type       reflect.Type
	Method     reflect.Method
	IsExported bool
	NumIn      int
	NumOut     int
	IsVariadic bool
	Parameters []ParameterInfo
	Returns    []ParameterInfo
}

// IsCallable checks if a variable is callable (similar to PHP's is_callable)
func IsCallable(variable interface{}, syntaxOnly ...bool) bool {
	return defaultReflector.IsCallable(variable, syntaxOnly...)
}

// IsCallable checks if a variable is callable
func (r *Reflector) IsCallable(variable interface{}, syntaxOnly ...bool) bool {
	if variable == nil {
		return false
	}

	syntaxCheck := len(syntaxOnly) > 0 && syntaxOnly[0]

	// Check if it's a function
	v := reflect.ValueOf(variable)
	if v.Kind() == reflect.Func {
		return true
	}

	// Check if it's a slice/array representing a callable [object, method]
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		if v.Len() != 2 {
			return false
		}

		obj := v.Index(0).Interface()
		methodName := v.Index(1).Interface()

		// Method name must be a string
		methodStr, ok := methodName.(string)
		if !ok {
			return false
		}

		if syntaxCheck {
			// Just check syntax - object exists and method name is string
			return obj != nil && methodStr != ""
		}

		// Check if method actually exists and is callable
		return r.HasMethod(obj, methodStr) && r.IsMethodPublic(obj, methodStr)
	}

	return false
}

// GetCallableInfo returns detailed information about a callable
func GetCallableInfo(variable interface{}) *CallableInfo {
	return defaultReflector.GetCallableInfo(variable)
}

// GetCallableInfo returns detailed information about a callable
func (r *Reflector) GetCallableInfo(variable interface{}) *CallableInfo {
	info := &CallableInfo{
		IsCallable: false,
	}

	if variable == nil {
		return info
	}

	v := reflect.ValueOf(variable)
	t := reflect.TypeOf(variable)

	// Check if it's a function
	if v.Kind() == reflect.Func {
		info.IsCallable = true
		info.IsFunction = true
		info.Function = variable
		info.Type = t
		info.NumIn = t.NumIn()
		info.NumOut = t.NumOut()
		info.IsVariadic = t.IsVariadic()
		return info
	}

	// Check if it's a callable array/slice [object, method]
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		if v.Len() == 2 {
			obj := v.Index(0).Interface()
			methodName := v.Index(1).Interface()

			if methodStr, ok := methodName.(string); ok && obj != nil {
				if r.HasMethod(obj, methodStr) && r.IsMethodPublic(obj, methodStr) {
					info.IsCallable = true
					info.IsMethod = true
					info.HasReceiver = true
					info.Object = obj
					info.Method = methodStr

					// Get method type information
					objType := reflect.TypeOf(obj)
					if method, found := objType.MethodByName(methodStr); found {
						info.Type = method.Type
						info.NumIn = method.Type.NumIn() - 1 // Subtract receiver
						info.NumOut = method.Type.NumOut()
						info.IsVariadic = method.Type.IsVariadic()
					}
				}
			}
		}
	}

	return info
}

// HasMethod checks if an object has a specific method
func HasMethod(obj interface{}, methodName string) bool {
	return defaultReflector.HasMethod(obj, methodName)
}

// HasMethod checks if an object has a specific method
func (r *Reflector) HasMethod(obj interface{}, methodName string) bool {
	if obj == nil {
		return false
	}

	objType := reflect.TypeOf(obj)
	_, found := objType.MethodByName(methodName)
	return found
}

// IsMethodPublic checks if a method is public (exported)
func IsMethodPublic(obj interface{}, methodName string) bool {
	return defaultReflector.IsMethodPublic(obj, methodName)
}

// IsMethodPublic checks if a method is public (exported)
func (r *Reflector) IsMethodPublic(obj interface{}, methodName string) bool {
	if obj == nil || methodName == "" {
		return false
	}

	// In Go, exported methods start with uppercase letters
	return methodName[0] >= 'A' && methodName[0] <= 'Z'
}

// GetMethodInfo returns detailed information about a method
func GetMethodInfo(obj interface{}, methodName string) (*MethodInfo, error) {
	return defaultReflector.GetMethodInfo(obj, methodName)
}

// GetMethodInfo returns detailed information about a method
func (r *Reflector) GetMethodInfo(obj interface{}, methodName string) (*MethodInfo, error) {
	if obj == nil {
		return nil, fmt.Errorf("object is nil")
	}

	objType := reflect.TypeOf(obj)
	method, found := objType.MethodByName(methodName)
	if !found {
		return nil, fmt.Errorf("method %s not found on type %s", methodName, objType.String())
	}

	info := &MethodInfo{
		Name:       method.Name,
		Type:       method.Type,
		Method:     method,
		IsExported: r.IsMethodPublic(obj, methodName),
		NumIn:      method.Type.NumIn() - 1, // Subtract receiver
		NumOut:     method.Type.NumOut(),
		IsVariadic: method.Type.IsVariadic(),
		Parameters: make([]ParameterInfo, 0),
		Returns:    make([]ParameterInfo, 0),
	}

	// Get parameter information (skip receiver at index 0)
	for i := 1; i < method.Type.NumIn(); i++ {
		paramType := method.Type.In(i)
		param := ParameterInfo{
			Name:      fmt.Sprintf("param%d", i-1),
			Type:      paramType,
			Kind:      paramType.Kind(),
			IsBuiltin: r.IsBuiltinType(paramType),
			IsPointer: paramType.Kind() == reflect.Ptr,
			IsStruct:  paramType.Kind() == reflect.Struct,
			Package:   paramType.PkgPath(),
			TypeName:  paramType.String(),
		}
		info.Parameters = append(info.Parameters, param)
	}

	// Get return type information
	for i := 0; i < method.Type.NumOut(); i++ {
		returnType := method.Type.Out(i)
		ret := ParameterInfo{
			Name:      fmt.Sprintf("return%d", i),
			Type:      returnType,
			Kind:      returnType.Kind(),
			IsBuiltin: r.IsBuiltinType(returnType),
			IsPointer: returnType.Kind() == reflect.Ptr,
			IsStruct:  returnType.Kind() == reflect.Struct,
			Package:   returnType.PkgPath(),
			TypeName:  returnType.String(),
		}
		info.Returns = append(info.Returns, ret)
	}

	return info, nil
}

// GetParameterTypes returns the types of method/function parameters
func GetParameterTypes(variable interface{}) ([]reflect.Type, error) {
	return defaultReflector.GetParameterTypes(variable)
}

// GetParameterTypes returns the types of method/function parameters
func (r *Reflector) GetParameterTypes(variable interface{}) ([]reflect.Type, error) {
	if variable == nil {
		return nil, fmt.Errorf("variable is nil")
	}

	var funcType reflect.Type

	v := reflect.ValueOf(variable)
	if v.Kind() == reflect.Func {
		funcType = reflect.TypeOf(variable)
	} else {
		// Try to handle callable arrays
		callableInfo := r.GetCallableInfo(variable)
		if !callableInfo.IsCallable {
			return nil, fmt.Errorf("variable is not callable")
		}
		funcType = callableInfo.Type
	}

	var types []reflect.Type
	startIndex := 0

	// Skip receiver for methods
	if funcType.NumIn() > 0 {
		firstParam := funcType.In(0)
		// If the first parameter looks like a receiver (struct or pointer to struct)
		if firstParam.Kind() == reflect.Struct ||
			(firstParam.Kind() == reflect.Ptr && firstParam.Elem().Kind() == reflect.Struct) {
			startIndex = 1
		}
	}

	for i := startIndex; i < funcType.NumIn(); i++ {
		types = append(types, funcType.In(i))
	}

	return types, nil
}

// GetReturnTypes returns the types of method/function return values
func GetReturnTypes(variable interface{}) ([]reflect.Type, error) {
	return defaultReflector.GetReturnTypes(variable)
}

// GetReturnTypes returns the types of method/function return values
func (r *Reflector) GetReturnTypes(variable interface{}) ([]reflect.Type, error) {
	if variable == nil {
		return nil, fmt.Errorf("variable is nil")
	}

	var funcType reflect.Type

	v := reflect.ValueOf(variable)
	if v.Kind() == reflect.Func {
		funcType = reflect.TypeOf(variable)
	} else {
		// Try to handle callable arrays
		callableInfo := r.GetCallableInfo(variable)
		if !callableInfo.IsCallable {
			return nil, fmt.Errorf("variable is not callable")
		}
		funcType = callableInfo.Type
	}

	var types []reflect.Type
	for i := 0; i < funcType.NumOut(); i++ {
		types = append(types, funcType.Out(i))
	}

	return types, nil
}

// IsBuiltinType checks if a type is a built-in Go type
func IsBuiltinType(t reflect.Type) bool {
	return defaultReflector.IsBuiltinType(t)
}

// IsBuiltinType checks if a type is a built-in Go type
func (r *Reflector) IsBuiltinType(t reflect.Type) bool {
	if t == nil {
		return false
	}

	// Built-in types have empty package path
	return t.PkgPath() == ""
}

// GetTypeName returns the fully qualified type name
func GetTypeName(obj interface{}) string {
	return defaultReflector.GetTypeName(obj)
}

// GetTypeName returns the fully qualified type name
func (r *Reflector) GetTypeName(obj interface{}) string {
	if obj == nil {
		return "nil"
	}

	t := reflect.TypeOf(obj)
	if t.PkgPath() == "" {
		return t.Name()
	}
	return t.PkgPath() + "." + t.Name()
}

// IsSubclassOf checks if a type is a subclass/implements another type
func IsSubclassOf(obj interface{}, parentType reflect.Type) bool {
	return defaultReflector.IsSubclassOf(obj, parentType)
}

// IsSubclassOf checks if a type is a subclass/implements another type
func (r *Reflector) IsSubclassOf(obj interface{}, parentType reflect.Type) bool {
	if obj == nil || parentType == nil {
		return false
	}

	objType := reflect.TypeOf(obj)

	// Check if it implements the interface
	if parentType.Kind() == reflect.Interface {
		return objType.Implements(parentType)
	}

	// Check if it's assignable to the parent type
	return objType.AssignableTo(parentType)
}

// GetReflectionResult gets or creates cached reflection information
func GetReflectionResult(obj interface{}) *ReflectionResult {
	return defaultReflector.GetReflectionResult(obj)
}

// GetReflectionResult gets or creates cached reflection information
func (r *Reflector) GetReflectionResult(obj interface{}) *ReflectionResult {
	if obj == nil {
		return &ReflectionResult{IsValid: false}
	}

	// Use type string as cache key
	t := reflect.TypeOf(obj)
	cacheKey := t.String()

	// Check cache first
	if cached, found := r.cache.Load(cacheKey); found {
		if result, ok := cached.(*ReflectionResult); ok {
			// Update value for this instance
			result.Value = reflect.ValueOf(obj)
			return result
		}
	}

	// Create new reflection result
	v := reflect.ValueOf(obj)
	result := &ReflectionResult{
		Type:        t,
		Value:       v,
		Methods:     make(map[string]reflect.Method),
		Fields:      make(map[string]reflect.StructField),
		IsValid:     v.IsValid(),
		IsCallable:  v.Kind() == reflect.Func,
		IsInterface: t.Kind() == reflect.Interface,
		IsStruct:    t.Kind() == reflect.Struct || (t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct),
		IsPointer:   t.Kind() == reflect.Ptr,
	}

	// Cache methods
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		result.Methods[method.Name] = method
	}

	// Cache fields if it's a struct
	structType := t
	if t.Kind() == reflect.Ptr {
		structType = t.Elem()
	}

	if structType.Kind() == reflect.Struct {
		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)
			result.Fields[field.Name] = field
		}
	}

	// Store in cache
	r.cache.Store(cacheKey, result)

	return result
}

// CallMethod calls a method on an object with the given arguments
func CallMethod(obj interface{}, methodName string, args ...interface{}) ([]reflect.Value, error) {
	return defaultReflector.CallMethod(obj, methodName, args...)
}

// CallMethod calls a method on an object with the given arguments
func (r *Reflector) CallMethod(obj interface{}, methodName string, args ...interface{}) ([]reflect.Value, error) {
	if obj == nil {
		return nil, fmt.Errorf("object is nil")
	}

	v := reflect.ValueOf(obj)
	method := v.MethodByName(methodName)

	if !method.IsValid() {
		return nil, fmt.Errorf("method %s not found on type %s", methodName, reflect.TypeOf(obj).String())
	}

	// Convert arguments to reflect.Values
	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		reflectArgs[i] = reflect.ValueOf(arg)
	}

	// Call the method
	return method.Call(reflectArgs), nil
}

// GetMethodByName returns a method by name
func GetMethodByName(obj interface{}, methodName string) (reflect.Method, bool) {
	return defaultReflector.GetMethodByName(obj, methodName)
}

// GetMethodByName returns a method by name
func (r *Reflector) GetMethodByName(obj interface{}, methodName string) (reflect.Method, bool) {
	if obj == nil {
		return reflect.Method{}, false
	}

	t := reflect.TypeOf(obj)
	return t.MethodByName(methodName)
}

// GetAllMethods returns all methods of an object
func GetAllMethods(obj interface{}) []reflect.Method {
	return defaultReflector.GetAllMethods(obj)
}

// GetAllMethods returns all methods of an object
func (r *Reflector) GetAllMethods(obj interface{}) []reflect.Method {
	if obj == nil {
		return nil
	}

	t := reflect.TypeOf(obj)
	methods := make([]reflect.Method, t.NumMethod())

	for i := 0; i < t.NumMethod(); i++ {
		methods[i] = t.Method(i)
	}

	return methods
}

// GetPublicMethods returns only public (exported) methods
func GetPublicMethods(obj interface{}) []reflect.Method {
	return defaultReflector.GetPublicMethods(obj)
}

// GetPublicMethods returns only public (exported) methods
func (r *Reflector) GetPublicMethods(obj interface{}) []reflect.Method {
	allMethods := r.GetAllMethods(obj)
	var publicMethods []reflect.Method

	for _, method := range allMethods {
		if r.IsMethodPublic(obj, method.Name) {
			publicMethods = append(publicMethods, method)
		}
	}

	return publicMethods
}

// GetFieldByName returns a struct field by name
func GetFieldByName(obj interface{}, fieldName string) (reflect.StructField, bool) {
	return defaultReflector.GetFieldByName(obj, fieldName)
}

// GetFieldByName returns a struct field by name
func (r *Reflector) GetFieldByName(obj interface{}, fieldName string) (reflect.StructField, bool) {
	if obj == nil {
		return reflect.StructField{}, false
	}

	t := reflect.TypeOf(obj)

	// Handle pointer to struct
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return reflect.StructField{}, false
	}

	return t.FieldByName(fieldName)
}

// IsExported checks if a name is exported (public)
func IsExported(name string) bool {
	return defaultReflector.IsExported(name)
}

// IsExported checks if a name is exported (public)
func (r *Reflector) IsExported(name string) bool {
	if name == "" {
		return false
	}
	return name[0] >= 'A' && name[0] <= 'Z'
}

// GetPackagePath returns the package path of a type
func GetPackagePath(obj interface{}) string {
	return defaultReflector.GetPackagePath(obj)
}

// GetPackagePath returns the package path of a type
func (r *Reflector) GetPackagePath(obj interface{}) string {
	if obj == nil {
		return ""
	}

	t := reflect.TypeOf(obj)
	return t.PkgPath()
}

// CreateInstance creates a new instance of the same type as the given object
func CreateInstance(obj interface{}) (interface{}, error) {
	return defaultReflector.CreateInstance(obj)
}

// CreateInstance creates a new instance of the same type as the given object
func (r *Reflector) CreateInstance(obj interface{}) (interface{}, error) {
	if obj == nil {
		return nil, fmt.Errorf("object is nil")
	}

	t := reflect.TypeOf(obj)

	// Handle pointer types
	if t.Kind() == reflect.Ptr {
		elem := t.Elem()
		instance := reflect.New(elem)
		return instance.Interface(), nil
	}

	// Handle value types
	instance := reflect.New(t)
	return instance.Elem().Interface(), nil
}

// DeepEqual performs deep comparison of two values
func DeepEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

// IsZero checks if a value is the zero value for its type
func IsZero(obj interface{}) bool {
	return defaultReflector.IsZero(obj)
}

// IsZero checks if a value is the zero value for its type
func (r *Reflector) IsZero(obj interface{}) bool {
	if obj == nil {
		return true
	}

	v := reflect.ValueOf(obj)
	return v.IsZero()
}

// GetInterfaceType returns the underlying type if obj is an interface
func GetInterfaceType(obj interface{}) reflect.Type {
	return defaultReflector.GetInterfaceType(obj)
}

// GetInterfaceType returns the underlying type if obj is an interface
func (r *Reflector) GetInterfaceType(obj interface{}) reflect.Type {
	if obj == nil {
		return nil
	}

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Interface {
		return v.Elem().Type()
	}
	return v.Type()
}

// ToString provides a debug-friendly string representation
func ToString(obj interface{}) string {
	return defaultReflector.ToString(obj)
}

// ToString provides a debug-friendly string representation
func (r *Reflector) ToString(obj interface{}) string {
	if obj == nil {
		return "nil"
	}

	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	return fmt.Sprintf("Type: %s, Kind: %s, Value: %v", t.String(), t.Kind().String(), v.Interface())
}
