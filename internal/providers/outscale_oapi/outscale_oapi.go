package outscale_oapi

import (
	"context"
	"errors"
	"fmt"
	"os"

	. "github.com/outscale-dev/frieza/internal/common"
	osc "github.com/outscale/osc-sdk-go/v2"
	"github.com/teris-io/cli"
)

const Name = "outscale_oapi"
const typeVm = "vm"

type OutscaleOAPI struct {
	client  *osc.APIClient
	context context.Context
}

func checkConfig(config ProviderConfig) error {
	if len(config["ak"]) == 0 {
		return errors.New("access key is needed")
	}
	if len(config["sk"]) == 0 {
		return errors.New("secret key is needed")
	}
	if len(config["region"]) == 0 {
		return errors.New("region is needed")
	}
	return nil
}

func New(config ProviderConfig) (*OutscaleOAPI, error) {
	if err := checkConfig(config); err != nil {
		return nil, err
	}
	oscConfig := osc.NewConfiguration()
	oscConfig.Debug = false
	client := osc.NewAPIClient(oscConfig)
	ctx := context.WithValue(context.Background(), osc.ContextAWSv4, osc.AWSv4{
		AccessKey: config["ak"],
		SecretKey: config["sk"],
	})
	ctx = context.WithValue(ctx, osc.ContextServerIndex, 0)
	ctx = context.WithValue(ctx, osc.ContextServerVariables, map[string]string{"region": config["region"]})
	return &OutscaleOAPI{
		client:  client,
		context: ctx,
	}, nil
}

func Types() []ObjectType {
	object_types := []ObjectType{typeVm}
	return object_types
}

func Cli() (string, cli.Command) {
	return Name, cli.NewCommand(Name, "create new Outscale API profile").
		WithOption(cli.NewOption("region", "Outscale region (e.g. eu-west-2)")).
		WithOption(cli.NewOption("ak", "access key")).
		WithOption(cli.NewOption("sk", "secret key"))
}

func (provider *OutscaleOAPI) Name() string {
	return Name
}

func (provider *OutscaleOAPI) Types() []ObjectType {
	return Types()
}

func (provider *OutscaleOAPI) AuthTest() error {
	_, httpRes, err := provider.client.AccountApi.ReadAccounts(provider.context).
		ReadAccountsRequest(osc.ReadAccountsRequest{}).
		Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:")
		if httpRes != nil {
			fmt.Fprintln(os.Stderr, httpRes.Status)
		}
	}
	return nil
}

func newObjects() Objects {
	objects := make(Objects)
	for _, typeName := range Types() {
		objects[typeName] = make([]Object, 0)
	}
	return objects
}

func (provider *OutscaleOAPI) Objects() Objects {
	objects := newObjects()
	objects[typeVm] = provider.getVms()
	return objects
}

func (provider *OutscaleOAPI) Delete(objects Objects) {
	vms := objects[typeVm]
	if vms != nil {
		provider.deleteVms(vms)
	}
}

func (provider *OutscaleOAPI) getVms() []Object {
	vms := make([]Object, 0)
	read, httpRes, err := provider.client.VmApi.ReadVms(provider.context).
		ReadVmsRequest(osc.ReadVmsRequest{}).
		Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while reading vms ")
		if httpRes != nil {
			fmt.Fprintln(os.Stderr, httpRes.Status)
		}
		return vms
	}
	for _, vm := range *read.Vms {
		switch *vm.State {
		case "pending", "running", "stopping", "stopped", "shutting-down", "quarantine":
			vms = append(vms, *vm.VmId)
		}
	}
	return vms
}

func (provider *OutscaleOAPI) deleteVms(vms []Object) {
	if len(vms) == 0 {
		return
	}
	fmt.Printf("Deleting virtual machines: %s ... ", vms)
	deletionOpts := osc.DeleteVmsRequest{VmIds: vms}
	_, httpRes, err := provider.client.VmApi.DeleteVms(provider.context).
		DeleteVmsRequest(deletionOpts).
		Execute()
	if err != nil {
		fmt.Fprint(os.Stderr, "Error while deleting vms:")
		if httpRes != nil {
			fmt.Fprintln(os.Stderr, httpRes.Status)
		}
	} else {
		fmt.Println("OK")
	}
}
