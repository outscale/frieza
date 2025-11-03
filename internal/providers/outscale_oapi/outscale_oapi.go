package outscale_oapi

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	. "github.com/outscale/frieza/internal/common"
	osc "github.com/outscale/osc-sdk-go/v2"
	"github.com/teris-io/cli"
)

const (
	Name = "outscale_oapi"

	typeVm              = "vm"
	typeLoadBalancer    = "load_balancer"
	typeNatService      = "nat_service"
	typeSecurityGroup   = "security_group"
	typePublicIp        = "public_ip"
	typeVolume          = "volume"
	typeKeypair         = "keypair"
	typeRouteTable      = "route_table"
	typeInternetService = "internet_service"
	typeSubnet          = "subnet"
	typeNet             = "net"
	typeImage           = "image"
	typeSnapshot        = "snapshot"
	typeVpnConnection   = "vpn_connection"
	typeVirtualGateway  = "virtual_gateway"
	typeNic             = "nic"
	typeAccessKey       = "access_key"
	typeNetAccessPoint  = "net_access_point"
	typeNetPeering      = "net_peering"
)

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
	ctx = context.WithValue(
		ctx,
		osc.ContextServerVariables,
		map[string]string{"region": config["region"]},
	)
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
		typeSecurityGroup,
		typeInternetService,
		typeRouteTable,
		typeNatService,
		typeNic,
		typeVpnConnection,
		typeVirtualGateway,
		typePublicIp,
		typeNetAccessPoint,
		typeNetPeering,
		typeSubnet,
		typeNet,
		typeVolume,
		typeImage,
		typeSnapshot,
		typeKeypair,
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

func (provider *OutscaleOAPI) ReadObjects(typeName string) ([]Object, error) {
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
	case typeNetAccessPoint:
		return provider.readNetAccessPoints()
	case typeNetPeering:
		return provider.readNetPeerings()
	}
	return []Object{}, nil
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
	case typeNetAccessPoint:
		provider.deleteNetAccessPoints(objects)
	case typeNetPeering:
		provider.deleteNetPeerings(objects)
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

func (provider *OutscaleOAPI) readVms() ([]Object, error) {
	vms := make([]Object, 0)
	read, httpRes, err := provider.client.VmApi.ReadVms(provider.context).
		ReadVmsRequest(osc.ReadVmsRequest{
			Filters: &osc.FiltersVm{
				VmStateNames: &[]string{
					"pending", "running", "stopping", "stopped", "shutting-down", "quarantine", // skipping terminated
				},
			},
		}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("read vms: %w", getErrorInfo(err, httpRes))
	}
	for i, vm := range *read.Vms {
		vms = append(vms, *vm.VmId)
		provider.cache.vms[*vm.VmId] = &(*read.Vms)[i]
	}
	return vms, nil
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
		log.Printf("Error while shutting down vms: %v\n", getErrorInfo(err, httpRes))
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
		log.Printf("Error while deleting vms: %v\n", getErrorInfo(err, httpRes))
	} else {
		log.Println("OK")
	}
}

