package common

type ObjectType = string
type Object = string

type ProviderConfig = map[string]string

type Provider interface {
	Name() string
	Types() []ObjectType
	AuthTest() error
	ReadObjects(typeName string) []Object
	DeleteObjects(typeName string, objects []Object)
}
