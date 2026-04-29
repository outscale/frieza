// This provider is just an example you can copy to start implementing your own.
// In this example, we simulate a provider providing "MyResource".
// Check docs/CONTRIBUTING.md for more details
package provider_example

import (
	"context"
	"errors"
	"log"

	. "github.com/outscale/frieza/internal/common"
	"github.com/teris-io/cli"
)

const (
	Name           = "provider_example"
	typeMyResource = "MyResource"
)

type ProviderExample struct {
	apiKey string
}

func checkConfig(config ProviderConfig) error {
	if len(config["api-key"]) == 0 {
		return errors.New("api key is needed")
	}
	return nil
}

func New(config ProviderConfig, debug bool) (*ProviderExample, error) {
	if err := checkConfig(config); err != nil {
		return nil, err
	}
	return &ProviderExample{
		apiKey: config["api-key"],
	}, nil
}

func Types() []ObjectType {
	object_types := []ObjectType{typeMyResource}
	return object_types
}

func Cli() (string, cli.Command) {
	return Name, cli.NewCommand(Name, "create new Example profile").
		WithOption(cli.NewOption("api-key", "Api key"))
}

func (provider *ProviderExample) Name() string {
	return Name
}

func (provider *ProviderExample) Types() []ObjectType {
	return Types()
}

func (provider *ProviderExample) AuthTest(ctx context.Context) error {
	if provider.apiKey != "123" {
		return errors.New("cannot authenticate with API Key")
	}
	return nil
}

func (provider *ProviderExample) ReadObjects(ctx context.Context, typeName string) ([]Object, error) {
	switch typeName {
	case typeMyResource:
		return provider.readMyResources(ctx)
	}
	return []Object{}, nil
}

func (provider *ProviderExample) DeleteObjects(ctx context.Context, typeName string, objects []Object) {
	switch typeName {
	case typeMyResource:
		provider.deleteMyResources(ctx, objects)
	}
}

func (provider *ProviderExample) StringObject(object string, typeName string) string {
	return object
}

func (provider *ProviderExample) readMyResources(ctx context.Context) ([]Object, error) {
	MyResources := make([]Object, 0, 2)
	// Get remote objects
	// ...
	MyResources = append(MyResources, "MyResource-id-1")
	MyResources = append(MyResources, "MyResource-id-2")
	return MyResources, nil
}

func (provider *ProviderExample) deleteMyResources(ctx context.Context, myResources []Object) {
	log.Printf("Deleting MyResources: %s ... ", myResources)
	log.Println("OK")
}
