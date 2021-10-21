package common

type ObjectType = string
type Object = string
type Objects = map[ObjectType][]Object

type ProviderConfig = map[string]string

type Provider interface {
	Name() string
	Types() []ObjectType
	AuthTest() error
	Objects() Objects
	Delete(objects Objects)
}
