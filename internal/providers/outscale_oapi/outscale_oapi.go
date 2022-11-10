package outscale_oapi

import (
	"context"
	"errors"
	"log"

	. "github.com/outscale-dev/frieza/internal/common"
	osc "github.com/outscale/osc-sdk-go/v2"
	"github.com/teris-io/cli"
)

const Name = "outscale_oapi"
const typeVm = "vm"
const typeLoadBalancer = "load_balancer"
const typeNatService = "nat_service"
const typeSecurityGroup = "security_group"
const typePublicIp = "public_ip"
const typeVolume = "volume"
const typeKeypair = "keypair"
const typeRouteTable = "route_table"
const typeInternetService = "internet_service"
const typeSubnet = "subnet"
const typeNet = "net"
const typeImage = "image"
const typeSnapshot = "snapshot"
const typeVpnConnection = "vpn_connection"
const typeVirtualGateway = "virtual_gateway"
const typeNic = "nic"
const typeAccessKey = "access_key"

type OutscaleOAPI struct {
	client  *osc.APIClient
	context context.Context
	cache   apiCache
}

type apiCache struct {
	accountId        *string
	internetServices map[Object]*osc.InternetService
	publicIps        map[Object]*osc.PublicIp
	vms              map[Object]*osc.Vm
	nics             map[Object]*osc.Nic
	routeTables      map[Object]*osc.RouteTable
	securityGroups   map[Object]*osc.SecurityGroup
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

func New(config ProviderConfig, debug bool) (*OutscaleOAPI, error) {
	if err := checkConfig(config); err != nil {
		return nil, err
	}
	oscConfig := osc.NewConfiguration()
	oscConfig.Debug = debug
	oscConfig.UserAgent = "frieza/" + FullVersion()
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
		cache:   newAPICache(),
	}, nil
}