func (provider *OutscaleOAPI) readLoadBalancers() ([]Object, error) {
	loadBalancers := make([]Object, 0)
	read, httpRes, err := provider.client.LoadBalancerApi.ReadLoadBalancers(provider.context).
		ReadLoadBalancersRequest(osc.ReadLoadBalancersRequest{}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("read load balancers: %w", getErrorInfo(err, httpRes))
	}
	for _, loadBalancer := range *read.LoadBalancers {
		loadBalancers = append(loadBalancers, *loadBalancer.LoadBalancerName)
	}
	return loadBalancers, nil
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
			log.Printf("Error while deleting load balancer: %v\n", getErrorInfo(err, httpRes))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readNatServices() ([]Object, error) {
	natServices := make([]Object, 0)
	read, httpRes, err := provider.client.NatServiceApi.ReadNatServices(provider.context).
		ReadNatServicesRequest(osc.ReadNatServicesRequest{
			Filters: &osc.FiltersNatService{
				States: &[]string{
					"pending", "available", // skipping deleting, deleted
				},
			},
		}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("read nat: %w", getErrorInfo(err, httpRes))
	}
	for _, natService := range *read.NatServices {
		natServices = append(natServices, *natService.NatServiceId)
	}
	return natServices, nil
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
			log.Printf("Error while deleting nat service: %v\n", getErrorInfo(err, httpRes))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readSecurityGroups() ([]Object, error) {
	securityGroups := make([]Object, 0)
	read, httpRes, err := provider.client.SecurityGroupApi.
		ReadSecurityGroups(provider.context).
		ReadSecurityGroupsRequest(osc.ReadSecurityGroupsRequest{}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("read security groups: %w", getErrorInfo(err, httpRes))
	}
	for _, sg := range *read.SecurityGroups {
		if *sg.SecurityGroupName == "default" {
			continue
		}
		copySg := sg
		securityGroups = append(securityGroups, *sg.SecurityGroupId)
		provider.cache.securityGroups[*sg.SecurityGroupId] = &copySg
	}
	return securityGroups, nil
}

func (provider *OutscaleOAPI) deleteSecurityGroupRules(securityGroupId string) error {
	securityGroup := provider.cache.securityGroups[securityGroupId]
	if securityGroup == nil ||
		(securityGroup.InboundRules == nil && securityGroup.OutboundRules == nil) {
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
			log.Printf(
				"Error while deleting inbound rules of security group route %s: ",
				securityGroupId,
			)
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
			log.Printf(
				"Error while deleting outbound rules of security group route %s: ",
				securityGroupId,
			)
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
			log.Printf("Error while deleting security groups: %v\n", getErrorInfo(err, httpRes))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readPublicIps() ([]Object, error) {
	publicIps := make([]Object, 0)
	read, httpRes, err := provider.client.PublicIpApi.
		ReadPublicIps(provider.context).
		ReadPublicIpsRequest(osc.ReadPublicIpsRequest{}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("read public ips: %w", getErrorInfo(err, httpRes))
	}
	for i, pip := range *read.PublicIps {
		publicIps = append(publicIps, *pip.PublicIp)
		provider.cache.publicIps[*pip.PublicIp] = &(*read.PublicIps)[i]
	}
	return publicIps, nil
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
		log.Printf("Error while unlinking public ip: %v\n", getErrorInfo(err, httpRes))
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
			log.Printf("Error while deleting public ip: %v\n", getErrorInfo(err, httpRes))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readVolumes() ([]Object, error) {
	volumes := make([]Object, 0)
	read, httpRes, err := provider.client.VolumeApi.
		ReadVolumes(provider.context).
		ReadVolumesRequest(osc.ReadVolumesRequest{
			Filters: &osc.FiltersVolume{
				VolumeStates: &[]string{
					"creating", "available", "in-use", "error", // skip deleting
				},
			},
		}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("read volumes: %w", getErrorInfo(err, httpRes))
	}
	for _, volume := range *read.Volumes {
		volumes = append(volumes, *volume.VolumeId)
	}
	return volumes, nil
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
			log.Printf("Error while deleting volume: %v\n", getErrorInfo(err, httpRes))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readKeypairs() ([]Object, error) {
	keypairs := make([]Object, 0)
	read, httpRes, err := provider.client.KeypairApi.ReadKeypairs(provider.context).
		ReadKeypairsRequest(osc.ReadKeypairsRequest{}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("read key pairs: %w", getErrorInfo(err, httpRes))
	}
	for _, keypair := range *read.Keypairs {
		keypairs = append(keypairs, *keypair.KeypairName)
	}
	return keypairs, nil
}

func (provider *OutscaleOAPI) deleteKeypairs(keypairs []Object) {
	if len(keypairs) == 0 {
		return
	}
	for _, keypair := range keypairs {
		log.Printf("Deleting keypair %s... ", keypair)
		deletionOpts := osc.DeleteKeypairRequest{KeypairName: &keypair}
		_, httpRes, err := provider.client.KeypairApi.
			DeleteKeypair(provider.context).
			DeleteKeypairRequest(deletionOpts).
			Execute()
		if err != nil {
			log.Printf("Error while deleting keypair: %v\n", getErrorInfo(err, httpRes))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readRouteTables() ([]Object, error) {
	routeTables := make([]Object, 0)
	read, httpRes, err := provider.client.RouteTableApi.ReadRouteTables(provider.context).
		ReadRouteTablesRequest(osc.ReadRouteTablesRequest{}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("read route tables: %w", getErrorInfo(err, httpRes))
	}
	for i, routeTable := range *read.RouteTables {
		if provider.isMainRouteTable(&routeTable) {
			continue
		}
		routeTables = append(routeTables, *routeTable.RouteTableId)
		provider.cache.routeTables[*routeTable.RouteTableId] = &(*read.RouteTables)[i]
	}
	return routeTables, nil
}

func (provider *OutscaleOAPI) unlinkRouteTable(RouteTableId string) error {
	routeTable := provider.cache.routeTables[RouteTableId]
	if routeTable == nil || routeTable.LinkRouteTables == nil {
		return nil
	}
	for _, link := range *routeTable.LinkRouteTables {
		if link.LinkRouteTableId == nil || link.GetMain() {
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
			log.Printf(
				"Error while unlinking route table %s (links %s): %v\n",
				RouteTableId,
				linkId,
				getErrorInfo(err, httpRes),
			)
			return err
		} else {
			log.Println("OK")
		}
	}
	return nil
}

func (provider *OutscaleOAPI) isMainRouteTable(routeTable *osc.RouteTable) bool {
	for _, link := range *routeTable.LinkRouteTables {
		if link.GetMain() {
			return true
		}
	}
	return false
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
			log.Printf("Error while deleting route table: %v\n", getErrorInfo(err, httpRes))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readInternetServices() ([]Object, error) {
	internetServices := make([]Object, 0)
	read, httpRes, err := provider.client.InternetServiceApi.ReadInternetServices(provider.context).
		ReadInternetServicesRequest(osc.ReadInternetServicesRequest{}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("read internet service: %w", getErrorInfo(err, httpRes))
	}
	for i, internetService := range *read.InternetServices {
		internetServices = append(internetServices, *internetService.InternetServiceId)
		provider.cache.internetServices[*internetService.InternetServiceId] = &(*read.InternetServices)[i]
	}
	return internetServices, nil
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
		log.Printf("Error while unlinking internet service: %v\n", getErrorInfo(err, httpRes))
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
			log.Printf("Error while deleting internet service: %v\n", getErrorInfo(err, httpRes))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readSubnets() ([]Object, error) {
	subnets := make([]Object, 0)
	read, httpRes, err := provider.client.SubnetApi.ReadSubnets(provider.context).
		ReadSubnetsRequest(osc.ReadSubnetsRequest{}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("read subnets: %w", getErrorInfo(err, httpRes))
	}
	for _, subnet := range *read.Subnets {
		subnets = append(subnets, *subnet.SubnetId)
	}
	return subnets, nil
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
			log.Printf("Error while deleting subnet: %v\n", getErrorInfo(err, httpRes))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readNets() ([]Object, error) {
	nets := make([]Object, 0)
	read, httpRes, err := provider.client.NetApi.ReadNets(provider.context).
		ReadNetsRequest(osc.ReadNetsRequest{
			Filters: &osc.FiltersNet{
				States: &[]string{"pending", "available"}, // skipping deleting
			},
		}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("read nets: %w", getErrorInfo(err, httpRes))
	}
	for _, net := range *read.Nets {
		nets = append(nets, *net.NetId)
	}
	return nets, nil
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
			log.Printf("Error while deleting net: %v\n", getErrorInfo(err, httpRes))
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
			return nil, fmt.Errorf("read vms: %w", getErrorInfo(err, httpRes))
		}
		if len(*read.Accounts) == 0 {
			log.Println("Error while reading account: no account listed")
			return nil, err
		}
		provider.cache.accountId = (*read.Accounts)[0].AccountId
	}
	return provider.cache.accountId, nil
}

func (provider *OutscaleOAPI) readImages() ([]Object, error) {
	images := make([]Object, 0)
	accountId, err := provider.readAccountId()
	if err != nil {
		return images, nil
	}
	var accountIds []string
	accountIds = append(accountIds, *accountId)
	read, httpRes, err := provider.client.ImageApi.ReadImages(provider.context).
		ReadImagesRequest(osc.ReadImagesRequest{
			Filters: &osc.FiltersImage{
				AccountIds: &accountIds,
			},
		}).
		Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading images: %v\n", getErrorInfo(err, httpRes))
		return nil, fmt.Errorf("read images: %w")
	}
	for _, image := range *read.Images {
		images = append(images, *image.ImageId)
	}
	return images, nil
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
			log.Printf("Error while deleting image: %v\n", getErrorInfo(err, httpRes))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readSnapshots() ([]Object, error) {
	snapshots := make([]Object, 0)
	accountId, err := provider.readAccountId()
	if err != nil {
		return snapshots, nil
	}
	var accountIds []string
	accountIds = append(accountIds, *accountId)
	read, httpRes, err := provider.client.SnapshotApi.ReadSnapshots(provider.context).
		ReadSnapshotsRequest(osc.ReadSnapshotsRequest{
			Filters: &osc.FiltersSnapshot{
				AccountIds: &accountIds,
				States: &[]string{
					"in-queue", "pending", "completed", "error", // skipping deleting
				},
			},
		}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("read snapshots: %w", getErrorInfo(err, httpRes))
	}
	for _, snapshot := range *read.Snapshots {
		snapshots = append(snapshots, *snapshot.SnapshotId)
	}
	return snapshots, nil
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
			log.Printf("Error while deleting snapshot: %v\n", getErrorInfo(err, httpRes))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readVpnConnections() ([]Object, error) {
	vpnConnections := make([]Object, 0)
	read, httpRes, err := provider.client.VpnConnectionApi.ReadVpnConnections(provider.context).
		ReadVpnConnectionsRequest(osc.ReadVpnConnectionsRequest{
			Filters: &osc.FiltersVpnConnection{
				States: &[]string{
					"pending", "available", // skipping deleting, deleted
				},
			},
		}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("read vpn connections: %w", getErrorInfo(err, httpRes))
	}
	for _, vpnConnection := range *read.VpnConnections {
		vpnConnections = append(vpnConnections, *vpnConnection.VpnConnectionId)
	}
	return vpnConnections, nil
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
			log.Printf("Error while deleting vpn connection: %v\n", getErrorInfo(err, httpRes))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readVirtualGateways() ([]Object, error) {
	virtualGateways := make([]Object, 0)
	read, httpRes, err := provider.client.VirtualGatewayApi.ReadVirtualGateways(provider.context).
		ReadVirtualGatewaysRequest(osc.ReadVirtualGatewaysRequest{
			Filters: &osc.FiltersVirtualGateway{
				States: &[]string{
					"pending", "available", // skipping deleting, deleted
				},
			},
		}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("read virtual gateways: %w", getErrorInfo(err, httpRes))
	}
	for _, virtualGateway := range *read.VirtualGateways {
		virtualGateways = append(virtualGateways, *virtualGateway.VirtualGatewayId)
	}
	return virtualGateways, nil
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
			log.Printf("Error while deleting virtual gateway: %v\n", getErrorInfo(err, httpRes))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readNics() ([]Object, error) {
	nics := make([]Object, 0)
	read, httpRes, err := provider.client.NicApi.ReadNics(provider.context).
		ReadNicsRequest(osc.ReadNicsRequest{}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("read nics: %w", getErrorInfo(err, httpRes))
	}
	for i, nic := range *read.Nics {
		nics = append(nics, *nic.NicId)
		provider.cache.nics[*nic.NicId] = &(*read.Nics)[i]
	}
	return nics, nil
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
			log.Printf("Error while unlinking nic: %v\n", getErrorInfo(err, httpRes))
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
			log.Printf("Error while deleting nic: %v\n", getErrorInfo(err, httpRes))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readAccessKeys() ([]Object, error) {
	accessKeys := make([]Object, 0)
	read, httpRes, err := provider.client.AccessKeyApi.ReadAccessKeys(provider.context).
		ReadAccessKeysRequest(osc.ReadAccessKeysRequest{}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("read ak: %w", getErrorInfo(err, httpRes))
	}
	for _, accessKey := range *read.AccessKeys {
		if *accessKey.State == "ACTIVE" {
			accessKeys = append(accessKeys, *accessKey.AccessKeyId)
		}
	}
	return accessKeys, nil
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
			log.Printf("Error while deleting access key: %v\n", getErrorInfo(err, httpRes))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readNetAccessPoints() ([]Object, error) {
	netAccessPoints := make([]Object, 0)
	read, httpRes, err := provider.client.NetAccessPointApi.ReadNetAccessPoints(provider.context).
		ReadNetAccessPointsRequest(osc.ReadNetAccessPointsRequest{
			Filters: &osc.FiltersNetAccessPoint{
				States: &[]string{
					"pending", "available", // skipping deleting, deleted
				},
			},
		}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("read net access points: %w", getErrorInfo(err, httpRes))
	}
	for _, netAccessPoint := range *read.NetAccessPoints {
		netAccessPoints = append(netAccessPoints, *netAccessPoint.NetAccessPointId)
	}
	return netAccessPoints, nil
}

func (provider *OutscaleOAPI) deleteNetAccessPoints(netAccessPoints []Object) {
	if len(netAccessPoints) == 0 {
		return
	}
	for _, netAccessPoint := range netAccessPoints {
		log.Printf("Deleting net access point %s... ", netAccessPoint)
		deletionOpts := osc.DeleteNetAccessPointRequest{NetAccessPointId: netAccessPoint}
		_, httpRes, err := provider.client.NetAccessPointApi.
			DeleteNetAccessPoint(provider.context).
			DeleteNetAccessPointRequest(deletionOpts).
			Execute()
		if err != nil {
			log.Print("Error while deleting net access point: ")
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readNetPeerings() ([]Object, error) {
	netPeerings := make([]Object, 0)
	read, httpRes, err := provider.client.NetPeeringApi.ReadNetPeerings(provider.context).
		ReadNetPeeringsRequest(osc.ReadNetPeeringsRequest{
			Filters: &osc.FiltersNetPeering{
				StateNames: &[]string{
					"pending-acceptance", "active", "rejected", "failed", "expired", // skipping deleted
				},
			},
		}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("read net peerings: %w", getErrorInfo(err, httpRes))
	}
	for _, netPeering := range *read.NetPeerings {
		netPeerings = append(netPeerings, *netPeering.NetPeeringId)
	}
	return netPeerings, nil
}

func (provider *OutscaleOAPI) deleteNetPeerings(netPeerings []Object) {
	if len(netPeerings) == 0 {
		return
	}
	for _, netPeering := range netPeerings {
		log.Printf("Deleting net peering %s... ", netPeering)
		deletionOpts := osc.DeleteNetPeeringRequest{NetPeeringId: netPeering}
		_, httpRes, err := provider.client.NetPeeringApi.
			DeleteNetPeering(provider.context).
			DeleteNetPeeringRequest(deletionOpts).
			Execute()
		if err != nil {
			log.Print("Error while deleting net peering: ")
			if httpRes != nil {
				log.Println(httpRes.Status)
			}
		} else {
			log.Println("OK")
		}
	}
}
