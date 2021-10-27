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
const typeSecurityGroup = "security_group"
const typePublicIp = "public_ip"
const typeVolume = "volume"
const typeKeypair = "keypair"
const typeRouteTable = "route_table"
const typeInternetService = "internet_service"
const typeSubnet = "subnet"
const typeNet = "net"
const typeImage = "image"

type OutscaleOAPI struct {
	client    *osc.APIClient
	context   context.Context
	accountId *string
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
	object_types := []ObjectType{
		typeVm,
		typeSecurityGroup,
		typePublicIp,
		typeVolume,
		typeKeypair,
		typeRouteTable,
		typeInternetService,
		typeSubnet,
		typeNet,
		typeImage,
	}
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
	objects[typeSecurityGroup] = provider.getSecurityGroups()
	objects[typePublicIp] = provider.getPublicIps()
	objects[typeVolume] = provider.getVolumes()
	objects[typeKeypair] = provider.getKeypairs()
	objects[typeRouteTable] = provider.getRouteTables()
	objects[typeInternetService] = provider.getInternetServices()
	objects[typeSubnet] = provider.getSubnets()
	objects[typeNet] = provider.getNets()
	objects[typeImage] = provider.getImages()
	return objects
}

func (provider *OutscaleOAPI) Delete(objects Objects) {
	provider.deleteVms(objects[typeVm])
	provider.deleteSecurityGroups(objects[typeSecurityGroup])
	provider.deletePublicIps(objects[typePublicIp])
	provider.deleteVolumes(objects[typeVolume])
	provider.deleteKeypairs(objects[typeKeypair])
	provider.deleteRouteTables(objects[typeRouteTable])
	provider.deleteInternetServices(objects[typeInternetService])
	provider.deleteSubnets(objects[typeSubnet])
	provider.deleteNets(objects[typeNet])
	provider.deleteImages(objects[typeImage])
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

func (provider *OutscaleOAPI) getSecurityGroups() []Object {
	securityGroups := make([]Object, 0)
	read, httpRes, err := provider.client.SecurityGroupApi.
		ReadSecurityGroups(provider.context).
		ReadSecurityGroupsRequest(osc.ReadSecurityGroupsRequest{}).
		Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while reading security groups")
		if httpRes != nil {
			fmt.Fprintln(os.Stderr, httpRes.Status)
		}
		return securityGroups
	}
	for _, sg := range *read.SecurityGroups {
		if *sg.SecurityGroupName == "default" {
			continue
		}
		securityGroups = append(securityGroups, *sg.SecurityGroupId)
	}
	return securityGroups
}