func Types() []ObjectType {
	object_types := []ObjectType{
		typeVm,
		typeLoadBalancer,
		typeNatService,
		typeSecurityGroup,
		typePublicIp,
		typeVolume,
		typeKeypair,
		typeRouteTable,
		typeInternetService,
		typeSubnet,
		typeNet,
		typeImage,
		typeSnapshot,
		typeVpnConnection,
		typeVirtualGateway,
		typeNic,
		typeAccessKey,
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
	_, err := provider.readAccountId()
	return err
}

func (provider *OutscaleOAPI) ReadObjects(typeName string) []Object {
	switch typeName {
	case typeVm:
		return provider.readVms()
	case typeLoadBalancer:
		return provider.readLoadBalancers()
	case typeNatService:
		return provider.readNatServices()
	case typeSecurityGroup:
		return provider.readSecurityGroups()
	case typePublicIp:
		return provider.readPublicIps()
	case typeVolume:
		return provider.readVolumes()
	case typeKeypair:
		return provider.readKeypairs()
	case typeRouteTable:
		return provider.readRouteTables()
	case typeInternetService:
		return provider.readInternetServices()
	case typeSubnet:
		return provider.readSubnets()
	case typeNet:
		return provider.readNets()
	case typeImage:
		return provider.readImages()
	case typeSnapshot:
		return provider.readSnapshots()
	case typeVpnConnection:
		return provider.readVpnConnections()
	case typeVirtualGateway:
		return provider.readVirtualGateways()
	case typeNic:
		return provider.readNics()
	case typeAccessKey:
		return provider.readAccessKeys()
	}
	return []Object{}
}

func (provider *OutscaleOAPI) DeleteObjects(typeName string, objects []Object) {
	switch typeName {
	case typeVm:
		provider.deleteVms(objects)
	case typeLoadBalancer:
		provider.deleteLoadBalancers(objects)
	case typeNatService:
		provider.deleteNatServices(objects)
	case typeSecurityGroup:
		provider.deleteSecurityGroups(objects)
	case typePublicIp:
		provider.deletePublicIps(objects)
	case typeVolume:
		provider.deleteVolumes(objects)
	case typeKeypair:
		provider.deleteKeypairs(objects)
	case typeRouteTable:
		provider.deleteRouteTables(objects)
	case typeInternetService:
		provider.deleteInternetServices(objects)
	case typeSubnet:
		provider.deleteSubnets(objects)
	case typeNet:
		provider.deleteNets(objects)
	case typeImage:
		provider.deleteImages(objects)
	case typeSnapshot:
		provider.deleteSnapshots(objects)
	case typeVpnConnection:
		provider.deleteVpnConnections(objects)
	case typeVirtualGateway:
		provider.deleteVirtualGateways(objects)
	case typeNic:
		provider.deleteNics(objects)
	case typeAccessKey:
		provider.deleteAccessKeys(objects)
	}
}

func (provider *OutscaleOAPI) StringObject(object string, typeName string) string {
	return object
}

func newAPICache() apiCache {
	return apiCache{
		internetServices: make(map[string]*osc.InternetService),
		publicIps:        make(map[string]*osc.PublicIp),
		vms:              make(map[string]*osc.Vm),
		nics:             make(map[string]*osc.Nic),
		routeTables:      make(map[string]*osc.RouteTable),
		securityGroups:   make(map[string]*osc.SecurityGroup),
	}
}

func (provider *OutscaleOAPI) readVms() []Object {
	vms := make([]Object, 0)
	read, httpRes, err := provider.client.VmApi.ReadVms(provider.context).
		ReadVmsRequest(osc.ReadVmsRequest{}).
		Execute()
	if err != nil {
		log.Println("Error while reading vms: ")
		if httpRes != nil {
			log.Println(httpRes.Status)
		}
		return vms
	}
	for i, vm := range *read.Vms {
		switch *vm.State {
		case "pending", "running", "stopping", "stopped", "shutting-down", "quarantine":
			vms = append(vms, *vm.VmId)
			provider.cache.vms[*vm.VmId] = &(*read.Vms)[i]
		}
	}
	return vms
}

func (provider *OutscaleOAPI) forceShutdownVms(vms []Object) {
	var vmsToForce []Object
	for _, vmId := range vms {
		vm := provider.cache.vms[vmId]
		if vm == nil {
			continue
		}
		switch *vm.State {
		case "pending", "running":
			vmsToForce = append(vmsToForce, vmId)
		}
	}
	log.Printf("Shutting down virtual machines: %s...\n", vmsToForce)
	forceStop := true
	stopOpts := osc.StopVmsRequest{
		VmIds:     vmsToForce,
		ForceStop: &forceStop,
	}
	_, httpRes, err := provider.client.VmApi.StopVms(provider.context).
		StopVmsRequest(stopOpts).
		Execute()
	if err != nil {
		log.Print("Error while shutting down vms: ")
		if httpRes != nil {
			log.Println(httpRes.Status)
		}
		return
	}
	log.Println("OK")
}

func (provider *OutscaleOAPI) deleteVms(vms []Object) {
	if len(vms) == 0 {
		return
	}
	provider.forceShutdownVms(vms)
	log.Printf("Deleting virtual machines: %s ... ", vms)
	deletionOpts := osc.DeleteVmsRequest{VmIds: vms}
	_, httpRes, err := provider.client.VmApi.DeleteVms(provider.context).
		DeleteVmsRequest(deletionOpts).
		Execute()
	if err != nil {
		log.Print("Error while deleting vms: ")
		if httpRes != nil {
			log.Println(httpRes.Status)
		}
	} else {
		log.Println("OK")
	}
}

func (provider *OutscaleOAPI) readLoadBalancers() []Object {
	loadBalancers := make([]Object, 0)
	read, httpRes, err := provider.client.LoadBalancerApi.ReadLoadBalancers(provider.context).
		ReadLoadBalancersRequest(osc.ReadLoadBalancersRequest{}).
		Execute()
	if err != nil {
		log.Print("Error while reading load balancers: ")
		if httpRes != nil {
			log.Println(httpRes.Status)
		}
		return loadBalancers
	}
	for _, loadBalancer := range *read.LoadBalancers {
		loadBalancers = append(loadBalancers, *loadBalancer.LoadBalancerName)
	}
	return loadBalancers
}

func (provider *OutscaleOAPI) deleteLoadBalancers(loadBalancers []Object) {
	if len(loadBalancers) == 0 {
		return
	}
	for _, loadBalancer := range loadBalancers {
		log.Printf("Deleting load balancer %s... ", loadBalancer)
		deletionOpts := osc.DeleteLoadBalancerRequest{LoadBalancerName: loadBalancer}
		_, httpRes, err := provider.client.LoadBalancerApi.
			DeleteLoadBalancer(provider.context).
			DeleteLoadBalancerRequest(deletionOpts).
			Execute()
		if err != nil {
			log.Print("Error while deleting load balancer: ")
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readNatServices() []Object {
	natServices := make([]Object, 0)
	read, httpRes, err := provider.client.NatServiceApi.ReadNatServices(provider.context).
		ReadNatServicesRequest(osc.ReadNatServicesRequest{}).
		Execute()
	if err != nil {
		log.Print("Error while reading nat services: ")
		if httpRes != nil {
			log.Println(httpRes.Status)
		}
		return natServices
	}
	for _, natService := range *read.NatServices {
		switch *natService.State {
		case "pending", "available":
			natServices = append(natServices, *natService.NatServiceId)
		}
	}
	return natServices
}

func (provider *OutscaleOAPI) deleteNatServices(natServices []Object) {
	if len(natServices) == 0 {
		return
	}
	for _, natService := range natServices {
		log.Printf("Deleting nat service %s... ", natService)
		deletionOpts := osc.DeleteNatServiceRequest{NatServiceId: natService}
		_, httpRes, err := provider.client.NatServiceApi.
			DeleteNatService(provider.context).
			DeleteNatServiceRequest(deletionOpts).
			Execute()
		if err != nil {
			log.Print("Error while deleting nat service: ")
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readSecurityGroups() []Object {
	securityGroups := make([]Object, 0)
	read, httpRes, err := provider.client.SecurityGroupApi.
		ReadSecurityGroups(provider.context).
		ReadSecurityGroupsRequest(osc.ReadSecurityGroupsRequest{}).
		Execute()
	if err != nil {
		log.Print("Error while reading security groups: ")
		if httpRes != nil {
			log.Println(httpRes.Status)
		}
		return securityGroups
	}
	for _, sg := range *read.SecurityGroups {
		if *sg.SecurityGroupName == "default" {
			continue
		}
		copySg := sg
		securityGroups = append(securityGroups, *sg.SecurityGroupId)
		provider.cache.securityGroups[*sg.SecurityGroupId] = &copySg
	}
	return securityGroups
}

func (provider *OutscaleOAPI) deleteSecurityGroupRules(securityGroupId string) error {
	securityGroup := provider.cache.securityGroups[securityGroupId]
	if securityGroup == nil || (securityGroup.InboundRules == nil && securityGroup.OutboundRules == nil) {
		return nil
	}

	if len(securityGroup.GetInboundRules()) != 0 {
		targetRules := []osc.SecurityGroupRule{}
		for _, rule := range securityGroup.GetInboundRules() {
			if len(rule.GetSecurityGroupsMembers()) == 0 {
				targetRules = append(targetRules, rule)
			}

			targetSecurityGroupMember := []osc.SecurityGroupsMember{}
			for _, sgMember := range rule.GetSecurityGroupsMembers() {
				sgMember.SetAccountId("")
				sgMember.SetSecurityGroupName("")
				targetSecurityGroupMember = append(targetSecurityGroupMember, sgMember)
			}

			rule.SetSecurityGroupsMembers(targetSecurityGroupMember)
			targetRules = append(targetRules, rule)

		}
		log.Printf("Deleting inbound security group rule from %s... ", securityGroupId)
		delete := osc.DeleteSecurityGroupRuleRequest{
			Flow:            "Inbound",
			Rules:           &targetRules,
			SecurityGroupId: securityGroupId,
		}

		_, httpRes, err := provider.client.SecurityGroupRuleApi.
			DeleteSecurityGroupRule(provider.context).
			DeleteSecurityGroupRuleRequest(delete).
			Execute()
		if err != nil {
			log.Printf("Error while deleting inbound rules of security group route %s: ", securityGroupId)
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
			return err
		} else {
			log.Println("OK")
		}
	}

	if len(securityGroup.GetOutboundRules()) != 0 {
		targetRules := []osc.SecurityGroupRule{}
		for _, rule := range securityGroup.GetOutboundRules() {
			if len(rule.GetSecurityGroupsMembers()) == 0 {
				targetRules = append(targetRules, rule)
			}

			targetSecurityGroupMember := []osc.SecurityGroupsMember{}
			for _, sgMember := range rule.GetSecurityGroupsMembers() {
				sgMember.SetAccountId("")
				sgMember.SetSecurityGroupName("")
				targetSecurityGroupMember = append(targetSecurityGroupMember, sgMember)
			}

			rule.SetSecurityGroupsMembers(targetSecurityGroupMember)
			targetRules = append(targetRules, rule)

		}
		log.Printf("Deleting outbound security group rule from %s... ", securityGroupId)
		delete := osc.DeleteSecurityGroupRuleRequest{
			Flow:            "Outbound",
			Rules:           &targetRules,
			SecurityGroupId: securityGroupId,
		}

		_, httpRes, err := provider.client.SecurityGroupRuleApi.
			DeleteSecurityGroupRule(provider.context).
			DeleteSecurityGroupRuleRequest(delete).
			Execute()
		if err != nil {
			log.Printf("Error while deleting outbound rules of security group route %s: ", securityGroupId)
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
			return err
		} else {
			log.Println("OK")
		}
	}
	return nil
}

func (provider *OutscaleOAPI) deleteSecurityGroups(securityGroups []Object) {
	if len(securityGroups) == 0 {
		return
	}
	for _, sg := range securityGroups {
		if provider.deleteSecurityGroupRules(sg) != nil {
			continue
		}
		log.Printf("Deleting security group %s... ", sg)
		deletionOpts := osc.DeleteSecurityGroupRequest{SecurityGroupId: &sg}
		_, httpRes, err := provider.client.SecurityGroupApi.
			DeleteSecurityGroup(provider.context).
			DeleteSecurityGroupRequest(deletionOpts).
			Execute()
		if err != nil {
			log.Print("Error while deleting security groups: ")
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readPublicIps() []Object {
	publicIps := make([]Object, 0)
	read, httpRes, err := provider.client.PublicIpApi.
		ReadPublicIps(provider.context).
		ReadPublicIpsRequest(osc.ReadPublicIpsRequest{}).
		Execute()
	if err != nil {
		log.Print("Error while reading public ips: ")
		if httpRes != nil {
			log.Println(httpRes.Status)
		}
		return publicIps
	}
	for i, pip := range *read.PublicIps {
		publicIps = append(publicIps, *pip.PublicIp)
		provider.cache.publicIps[*pip.PublicIp] = &(*read.PublicIps)[i]
	}
	return publicIps
}

func (provider *OutscaleOAPI) unlinkPublicIp(publicIP *string) error {
	cache := provider.cache.publicIps[*publicIP]
	if cache == nil {
		return nil
	}
	if cache.LinkPublicIpId == nil &&
		cache.NicId == nil &&
		cache.VmId == nil {
		return nil
	}
	log.Printf("Unlinking public ip %s... ", *publicIP)
	unlinkOpts := osc.UnlinkPublicIpRequest{PublicIp: publicIP}
	_, httpRes, err := provider.client.PublicIpApi.
		UnlinkPublicIp(provider.context).
		UnlinkPublicIpRequest(unlinkOpts).
		Execute()
	if err != nil {
		log.Print("Error while unlinking public ip: ")
		if httpRes != nil {
			log.Println(httpRes.Status)
		}
		return err
	}
	log.Println("OK")
	return nil
}

func (provider *OutscaleOAPI) deletePublicIps(publicIps []Object) {
	if len(publicIps) == 0 {
		return
	}
	for _, publicIP := range publicIps {
		if provider.unlinkPublicIp(&publicIP) != nil {
			continue
		}
		log.Printf("Deleting public ip %s... ", publicIP)
		deletionOpts := osc.DeletePublicIpRequest{PublicIp: &publicIP}
		_, httpRes, err := provider.client.PublicIpApi.
			DeletePublicIp(provider.context).
			DeletePublicIpRequest(deletionOpts).
			Execute()
		if err != nil {
			log.Print("Error while deleting public ip: ")
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readVolumes() []Object {
	volumes := make([]Object, 0)
	read, httpRes, err := provider.client.VolumeApi.
		ReadVolumes(provider.context).
		ReadVolumesRequest(osc.ReadVolumesRequest{}).
		Execute()
	if err != nil {
		log.Print("Error while reading volumes: ")
		if httpRes != nil {
			log.Println(httpRes.Status)
		}
		return volumes
	}
	for _, volume := range *read.Volumes {
		switch *volume.State {
		case "creating", "available", "in-use", "updating", "error":
			volumes = append(volumes, *volume.VolumeId)
		}
	}
	return volumes
}

func (provider *OutscaleOAPI) deleteVolumes(volumes []Object) {
	if len(volumes) == 0 {
		return
	}
	for _, volume := range volumes {
		log.Printf("Deleting volume %s... ", volume)
		deletionOpts := osc.DeleteVolumeRequest{VolumeId: volume}
		_, httpRes, err := provider.client.VolumeApi.
			DeleteVolume(provider.context).
			DeleteVolumeRequest(deletionOpts).
			Execute()
		if err != nil {
			log.Print("Error while deleting volume: ")
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readKeypairs() []Object {
	keypairs := make([]Object, 0)
	read, httpRes, err := provider.client.KeypairApi.ReadKeypairs(provider.context).
		ReadKeypairsRequest(osc.ReadKeypairsRequest{}).
		Execute()
	if err != nil {
		log.Print("Error while reading keypairs: ")
		if httpRes != nil {
			log.Println(httpRes.Status)
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
		log.Printf("Deleting keypair %s... ", keypair)
		deletionOpts := osc.DeleteKeypairRequest{KeypairName: keypair}
		_, httpRes, err := provider.client.KeypairApi.
			DeleteKeypair(provider.context).
			DeleteKeypairRequest(deletionOpts).
			Execute()
		if err != nil {
			log.Print("Error while deleting keypair: ")
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readRouteTables() []Object {
	routeTables := make([]Object, 0)
	read, httpRes, err := provider.client.RouteTableApi.ReadRouteTables(provider.context).
		ReadRouteTablesRequest(osc.ReadRouteTablesRequest{}).
		Execute()
	if err != nil {
		log.Print("Error while reading route tables: ")
		if httpRes != nil {
			log.Println(httpRes.Status)
		}
		return routeTables
	}
	for i, routeTable := range *read.RouteTables {
		routeTables = append(routeTables, *routeTable.RouteTableId)
		provider.cache.routeTables[*routeTable.RouteTableId] = &(*read.RouteTables)[i]
	}
	return routeTables
}

func (provider *OutscaleOAPI) unlinkRouteTable(RouteTableId string) error {
	routeTable := provider.cache.routeTables[RouteTableId]
	if routeTable == nil || routeTable.LinkRouteTables == nil {
		return nil
	}
	for _, link := range *routeTable.LinkRouteTables {
		if link.LinkRouteTableId == nil {
			continue
		}
		linkId := *link.LinkRouteTableId
		log.Printf("Unlinking route table %s (link %s)... ", RouteTableId, linkId)
		unlinkOps := osc.UnlinkRouteTableRequest{
			LinkRouteTableId: *link.LinkRouteTableId,
		}
		_, httpRes, err := provider.client.RouteTableApi.
			UnlinkRouteTable(provider.context).
			UnlinkRouteTableRequest(unlinkOps).
			Execute()
		if err != nil {
			log.Printf("Error while unlinking route table %s (links %s): ", RouteTableId, linkId)
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
			return err
		} else {
			log.Println("OK")
		}
	}
	return nil
}

func (provider *OutscaleOAPI) deleteRouteTables(routeTables []Object) {
	if len(routeTables) == 0 {
		return
	}
	for _, routeTable := range routeTables {
		if provider.unlinkRouteTable(routeTable) != nil {
			continue
		}
		log.Printf("Deleting route table %s... ", routeTable)
		deletionOpts := osc.DeleteRouteTableRequest{RouteTableId: routeTable}
		_, httpRes, err := provider.client.RouteTableApi.
			DeleteRouteTable(provider.context).
			DeleteRouteTableRequest(deletionOpts).
			Execute()
		if err != nil {
			log.Print("Error while deleting route table: ")
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readInternetServices() []Object {
	internetServices := make([]Object, 0)
	read, httpRes, err := provider.client.InternetServiceApi.ReadInternetServices(provider.context).
		ReadInternetServicesRequest(osc.ReadInternetServicesRequest{}).
		Execute()
	if err != nil {
		log.Print("Error while reading internet services: ")
		if httpRes != nil {
			log.Println(httpRes.Status)
		}
		return internetServices
	}
	for i, internetService := range *read.InternetServices {
		internetServices = append(internetServices, *internetService.InternetServiceId)
		provider.cache.internetServices[*internetService.InternetServiceId] = &(*read.InternetServices)[i]
	}
	return internetServices
}

func (provider *OutscaleOAPI) unlinkInternetSevice(internetServiceId string) error {
	internetService := provider.cache.internetServices[internetServiceId]
	if internetService == nil || internetService.NetId == nil {
		return nil
	}
	log.Printf("Unlinking internet service %s... ", internetServiceId)
	unlinkOps := osc.UnlinkInternetServiceRequest{
		InternetServiceId: internetServiceId,
		NetId:             *internetService.NetId,
	}
	_, httpRes, err := provider.client.InternetServiceApi.
		UnlinkInternetService(provider.context).
		UnlinkInternetServiceRequest(unlinkOps).
		Execute()
	if err != nil {
		log.Print("Error while unlinking internet service: ")
		if httpRes != nil {
			log.Println(httpRes.Status)
		}
		return err
	} else {
		log.Println("OK")
	}
	return nil
}

func (provider *OutscaleOAPI) deleteInternetServices(internetServices []Object) {
	if len(internetServices) == 0 {
		return
	}
	for _, internetService := range internetServices {
		if provider.unlinkInternetSevice(internetService) != nil {
			continue
		}
		log.Printf("Deleting internet service %s... ", internetService)
		deletionOpts := osc.DeleteInternetServiceRequest{InternetServiceId: internetService}
		_, httpRes, err := provider.client.InternetServiceApi.
			DeleteInternetService(provider.context).
			DeleteInternetServiceRequest(deletionOpts).
			Execute()
		if err != nil {
			log.Print("Error while deleting internet service: ")
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readSubnets() []Object {
	subnets := make([]Object, 0)
	read, httpRes, err := provider.client.SubnetApi.ReadSubnets(provider.context).
		ReadSubnetsRequest(osc.ReadSubnetsRequest{}).
		Execute()
	if err != nil {
		log.Print("Error while reading subnets: ")
		if httpRes != nil {
			log.Println(httpRes.Status)
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
		log.Printf("Deleting subnet %s... ", subnet)
		deletionOpts := osc.DeleteSubnetRequest{SubnetId: subnet}
		_, httpRes, err := provider.client.SubnetApi.
			DeleteSubnet(provider.context).
			DeleteSubnetRequest(deletionOpts).
			Execute()
		if err != nil {
			log.Print("Error while deleting subnet: ")
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readNets() []Object {
	nets := make([]Object, 0)
	read, httpRes, err := provider.client.NetApi.ReadNets(provider.context).
		ReadNetsRequest(osc.ReadNetsRequest{}).
		Execute()
	if err != nil {
		log.Print("Error while reading nets: ")
		if httpRes != nil {
			log.Println(httpRes.Status)
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
		log.Printf("Deleting net %s... ", net)
		deletionOpts := osc.DeleteNetRequest{NetId: net}
		_, httpRes, err := provider.client.NetApi.
			DeleteNet(provider.context).
			DeleteNetRequest(deletionOpts).
			Execute()
		if err != nil {
			log.Print("Error while deleting net: ")
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readAccountId() (*string, error) {
	if provider.cache.accountId == nil {
		read, httpRes, err := provider.client.AccountApi.ReadAccounts(provider.context).
			ReadAccountsRequest(osc.ReadAccountsRequest{}).
			Execute()
		if err != nil {
			log.Print("Error while reading account: ")
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
			return nil, err
		}
		if len(*read.Accounts) == 0 {
			log.Println("Error while reading account: no account listed")
			return nil, err
		}
		provider.cache.accountId = (*read.Accounts)[0].AccountId
	}
	return provider.cache.accountId, nil
}

func (provider *OutscaleOAPI) readImages() []Object {
	images := make([]Object, 0)
	accountId, err := provider.readAccountId()
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
		log.Print("Error while reading images: ")
		if httpRes != nil {
			log.Println(httpRes.Status)
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
		log.Printf("Deleting image %s... ", image)
		deletionOpts := osc.DeleteImageRequest{ImageId: image}
		_, httpRes, err := provider.client.ImageApi.
			DeleteImage(provider.context).
			DeleteImageRequest(deletionOpts).
			Execute()
		if err != nil {
			log.Print("Error while deleting image: ")
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readSnapshots() []Object {
	snapshots := make([]Object, 0)
	accountId, err := provider.readAccountId()
	if err != nil {
		return snapshots
	}
	var accountIds []string
	accountIds = append(accountIds, *accountId)
	read, httpRes, err := provider.client.SnapshotApi.ReadSnapshots(provider.context).
		ReadSnapshotsRequest(osc.ReadSnapshotsRequest{
			Filters: &osc.FiltersSnapshot{
				AccountIds: &accountIds,
			},
		}).
		Execute()
	if err != nil {
		log.Print("Error while reading snapshots: ")
		if httpRes != nil {
			log.Println(httpRes.Status)
		}
		return snapshots
	}
	for _, snapshot := range *read.Snapshots {
		snapshots = append(snapshots, *snapshot.SnapshotId)
	}
	return snapshots
}

func (provider *OutscaleOAPI) deleteSnapshots(snapshots []Object) {
	if len(snapshots) == 0 {
		return
	}
	for _, snapshot := range snapshots {
		log.Printf("Deleting snapshot %s... ", snapshot)
		deletionOpts := osc.DeleteSnapshotRequest{SnapshotId: snapshot}
		_, httpRes, err := provider.client.SnapshotApi.
			DeleteSnapshot(provider.context).
			DeleteSnapshotRequest(deletionOpts).
			Execute()
		if err != nil {
			log.Print("Error while deleting snapshot: ")
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readVpnConnections() []Object {
	vpnConnections := make([]Object, 0)
	read, httpRes, err := provider.client.VpnConnectionApi.ReadVpnConnections(provider.context).
		ReadVpnConnectionsRequest(osc.ReadVpnConnectionsRequest{}).
		Execute()
	if err != nil {
		log.Print("Error while reading vpn connections: ")
		if httpRes != nil {
			log.Println(httpRes.Status)
		}
		return vpnConnections
	}
	for _, vpnConnection := range *read.VpnConnections {
		switch *vpnConnection.State {
		case "pending", "available":
			vpnConnections = append(vpnConnections, *vpnConnection.VpnConnectionId)
		}
	}
	return vpnConnections
}

func (provider *OutscaleOAPI) deleteVpnConnections(vpnConnections []Object) {
	if len(vpnConnections) == 0 {
		return
	}
	for _, vpnConnection := range vpnConnections {
		log.Printf("Deleting vpn connection %s... ", vpnConnection)
		deletionOpts := osc.DeleteVpnConnectionRequest{VpnConnectionId: vpnConnection}
		_, httpRes, err := provider.client.VpnConnectionApi.
			DeleteVpnConnection(provider.context).
			DeleteVpnConnectionRequest(deletionOpts).
			Execute()
		if err != nil {
			log.Print("Error while deleting vpn connection: ")
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readVirtualGateways() []Object {
	virtualGateways := make([]Object, 0)
	read, httpRes, err := provider.client.VirtualGatewayApi.ReadVirtualGateways(provider.context).
		ReadVirtualGatewaysRequest(osc.ReadVirtualGatewaysRequest{}).
		Execute()
	if err != nil {
		log.Print("Error while reading virtual gateways: ")
		if httpRes != nil {
			log.Println(httpRes.Status)
		}
		return virtualGateways
	}
	for _, virtualGateway := range *read.VirtualGateways {
		switch *virtualGateway.State {
		case "pending", "available":
			virtualGateways = append(virtualGateways, *virtualGateway.VirtualGatewayId)
		}
	}
	return virtualGateways
}

func (provider *OutscaleOAPI) deleteVirtualGateways(virtualGateways []Object) {
	if len(virtualGateways) == 0 {
		return
	}
	for _, virtualGateway := range virtualGateways {
		log.Printf("Deleting virtual gateway %s... ", virtualGateway)
		deletionOpts := osc.DeleteVirtualGatewayRequest{VirtualGatewayId: virtualGateway}
		_, httpRes, err := provider.client.VirtualGatewayApi.
			DeleteVirtualGateway(provider.context).
			DeleteVirtualGatewayRequest(deletionOpts).
			Execute()
		if err != nil {
			log.Print("Error while deleting virtual gateway: ")
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readNics() []Object {
	nics := make([]Object, 0)
	read, httpRes, err := provider.client.NicApi.ReadNics(provider.context).
		ReadNicsRequest(osc.ReadNicsRequest{}).
		Execute()
	if err != nil {
		log.Print("Error while reading nics: ")
		if httpRes != nil {
			log.Println(httpRes.Status)
		}
		return nics
	}
	for i, nic := range *read.Nics {
		nics = append(nics, *nic.NicId)
		provider.cache.nics[*nic.NicId] = &(*read.Nics)[i]
	}
	return nics
}

func (provider *OutscaleOAPI) unlinkNics(nics []Object) {
	for _, nicId := range nics {
		nic := provider.cache.nics[nicId]
		if nic == nil {
			continue
		}
		switch *nic.State {
		case "attaching", "in-use":
		default:
			continue
		}
		if nic.LinkNic == nil || nic.LinkNic.LinkNicId == nil {
			continue
		}
		log.Printf("Unlinking nic %s... ", nicId)
		unlinkOpts := osc.UnlinkNicRequest{LinkNicId: *nic.LinkNic.LinkNicId}
		_, httpRes, err := provider.client.NicApi.
			UnlinkNic(provider.context).
			UnlinkNicRequest(unlinkOpts).
			Execute()
		if err != nil {
			log.Print("Error while unlinking nic: ")
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
			continue
		}
		log.Println("OK")
	}
}

func (provider *OutscaleOAPI) deleteNics(nics []Object) {
	if len(nics) == 0 {
		return
	}
	provider.unlinkNics(nics)
	for _, nicId := range nics {
		log.Printf("Deleting nic %s... ", nicId)
		deletionOpts := osc.DeleteNicRequest{NicId: nicId}
		_, httpRes, err := provider.client.NicApi.
			DeleteNic(provider.context).
			DeleteNicRequest(deletionOpts).
			Execute()
		if err != nil {
			log.Print("Error while deleting nic: ")
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readAccessKeys() []Object {
	accessKeys := make([]Object, 0)
	read, httpRes, err := provider.client.AccessKeyApi.ReadAccessKeys(provider.context).
		ReadAccessKeysRequest(osc.ReadAccessKeysRequest{}).
		Execute()
	if err != nil {
		log.Print("Error while reading access keys: ")
		if httpRes != nil {
			log.Println(httpRes.Status)
		}
		return accessKeys
	}
	for _, accessKey := range *read.AccessKeys {
		if *accessKey.State == "ACTIVE" {
			accessKeys = append(accessKeys, *accessKey.AccessKeyId)
		}
	}
	return accessKeys
}

func (provider *OutscaleOAPI) deleteAccessKeys(accessKeys []Object) {
	if len(accessKeys) == 0 {
		return
	}
	for _, accessKey := range accessKeys {
		log.Printf("Deleting access key %s... ", accessKey)
		deletionOpts := osc.DeleteAccessKeyRequest{AccessKeyId: accessKey}
		_, httpRes, err := provider.client.AccessKeyApi.
			DeleteAccessKey(provider.context).
			DeleteAccessKeyRequest(deletionOpts).
			Execute()
		if err != nil {
			log.Print("Error while deleting access key: ")
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
		} else {
			log.Println("OK")
		}
	}
}
