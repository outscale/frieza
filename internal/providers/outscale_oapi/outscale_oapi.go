package outscale_oapi

import (
	"context"
	"fmt"
	"log"
	"os"

	. "github.com/outscale/frieza/internal/common"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/osc-sdk-go/v3/pkg/profile"
	oscutils "github.com/outscale/osc-sdk-go/v3/pkg/utils"
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
	typeClientGateway   = "client_gateway"
	typeNic             = "nic"
	typeAccessKey       = "access_key"
	typeNetAccessPoint  = "net_access_point"
	typeNetPeering      = "net_peering"
	typeUser            = "user"
	typeUserAccessKey   = "user_access_key"
	typePolicy          = "policy"
	typePolicyLink      = "policy_link"
	typePolicyVersion   = "policy_version"
	typeFlexibleGpu     = "flexible_gpu"
)

type OutscaleOAPI struct {
	client *osc.Client
	cache  apiCache
}

type apiCache struct {
	accountId        *string
	internetServices map[Object]*osc.InternetService
	publicIps        map[Object]*osc.PublicIp
	vms              map[Object]*osc.Vm
	nics             map[Object]*osc.Nic
	routeTables      map[Object]*osc.RouteTable
	securityGroups   map[Object]*osc.SecurityGroup
	flexibleGpus     map[Object]*osc.FlexibleGpu
}

