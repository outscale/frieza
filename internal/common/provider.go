package common

import "context"

type (
	ObjectType = string
	Object     = string
)

type ProviderConfig = map[string]string

type Provider interface {
	Name() string
	Types() []ObjectType
	AuthTest(ctx context.Context) error
	ReadObjects(ctx context.Context, typeName string) ([]Object, error)
	DeleteObjects(ctx context.Context, typeName string, objects []Object)
	StringObject(object string, typeName string) string
}
