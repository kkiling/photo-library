package descriptor

import (
	"errors"
	"reflect"
	"runtime"
	"strings"
)

var (
	ErrMethodDescriptorNotFound = errors.New("method descriptor not found")
)

type MethodDescriptor struct {
	Method interface{}
	// TODO: template
	UseAuth int
}

func (m *MethodDescriptor) GetName() (string, error) {
	if m.Method == nil {
		return "", ErrMethodDescriptorNotFound
	}
	methodPointer := reflect.ValueOf(m.Method).Pointer()
	fullName := runtime.FuncForPC(methodPointer).Name()
	methodNameParts := strings.Split(fullName, ".")
	return methodNameParts[len(methodNameParts)-1], nil
}

type MethodDescriptorMap map[string]MethodDescriptor

func (m MethodDescriptorMap) GetByFullName(fullName string) (MethodDescriptor, bool) {
	methodNameParts := strings.Split(fullName, "/")
	methodName := methodNameParts[len(methodNameParts)-1]
	methodDescriptor, ok := m[methodName]
	return methodDescriptor, ok
}

func NewMethodDescriptorMap(
	methodDescriptors []MethodDescriptor,
) (MethodDescriptorMap, error) {
	m := make(MethodDescriptorMap)
	for _, methodDescriptor := range methodDescriptors {
		methodName, err := methodDescriptor.GetName()
		if err != nil {
			return nil, err
		}
		m[methodName] = methodDescriptor
	}
	return m, nil
}