func New(config ProviderConfig, debug bool) (*OutscaleOAPI, error) {
	profileName := config["profile"]
	profilePath := config["path"]
	profile, err := profile.NewFrom(profileName, profilePath)
	if err != nil {
		return nil, err
	}

	if ak, ok := config["ak"]; ok {
		profile.AccessKey = ak
	}

	if sk, ok := config["sk"]; ok {
		profile.SecretKey = sk
	}

	if region, ok := config["region"]; ok {
		profile.Region = region
	}

	client, err := osc.NewClient(profile, oscutils.WithUseragent("frieza/"+FullVersion()))
	if err != nil {
		return nil, err
	}

	return &OutscaleOAPI{
		client: client,
		cache:  newAPICache(),
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
		typeClientGateway,
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
		typeUserAccessKey,
		typeUser,
		typePolicyLink,
		typePolicy,
		typePolicyVersion,
		typeFlexibleGpu,
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
	case typeClientGateway:
		return provider.readClientGateways()
	case typeNic:
		return provider.readNics()
	case typeAccessKey:
		return provider.readAccessKeys()
	case typeNetAccessPoint:
		return provider.readNetAccessPoints()
	case typeNetPeering:
		return provider.readNetPeerings()
	case typeUser:
		return provider.readUsers()
	case typeUserAccessKey:
		return provider.readUserAccessKeys()
	case typePolicy:
		return provider.readPolicies()
	case typePolicyLink:
		return provider.readPolicyLinks()
	case typePolicyVersion:
		return provider.readPolicyVersions()
	case typeFlexibleGpu:
		return provider.readFlexibleGpus()
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
	case typeClientGateway:
		provider.deleteClientGateways(objects)
	case typeNic:
		provider.deleteNics(objects)
	case typeAccessKey:
		provider.deleteAccessKeys(objects)
	case typeNetAccessPoint:
		provider.deleteNetAccessPoints(objects)
	case typeNetPeering:
		provider.deleteNetPeerings(objects)
	case typeUser:
		provider.deleteUsers(objects)
	case typeUserAccessKey:
		provider.deleteUserAccessKeys(objects)
	case typePolicy:
		provider.deletePolicies(objects)
	case typePolicyLink:
		provider.deletePolicyLinks(objects)
	case typePolicyVersion:
		provider.deletePolicyVersions(objects)
	case typeFlexibleGpu:
		provider.deleteFlexibleGpus(objects)
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
	read, err := provider.client.ReadVms(context.Background(), osc.ReadVmsRequest{
		Filters: &osc.FiltersVm{
			VmStateNames: &[]osc.VmState{
				"pending", "running", "stopping", "stopped", "shutting-down", "quarantine", // skipping terminated
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("read vms: %w", getErrorInfo(err))
	}
	for i, vm := range *read.Vms {
		vms = append(vms, vm.VmId)
		provider.cache.vms[vm.VmId] = &(*read.Vms)[i]
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
		switch vm.State {
		case "pending", "running":
			vmsToForce = append(vmsToForce, vmId)
		}
	}
	if len(vmsToForce) == 0 {
		return
	}

	log.Printf("Shutting down virtual machines: %v...\n", vmsToForce)
	forceStop := true
	stopOpts := osc.StopVmsRequest{
		VmIds:     vmsToForce,
		ForceStop: &forceStop,
	}
	_, err := provider.client.StopVms(context.Background(), stopOpts)
	if err != nil {
		log.Printf("Error while shutting down vms: %v\n", getErrorInfo(err))
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
	_, err := provider.client.DeleteVms(context.Background(), deletionOpts)
	if err != nil {
		log.Printf("Error while deleting vms: %v\n", getErrorInfo(err))
	} else {
		log.Println("OK")
	}
}

func (provider *OutscaleOAPI) readLoadBalancers() ([]Object, error) {
	loadBalancers := make([]Object, 0)
	read, err := provider.client.ReadLoadBalancers(
		context.Background(),
		osc.ReadLoadBalancersRequest{},
	)
	if err != nil {
		return nil, fmt.Errorf("read load balancers: %w", getErrorInfo(err))
	}
	for _, loadBalancer := range *read.LoadBalancers {
		loadBalancers = append(loadBalancers, loadBalancer.LoadBalancerName)
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
		_, err := provider.client.DeleteLoadBalancer(context.Background(), deletionOpts)
		if err != nil {
			log.Printf("Error while deleting load balancer: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readNatServices() ([]Object, error) {
	natServices := make([]Object, 0)
	read, err := provider.client.ReadNatServices(
		context.Background(),
		osc.ReadNatServicesRequest{
			Filters: &osc.FiltersNatService{
				States: &[]osc.NatServiceState{
					"pending", "available", // skipping deleting, deleted
				},
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("read nat: %w", getErrorInfo(err))
	}
	for _, natService := range *read.NatServices {
		natServices = append(natServices, natService.NatServiceId)
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
		_, err := provider.client.DeleteNatService(context.Background(), deletionOpts)
		if err != nil {
			log.Printf("Error while deleting nat service: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readSecurityGroups() ([]Object, error) {
	securityGroups := make([]Object, 0)
	read, err := provider.client.ReadSecurityGroups(
		context.Background(),
		osc.ReadSecurityGroupsRequest{},
	)
	if err != nil {
		return nil, fmt.Errorf("read security groups: %w", getErrorInfo(err))
	}
	for _, sg := range *read.SecurityGroups {
		if sg.SecurityGroupName == "default" {
			continue
		}
		copySg := sg
		securityGroups = append(securityGroups, sg.SecurityGroupId)
		provider.cache.securityGroups[sg.SecurityGroupId] = &copySg
	}
	return securityGroups, nil
}

func ptrString(value string) *string {
	return &value
}

func (provider *OutscaleOAPI) deleteSecurityGroupRules(securityGroupId string) error {
	securityGroup := provider.cache.securityGroups[securityGroupId]
	if securityGroup == nil ||
		(securityGroup.InboundRules == nil && securityGroup.OutboundRules == nil) {
		return nil
	}

	if len(securityGroup.InboundRules) != 0 {
		targetRules := []osc.SecurityGroupRule{}
		for _, rule := range securityGroup.InboundRules {
			if len(rule.SecurityGroupsMembers) == 0 {
				targetRules = append(targetRules, rule)
			}

			targetSecurityGroupMember := []osc.SecurityGroupsMember{}
			for _, sgMember := range rule.SecurityGroupsMembers {
				sgMember.AccountId = ptrString("")
				sgMember.SecurityGroupName = ptrString("")
				targetSecurityGroupMember = append(targetSecurityGroupMember, sgMember)
			}

			rule.SecurityGroupsMembers = targetSecurityGroupMember
			targetRules = append(targetRules, rule)

		}
		log.Printf("Deleting inbound security group rule from %s... ", securityGroupId)
		delete := osc.DeleteSecurityGroupRuleRequest{
			Flow:            "Inbound",
			Rules:           targetRules,
			SecurityGroupId: securityGroupId,
		}

		_, err := provider.client.DeleteSecurityGroupRule(context.Background(), delete)
		if err != nil {
			log.Printf(
				"Error while deleting inbound rules of security group route %s: ",
				securityGroupId,
			)
			return err
		} else {
			log.Println("OK")
		}
	}

	if len(securityGroup.OutboundRules) != 0 {
		targetRules := []osc.SecurityGroupRule{}
		for _, rule := range securityGroup.OutboundRules {
			if len(rule.SecurityGroupsMembers) == 0 {
				targetRules = append(targetRules, rule)
			}

			targetSecurityGroupMember := []osc.SecurityGroupsMember{}
			for _, sgMember := range rule.SecurityGroupsMembers {
				sgMember.AccountId = ptrString("")
				sgMember.SecurityGroupName = ptrString("")
				targetSecurityGroupMember = append(targetSecurityGroupMember, sgMember)
			}

			rule.SecurityGroupsMembers = targetSecurityGroupMember
			targetRules = append(targetRules, rule)

		}
		log.Printf("Deleting outbound security group rule from %s... ", securityGroupId)
		delete := osc.DeleteSecurityGroupRuleRequest{
			Flow:            "Outbound",
			Rules:           targetRules,
			SecurityGroupId: securityGroupId,
		}

		_, err := provider.client.DeleteSecurityGroupRule(context.Background(), delete)
		if err != nil {
			log.Printf(
				"Error while deleting outbound rules of security group route %s: ",
				securityGroupId,
			)
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
		_, err := provider.client.DeleteSecurityGroup(context.Background(), deletionOpts)
		if err != nil {
			log.Printf("Error while deleting security groups: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readPublicIps() ([]Object, error) {
	publicIps := make([]Object, 0)
	read, err := provider.client.ReadPublicIps(
		context.Background(),
		osc.ReadPublicIpsRequest{},
	)
	if err != nil {
		return nil, fmt.Errorf("read public ips: %w", getErrorInfo(err))
	}
	for i, pip := range *read.PublicIps {
		publicIps = append(publicIps, pip.PublicIp)
		provider.cache.publicIps[pip.PublicIp] = &(*read.PublicIps)[i]
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
	_, err := provider.client.UnlinkPublicIp(context.Background(), unlinkOpts)
	if err != nil {
		log.Printf("Error while unlinking public ip: %v\n", getErrorInfo(err))
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
		_, err := provider.client.DeletePublicIp(context.Background(), deletionOpts)
		if err != nil {
			log.Printf("Error while deleting public ip: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readVolumes() ([]Object, error) {
	volumes := make([]Object, 0)
	read, err := provider.client.ReadVolumes(context.Background(), osc.ReadVolumesRequest{
		Filters: &osc.FiltersVolume{
			VolumeStates: &[]osc.VolumeState{
				"creating", "available", "in-use", "error",
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("read volumes: %w", getErrorInfo(err))
	}
	for _, volume := range *read.Volumes {
		// When a volume created from a snapshot is in the deleting state,
		// it will be returned even if the "deleting" filter is missing from Filters.VolumeStates
		if volume.State == "deleting" {
			continue
		}
		volumes = append(volumes, volume.VolumeId)
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
		_, err := provider.client.DeleteVolume(context.Background(), deletionOpts)
		if err != nil {
			log.Printf("Error while deleting volume: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readKeypairs() ([]Object, error) {
	keypairs := make([]Object, 0)
	read, err := provider.client.ReadKeypairs(context.Background(), osc.ReadKeypairsRequest{})
	if err != nil {
		return nil, fmt.Errorf("read key pairs: %w", getErrorInfo(err))
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
		_, err := provider.client.DeleteKeypair(context.Background(), deletionOpts)
		if err != nil {
			log.Printf("Error while deleting keypair: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readRouteTables() ([]Object, error) {
	routeTables := make([]Object, 0)
	read, err := provider.client.ReadRouteTables(
		context.Background(),
		osc.ReadRouteTablesRequest{},
	)
	if err != nil {
		return nil, fmt.Errorf("read route tables: %w", getErrorInfo(err))
	}
	for i, routeTable := range *read.RouteTables {
		if provider.isMainRouteTable(&routeTable) {
			continue
		}
		routeTables = append(routeTables, routeTable.RouteTableId)
		provider.cache.routeTables[routeTable.RouteTableId] = &(*read.RouteTables)[i]
	}
	return routeTables, nil
}

func (provider *OutscaleOAPI) unlinkRouteTable(RouteTableId string) error {
	routeTable := provider.cache.routeTables[RouteTableId]
	if routeTable == nil || routeTable.LinkRouteTables == nil {
		return nil
	}
	for _, link := range routeTable.LinkRouteTables {
		if link.Main {
			continue
		}
		linkId := link.LinkRouteTableId
		log.Printf("Unlinking route table %s (link %s)... ", RouteTableId, linkId)
		unlinkOps := osc.UnlinkRouteTableRequest{
			LinkRouteTableId: link.LinkRouteTableId,
		}
		_, err := provider.client.UnlinkRouteTable(context.Background(), unlinkOps)
		if err != nil {
			log.Printf(
				"Error while unlinking route table %s (links %s): %v\n",
				RouteTableId,
				linkId,
				getErrorInfo(err),
			)
			return err
		} else {
			log.Println("OK")
		}
	}
	return nil
}

func (provider *OutscaleOAPI) isMainRouteTable(routeTable *osc.RouteTable) bool {
	for _, link := range routeTable.LinkRouteTables {
		if link.Main {
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
		_, err := provider.client.DeleteRouteTable(context.Background(), deletionOpts)
		if err != nil {
			log.Printf("Error while deleting route table: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readInternetServices() ([]Object, error) {
	internetServices := make([]Object, 0)
	read, err := provider.client.ReadInternetServices(
		context.Background(),
		osc.ReadInternetServicesRequest{},
	)
	if err != nil {
		return nil, fmt.Errorf("read internet service: %w", getErrorInfo(err))
	}
	for i, internetService := range *read.InternetServices {
		internetServices = append(internetServices, internetService.InternetServiceId)
		provider.cache.internetServices[internetService.InternetServiceId] = &(*read.InternetServices)[i]
	}
	return internetServices, nil
}

func (provider *OutscaleOAPI) unlinkInternetSevice(internetServiceId string) error {
	internetService := provider.cache.internetServices[internetServiceId]
	if internetService == nil {
		return nil
	}
	log.Printf("Unlinking internet service %s... ", internetServiceId)
	unlinkOps := osc.UnlinkInternetServiceRequest{
		InternetServiceId: internetServiceId,
		NetId:             internetService.NetId,
	}
	_, err := provider.client.UnlinkInternetService(context.Background(), unlinkOps)
	if err != nil {
		log.Printf("Error while unlinking internet service: %v\n", getErrorInfo(err))
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
		_, err := provider.client.DeleteInternetService(context.Background(), deletionOpts)
		if err != nil {
			log.Printf("Error while deleting internet service: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readSubnets() ([]Object, error) {
	subnets := make([]Object, 0)
	read, err := provider.client.ReadSubnets(context.Background(), osc.ReadSubnetsRequest{})
	if err != nil {
		return nil, fmt.Errorf("read subnets: %w", getErrorInfo(err))
	}
	for _, subnet := range *read.Subnets {
		subnets = append(subnets, subnet.SubnetId)
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
		_, err := provider.client.DeleteSubnet(context.Background(), deletionOpts)
		if err != nil {
			log.Printf("Error while deleting subnet: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readNets() ([]Object, error) {
	nets := make([]Object, 0)
	read, err := provider.client.ReadNets(context.Background(), osc.ReadNetsRequest{
		Filters: &osc.FiltersNet{
			States: &[]osc.NetState{"pending", "available"}, // skipping deleting
		},
	})
	if err != nil {
		return nil, fmt.Errorf("read nets: %w", getErrorInfo(err))
	}
	for _, net := range *read.Nets {
		nets = append(nets, net.NetId)
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
		_, err := provider.client.DeleteNet(context.Background(), deletionOpts)
		if err != nil {
			log.Printf("Error while deleting net: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readAccountId() (*string, error) {
	if provider.cache.accountId == nil {
		read, err := provider.client.ReadAccounts(
			context.Background(),
			osc.ReadAccountsRequest{},
		)
		if err != nil {
			return nil, fmt.Errorf("read vms: %w", getErrorInfo(err))
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
	read, err := provider.client.ReadImages(context.Background(), osc.ReadImagesRequest{
		Filters: &osc.FiltersImage{
			AccountIds: &accountIds,
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading images: %v\n", getErrorInfo(err))
		return nil, fmt.Errorf("read images: %w", err)
	}
	for _, image := range *read.Images {
		images = append(images, image.ImageId)
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
		_, err := provider.client.DeleteImage(context.Background(), deletionOpts)
		if err != nil {
			log.Printf("Error while deleting image: %v\n", getErrorInfo(err))
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
	read, err := provider.client.ReadSnapshots(context.Background(), osc.ReadSnapshotsRequest{
		Filters: &osc.FiltersSnapshot{
			AccountIds: &accountIds,
			States: &[]osc.SnapshotState{
				"in-queue", "pending", "completed", "error", // skipping deleting
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("read snapshots: %w", getErrorInfo(err))
	}
	for _, snapshot := range *read.Snapshots {
		snapshots = append(snapshots, snapshot.SnapshotId)
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
		_, err := provider.client.DeleteSnapshot(context.Background(), deletionOpts)
		if err != nil {
			log.Printf("Error while deleting snapshot: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readVpnConnections() ([]Object, error) {
	vpnConnections := make([]Object, 0)
	read, err := provider.client.ReadVpnConnections(
		context.Background(),
		osc.ReadVpnConnectionsRequest{
			Filters: &osc.FiltersVpnConnection{
				States: &[]string{
					"pending", "available", // skipping deleting, deleted
				},
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("read vpn connections: %w", getErrorInfo(err))
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
		_, err := provider.client.DeleteVpnConnection(context.Background(), deletionOpts)
		if err != nil {
			log.Printf("Error while deleting vpn connection: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readVirtualGateways() ([]Object, error) {
	virtualGateways := make([]Object, 0)
	read, err := provider.client.ReadVirtualGateways(
		context.Background(),
		osc.ReadVirtualGatewaysRequest{
			Filters: &osc.FiltersVirtualGateway{
				States: &[]string{
					"pending", "available", // skipping deleting, deleted
				},
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("read virtual gateways: %w", getErrorInfo(err))
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
		_, err := provider.client.DeleteVirtualGateway(context.Background(), deletionOpts)
		if err != nil {
			log.Printf("Error while deleting virtual gateway: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readClientGateways() ([]Object, error) {
	clientGateways := make([]Object, 0)
	read, err := provider.client.ReadClientGateways(
		context.Background(),
		osc.ReadClientGatewaysRequest{
			Filters: &osc.FiltersClientGateway{
				States: &[]string{
					"pending", "available", // skipping deleting, deleted
				},
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("read client gateways: %w", getErrorInfo(err))
	}
	for _, clientGateway := range *read.ClientGateways {
		clientGateways = append(clientGateways, *clientGateway.ClientGatewayId)
	}
	return clientGateways, nil
}

func (provider *OutscaleOAPI) deleteClientGateways(clientGateways []Object) {
	if len(clientGateways) == 0 {
		return
	}
	for _, clientGateway := range clientGateways {
		log.Printf("Deleting client gateway %s... ", clientGateway)
		deletionOpts := osc.DeleteClientGatewayRequest{ClientGatewayId: clientGateway}
		_, err := provider.client.DeleteClientGateway(context.Background(), deletionOpts)
		if err != nil {
			log.Printf("Error while deleting client gateway: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readNics() ([]Object, error) {
	nics := make([]Object, 0)
	read, err := provider.client.ReadNics(context.Background(), osc.ReadNicsRequest{})
	if err != nil {
		return nil, fmt.Errorf("read nics: %w", getErrorInfo(err))
	}
	for i, nic := range *read.Nics {
		nics = append(nics, nic.NicId)
		provider.cache.nics[nic.NicId] = &(*read.Nics)[i]
	}
	return nics, nil
}

func (provider *OutscaleOAPI) unlinkNics(nics []Object) {
	for _, nicId := range nics {
		nic := provider.cache.nics[nicId]
		if nic == nil {
			continue
		}
		switch nic.State {
		case "attaching", "in-use":
		default:
			continue
		}
		if nic.LinkNic == nil {
			continue
		}
		log.Printf("Unlinking nic %s... ", nicId)
		unlinkOpts := osc.UnlinkNicRequest{LinkNicId: nic.LinkNic.LinkNicId}
		_, err := provider.client.UnlinkNic(context.Background(), unlinkOpts)
		if err != nil {
			log.Printf("Error while unlinking nic: %v\n", getErrorInfo(err))
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
		_, err := provider.client.DeleteNic(context.Background(), deletionOpts)
		if err != nil {
			log.Printf("Error while deleting nic: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readAccessKeys() ([]Object, error) {
	accessKeys := make([]Object, 0)
	read, err := provider.client.ReadAccessKeys(
		context.Background(),
		osc.ReadAccessKeysRequest{},
	)
	if err != nil {
		return nil, fmt.Errorf("read ak: %w", getErrorInfo(err))
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
		_, err := provider.client.DeleteAccessKey(context.Background(), deletionOpts)
		if err != nil {
			log.Printf("Error while deleting access key: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readNetAccessPoints() ([]Object, error) {
	netAccessPoints := make([]Object, 0)
	read, err := provider.client.ReadNetAccessPoints(
		context.Background(),
		osc.ReadNetAccessPointsRequest{
			Filters: &osc.FiltersNetAccessPoint{
				States: &[]osc.NetAccessPointState{
					"pending", "available", // skipping deleting, deleted
				},
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("read net access points: %w", getErrorInfo(err))
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
		_, err := provider.client.DeleteNetAccessPoint(context.Background(), deletionOpts)
		if err != nil {
			log.Print("Error while deleting net access point: ")
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readNetPeerings() ([]Object, error) {
	netPeerings := make([]Object, 0)
	read, err := provider.client.ReadNetPeerings(
		context.Background(),
		osc.ReadNetPeeringsRequest{
			Filters: &osc.FiltersNetPeering{
				StateNames: &[]osc.NetPeeringStateName{
					"pending-acceptance", "active", "rejected", "failed", "expired", // skipping deleted
				},
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("read net peerings: %w", getErrorInfo(err))
	}
	for _, netPeering := range *read.NetPeerings {
		netPeerings = append(netPeerings, netPeering.NetPeeringId)
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
		_, err := provider.client.DeleteNetPeering(context.Background(), deletionOpts)
		if err != nil {
			log.Print("Error while deleting net peering: %w", err)
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readUsers() ([]Object, error) {
	users := make([]Object, 0)
	read, err := provider.client.ReadUsers(context.Background(), osc.ReadUsersRequest{})
	if err != nil {
		return nil, fmt.Errorf("read users: %w", getErrorInfo(err))
	}
	for _, user := range *read.Users {
		users = append(users, *user.UserName)
	}
	return users, nil
}

func (provider *OutscaleOAPI) deleteUsers(users []Object) {
	if len(users) == 0 {
		return
	}
	for _, user := range users {
		log.Printf("Deleting user %s... ", user)
		deleteOpts := osc.DeleteUserRequest{UserName: user}
		_, err := provider.client.DeleteUser(context.Background(), deleteOpts)
		if err != nil {
			log.Print("Error while deleting user: %w", err)
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readUserAccessKeys() ([]Object, error) {
	accessKeys := make([]Object, 0)

	readUser, err := provider.client.ReadUsers(context.Background(), osc.ReadUsersRequest{})
	if err != nil {
		return nil, fmt.Errorf("read users: %w", getErrorInfo(err))
	}

	for _, user := range *readUser.Users {
		read, err := provider.client.ReadAccessKeys(
			context.Background(),
			osc.ReadAccessKeysRequest{UserName: user.UserName},
		)
		if err != nil {
			return nil, fmt.Errorf("read user ak: %w", getErrorInfo(err))
		}
		for _, accessKey := range *read.AccessKeys {
			if *accessKey.State == "ACTIVE" {
				composed := fmt.Sprintf("%s,%s", *user.UserName, *accessKey.AccessKeyId)
				accessKeys = append(accessKeys, composed)
			}
		}
	}
	return accessKeys, nil
}

func (provider *OutscaleOAPI) deleteUserAccessKeys(accessKeys []Object) {
	if len(accessKeys) == 0 {
		return
	}
	for _, accessKey := range accessKeys {
		var userName, accessKeyId string
		log.Printf("Deleting user access key %s... ", accessKey)
		_, err := fmt.Scanf("%s,%s", &userName, &accessKeyId)
		if err != nil {
			continue
		}

		deletionOpts := osc.DeleteAccessKeyRequest{AccessKeyId: accessKeyId, UserName: &userName}
		_, err = provider.client.DeleteAccessKey(context.Background(), deletionOpts)
		if err != nil {
			log.Printf("Error while deleting user access key: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readPolicies() ([]Object, error) {
	policies := make([]Object, 0)

	read, err := provider.client.ReadPolicies(context.Background(), osc.ReadPoliciesRequest{})
	if err != nil {
		return nil, fmt.Errorf("read policies: %w", getErrorInfo(err))
	}
	for _, policy := range *read.Policies {
		policies = append(policies, *policy.Orn)
	}
	return policies, nil
}

func (provider *OutscaleOAPI) deletePolicies(policies []Object) {
	if len(policies) == 0 {
		return
	}
	for _, policy := range policies {
		log.Printf("Deleting policy %s... ", policy)
		deleteOpts := osc.DeletePolicyRequest{PolicyOrn: policy}
		_, err := provider.client.DeletePolicy(context.Background(), deleteOpts)
		if err != nil {
			log.Print("Error while deleting policy: %w", err)
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readPolicyLinks() ([]Object, error) {
	policyLinks := make([]Object, 0)

	read, err := provider.client.ReadPolicies(context.Background(), osc.ReadPoliciesRequest{})
	if err != nil {
		return nil, fmt.Errorf("read policies: %w", getErrorInfo(err))
	}
	for _, policy := range *read.Policies {
		read, err := provider.client.ReadEntitiesLinkedToPolicy(
			context.Background(),
			osc.ReadEntitiesLinkedToPolicyRequest{
				EntitiesType: &[]osc.ReadEntitiesLinkedToPolicyRequestEntitiesType{"USER", "GROUP"},
				PolicyOrn:    policy.Orn,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("read policy links: %w", getErrorInfo(err))
		}
		for _, policyLink := range *read.PolicyEntities.Groups {
			policyLinks = append(
				policyLinks,
				fmt.Sprintf("GROUP,%s,%s", *policy.Orn, *policyLink.Name),
			)
		}
		for _, policyLink := range *read.PolicyEntities.Users {
			policyLinks = append(
				policyLinks,
				fmt.Sprintf("USER,%s,%s", *policy.Orn, *policyLink.Name),
			)
		}
	}
	return policyLinks, nil
}

func (provider *OutscaleOAPI) deletePolicyLinks(policyLinks []Object) {
	if len(policyLinks) == 0 {
		return
	}

	for _, policylink := range policyLinks {
		var linkType, policyOrn, linkName string
		log.Printf("Deleting policy link %s... ", policylink)
		_, err := fmt.Scanf("%s,%s,%s", &linkType, &policyOrn, &linkName)
		if err != nil {
			continue
		}

		switch linkType {
		case "USER":
			deleteOpts := osc.UnlinkPolicyRequest{
				PolicyOrn: policyOrn,
				UserName:  linkName,
			}
			_, err := provider.client.UnlinkPolicy(context.Background(), deleteOpts)
			if err != nil {
				log.Print("Error while unlinking policy: %w", err)
			}

		case "GROUP":
			deleteOpts := osc.UnlinkManagedPolicyFromUserGroupRequest{
				PolicyOrn:     policyOrn,
				UserGroupName: linkName,
			}
			_, err := provider.client.UnlinkManagedPolicyFromUserGroup(
				context.Background(),
				deleteOpts,
			)
			if err != nil {
				log.Print("Error while unlinking policy: %w", err)
			}
		}
	}
}

func (provider *OutscaleOAPI) readPolicyVersions() ([]Object, error) {
	policyVersions := make([]Object, 0)

	read, err := provider.client.ReadPolicies(context.Background(), osc.ReadPoliciesRequest{})
	if err != nil {
		return nil, fmt.Errorf("read policies: %w", getErrorInfo(err))
	}
	for _, policy := range *read.Policies {
		read, err := provider.client.ReadPolicyVersions(
			context.Background(),
			osc.ReadPolicyVersionsRequest{
				PolicyOrn: *policy.Orn,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("read policy version: %w", getErrorInfo(err))
		}
		for _, policyVersion := range *read.PolicyVersions {
			if *policyVersion.DefaultVersion {
				continue
			}

			policyVersions = append(
				policyVersions,
				fmt.Sprintf("%s,%s", *policy.Orn, *policyVersion.VersionId),
			)
		}
	}
	return policyVersions, nil
}

func (provider *OutscaleOAPI) deletePolicyVersions(policyVersions []Object) {
	if len(policyVersions) == 0 {
		return
	}

	for _, policyVersion := range policyVersions {
		var policyOrn, version string
		log.Printf("Deleting policy version %s... ", policyVersion)
		_, err := fmt.Scanf("%s,%s", &policyOrn, &version)
		if err != nil {
			continue
		}

		deleteOpts := osc.DeletePolicyVersionRequest{
			PolicyOrn: policyOrn,
			VersionId: version,
		}
		_, err = provider.client.DeletePolicyVersion(context.Background(), deleteOpts)
		if err != nil {
			log.Print("Error while deleting policy version: %w", err)
		}

	}
}

func (provider *OutscaleOAPI) readFlexibleGpus() ([]Object, error) {
	flexibleGpus := make([]Object, 0)

	read, err := provider.client.ReadFlexibleGpus(context.Background(), osc.ReadFlexibleGpusRequest{})
	if err != nil {
		return nil, fmt.Errorf("read flexible gpus: %w", getErrorInfo(err))
	}
	for i, gpu := range *read.FlexibleGpus {
		flexibleGpus = append(flexibleGpus, *gpu.FlexibleGpuId)
		provider.cache.flexibleGpus[*gpu.FlexibleGpuId] = &(*read.FlexibleGpus)[i]
	}
	return flexibleGpus, nil
}

func (provider *OutscaleOAPI) unlinkFlexibleGpus(flexibleGpus []Object) {
	for _, gpuObj := range flexibleGpus {
		gpu := provider.cache.flexibleGpus[gpuObj]
		if gpu == nil {
			continue
		}
		switch *gpu.State {
		case "attaching", "attached":
		default:
			continue
		}
		log.Printf("Unlinking flexible gpu %s... ", *gpu.FlexibleGpuId)
		unlinkOpts := osc.UnlinkFlexibleGpuRequest{FlexibleGpuId: *gpu.FlexibleGpuId}
		_, err := provider.client.UnlinkFlexibleGpu(context.Background(), unlinkOpts)
		if err != nil {
			log.Printf("Error while unlinking flexible gpu: %v\n", getErrorInfo(err))
			continue
		}
		log.Println("OK")
	}
}

func (provider *OutscaleOAPI) deleteFlexibleGpus(flexibleGpus []Object) {
	if len(flexibleGpus) == 0 {
		return
	}
	provider.unlinkFlexibleGpus(flexibleGpus)
	for _, gpu := range flexibleGpus {
		log.Printf("Releasing flexible gpu %s... ", gpu)
		deleteOpts := osc.DeleteFlexibleGpuRequest{FlexibleGpuId: gpu}
		_, err := provider.client.DeleteFlexibleGpu(context.Background(), deleteOpts)
		if err != nil {
			log.Print("Error while deleting flexible gpu: %w", err)
		} else {
			log.Println("OK")
		}
	}
}