func (provider *OutscaleOAPI) deleteSecurityGroups(securityGroups []Object) {
	if len(securityGroups) == 0 {
		return
	}
	for _, sg := range securityGroups {
		fmt.Printf("Deleting security group %s... ", sg)
		deletionOpts := osc.DeleteSecurityGroupRequest{SecurityGroupId: &sg}
		_, httpRes, err := provider.client.SecurityGroupApi.
			DeleteSecurityGroup(provider.context).
			DeleteSecurityGroupRequest(deletionOpts).
			Execute()
		if err != nil {
			fmt.Fprint(os.Stderr, "Error while deleting security groups")
			if httpRes != nil {
				fmt.Fprintln(os.Stderr, httpRes.Status)
			}
		} else {
			fmt.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) getPublicIps() []Object {
	publicIps := make([]Object, 0)
	read, httpRes, err := provider.client.PublicIpApi.
		ReadPublicIps(provider.context).
		ReadPublicIpsRequest(osc.ReadPublicIpsRequest{}).
		Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while reading public ips")
		if httpRes != nil {
			fmt.Fprintln(os.Stderr, httpRes.Status)
		}
		return publicIps
	}
	for _, pip := range *read.PublicIps {
		publicIps = append(publicIps, *pip.PublicIpId)
	}
	return publicIps
}

func (provider *OutscaleOAPI) deletePublicIps(publicIps []Object) {
	if len(publicIps) == 0 {
		return
	}
	for _, pip := range publicIps {
		fmt.Printf("Deleting public ip %s... ", pip)
		deletionOpts := osc.DeletePublicIpRequest{PublicIpId: &pip}
		_, httpRes, err := provider.client.PublicIpApi.
			DeletePublicIp(provider.context).
			DeletePublicIpRequest(deletionOpts).
			Execute()
		if err != nil {
			fmt.Fprint(os.Stderr, "Error while deleting public ips")
			if httpRes != nil {
				fmt.Fprintln(os.Stderr, httpRes.Status)
			}
		} else {
			fmt.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) getVolumes() []Object {
	volumes := make([]Object, 0)
	read, httpRes, err := provider.client.VolumeApi.
		ReadVolumes(provider.context).
		ReadVolumesRequest(osc.ReadVolumesRequest{}).
		Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while reading volumes")
		if httpRes != nil {
			fmt.Fprintln(os.Stderr, httpRes.Status)
		}
		return volumes
	}
	for _, volume := range *read.Volumes {
		volumes = append(volumes, *volume.VolumeId)
	}
	return volumes
}

func (provider *OutscaleOAPI) deleteVolumes(volumes []Object) {
	if len(volumes) == 0 {
		return
	}
	for _, volume := range volumes {
		fmt.Printf("Deleting public ip %s... ", volume)
		deletionOpts := osc.DeleteVolumeRequest{VolumeId: volume}
		_, httpRes, err := provider.client.VolumeApi.
			DeleteVolume(provider.context).
			DeleteVolumeRequest(deletionOpts).
			Execute()
		if err != nil {
			fmt.Fprint(os.Stderr, "Error while deleting volume")
			if httpRes != nil {
				fmt.Fprintln(os.Stderr, httpRes.Status)
			}
		} else {
			fmt.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) getKeypairs() []Object {
	keypairs := make([]Object, 0)
	read, httpRes, err := provider.client.KeypairApi.ReadKeypairs(provider.context).
		ReadKeypairsRequest(osc.ReadKeypairsRequest{}).
		Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while reading keypairs ")
		if httpRes != nil {
			fmt.Fprintln(os.Stderr, httpRes.Status)
		}
		return keypairs
	}
	for _, keypair := range *read.Keypairs {
		keypairs = append(keypairs, *keypair.KeypairName)
	}
	return keypairs
}

func (provider *OutscaleOAPI) deleteKeypairs(keypairs []Object) {
	if len(keypairs) == 0 {
		return
	}
	for _, keypair := range keypairs {
		fmt.Printf("Deleting keypair %s... ", keypair)
		deletionOpts := osc.DeleteKeypairRequest{KeypairName: keypair}
		_, httpRes, err := provider.client.KeypairApi.
			DeleteKeypair(provider.context).
			DeleteKeypairRequest(deletionOpts).
			Execute()
		if err != nil {
			fmt.Fprint(os.Stderr, "Error while deleting keypair")
			if httpRes != nil {
				fmt.Fprintln(os.Stderr, httpRes.Status)
			}
		} else {
			fmt.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) getRouteTables() []Object {
	routeTables := make([]Object, 0)
	read, httpRes, err := provider.client.RouteTableApi.ReadRouteTables(provider.context).
		ReadRouteTablesRequest(osc.ReadRouteTablesRequest{}).
		Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while reading route tables ")
		if httpRes != nil {
			fmt.Fprintln(os.Stderr, httpRes.Status)
		}
		return routeTables
	}
	for _, routeTable := range *read.RouteTables {
		routeTables = append(routeTables, *routeTable.RouteTableId)
	}
	return routeTables
}

func (provider *OutscaleOAPI) deleteRouteTables(routeTables []Object) {
	if len(routeTables) == 0 {
		return
	}
	for _, routeTable := range routeTables {
		fmt.Printf("Deleting route table %s... ", routeTable)
		deletionOpts := osc.DeleteRouteTableRequest{RouteTableId: routeTable}
		_, httpRes, err := provider.client.RouteTableApi.
			DeleteRouteTable(provider.context).
			DeleteRouteTableRequest(deletionOpts).
			Execute()
		if err != nil {
			fmt.Fprint(os.Stderr, "Error while deleting route table")
			if httpRes != nil {
				fmt.Fprintln(os.Stderr, httpRes.Status)
			}
		} else {
			fmt.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) getInternetServices() []Object {
	internetServices := make([]Object, 0)
	read, httpRes, err := provider.client.InternetServiceApi.ReadInternetServices(provider.context).
		ReadInternetServicesRequest(osc.ReadInternetServicesRequest{}).
		Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while reading internet services ")
		if httpRes != nil {
			fmt.Fprintln(os.Stderr, httpRes.Status)
		}
		return internetServices
	}
	for _, internetService := range *read.InternetServices {
		internetServices = append(internetServices, *internetService.InternetServiceId)
	}
	return internetServices
}

func (provider *OutscaleOAPI) deleteInternetServices(internetServices []Object) {
	if len(internetServices) == 0 {
		return
	}
	for _, internetService := range internetServices {
		fmt.Printf("Deleting internet service %s... ", internetService)
		deletionOpts := osc.DeleteInternetServiceRequest{InternetServiceId: internetService}
		_, httpRes, err := provider.client.InternetServiceApi.
			DeleteInternetService(provider.context).
			DeleteInternetServiceRequest(deletionOpts).
			Execute()
		if err != nil {
			fmt.Fprint(os.Stderr, "Error while deleting internet service")
			if httpRes != nil {
				fmt.Fprintln(os.Stderr, httpRes.Status)
			}
		} else {
			fmt.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) getSubnets() []Object {
	subnets := make([]Object, 0)
	read, httpRes, err := provider.client.SubnetApi.ReadSubnets(provider.context).
		ReadSubnetsRequest(osc.ReadSubnetsRequest{}).
		Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while reading subnets ")
		if httpRes != nil {
			fmt.Fprintln(os.Stderr, httpRes.Status)
		}
		return subnets
	}
	for _, subnet := range *read.Subnets {
		subnets = append(subnets, *subnet.SubnetId)
	}
	return subnets
}

func (provider *OutscaleOAPI) deleteSubnets(subnets []Object) {
	if len(subnets) == 0 {
		return
	}
	for _, subnet := range subnets {
		fmt.Printf("Deleting subnet %s... ", subnet)
		deletionOpts := osc.DeleteSubnetRequest{SubnetId: subnet}
		_, httpRes, err := provider.client.SubnetApi.
			DeleteSubnet(provider.context).
			DeleteSubnetRequest(deletionOpts).
			Execute()
		if err != nil {
			fmt.Fprint(os.Stderr, "Error while deleting subnet")
			if httpRes != nil {
				fmt.Fprintln(os.Stderr, httpRes.Status)
			}
		} else {
			fmt.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) getNets() []Object {
	nets := make([]Object, 0)
	read, httpRes, err := provider.client.NetApi.ReadNets(provider.context).
		ReadNetsRequest(osc.ReadNetsRequest{}).
		Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while reading nets ")
		if httpRes != nil {
			fmt.Fprintln(os.Stderr, httpRes.Status)
		}
		return nets
	}
	for _, net := range *read.Nets {
		nets = append(nets, *net.NetId)
	}
	return nets
}

func (provider *OutscaleOAPI) deleteNets(nets []Object) {
	if len(nets) == 0 {
		return
	}
	for _, net := range nets {
		fmt.Printf("Deleting net %s... ", net)
		deletionOpts := osc.DeleteNetRequest{NetId: net}
		_, httpRes, err := provider.client.NetApi.
			DeleteNet(provider.context).
			DeleteNetRequest(deletionOpts).
			Execute()
		if err != nil {
			fmt.Fprint(os.Stderr, "Error while deleting net")
			if httpRes != nil {
				fmt.Fprintln(os.Stderr, httpRes.Status)
			}
		} else {
			fmt.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) getAccountId() (error, *string) {
	if provider.accountId == nil {
		read, httpRes, err := provider.client.AccountApi.ReadAccounts(provider.context).
			ReadAccountsRequest(osc.ReadAccountsRequest{}).
			Execute()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error while reading account")
			if httpRes != nil {
				fmt.Fprintln(os.Stderr, httpRes.Status)
			}
			return err, nil
		}
		if len(*read.Accounts) == 0 {
			fmt.Fprintln(os.Stderr, "Error while reading account: no account listed")
			return err, nil
		}
		provider.accountId = (*read.Accounts)[0].AccountId
	}
	return nil, provider.accountId
}

func (provider *OutscaleOAPI) getImages() []Object {
	images := make([]Object, 0)
	err, accountId := provider.getAccountId()
	if err != nil {
		return images
	}
	var accountIds []string
	accountIds = append(accountIds, *accountId)
	read, httpRes, err := provider.client.ImageApi.ReadImages(provider.context).
		ReadImagesRequest(osc.ReadImagesRequest{
			Filters: &osc.FiltersImage{
				AccountIds: &accountIds,
			}}).
		Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while reading images ")
		if httpRes != nil {
			fmt.Fprintln(os.Stderr, httpRes.Status)
		}
		return images
	}
	for _, image := range *read.Images {
		images = append(images, *image.ImageId)
	}
	return images
}

func (provider *OutscaleOAPI) deleteImages(images []Object) {
	if len(images) == 0 {
		return
	}
	for _, image := range images {
		fmt.Printf("Deleting image %s... ", image)
		deletionOpts := osc.DeleteImageRequest{ImageId: image}
		_, httpRes, err := provider.client.ImageApi.
			DeleteImage(provider.context).
			DeleteImageRequest(deletionOpts).
			Execute()
		if err != nil {
			fmt.Fprint(os.Stderr, "Error while deleting image")
			if httpRes != nil {
				fmt.Fprintln(os.Stderr, httpRes.Status)
			}
		} else {
			fmt.Println("OK")
		}
	}
}
