// This provider is just an example you can copy to start implementing your own.
// In this example, we simulate a provider providing "MyResource".
// Check docs/CONTRIBUTING.md for more details
package provider_example

import (
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

func (provider *ProviderExample) AuthTest() error {
	if provider.apiKey != "123" {
		return errors.New("cannot authenticate with API Key")
	}
	return nil
}

func (provider *ProviderExample) ReadObjects(typeName string) []Object {
	switch typeName {
	case typeMyResource:
		return provider.readMyResources()
	}
	return []Object{}
}

func (provider *ProviderExample) DeleteObjects(typeName string, objects []Object) {
	switch typeName {
	case typeMyResource:
		provider.deleteMyResources(objects)
	}
}

func (provider *ProviderExample) StringObject(object string, typeName string) string {
	return object
}

func (provider *ProviderExample) readMyResources() []Object {
	MyResources := make([]Object, 0)
	// Get remote objects
	// ...
	MyResources = append(MyResources, "MyResource-id-1")
	MyResources = append(MyResources, "MyResource-id-2")
	return MyResources
}

func (provider *ProviderExample) deleteMyResources(MyResources []Object) {
	log.Printf("Deleting MyResources: %s ... ", MyResources)
	log.Println("OK")
}
