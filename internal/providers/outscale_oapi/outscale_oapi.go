package outscale_oapi

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	. "github.com/outscale/frieza/internal/common"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/osc-sdk-go/v3/pkg/profile"
	"github.com/teris-io/cli"
)

const (
	Name = "outscale_oapi"

	typeVm                = "vm"
	typeLoadBalancer      = "load_balancer"
	typeNatService        = "nat_service"
	typeSecurityGroup     = "security_group"
	typePublicIp          = "public_ip"
	typeVolume            = "volume"
	typeKeypair           = "keypair"
	typeRouteTable        = "route_table"
	typeInternetService   = "internet_service"
	typeSubnet            = "subnet"
	typeNet               = "net"
	typeImage             = "image"
	typeSnapshot          = "snapshot"
	typeVpnConnection     = "vpn_connection"
	typeVirtualGateway    = "virtual_gateway"
	typeClientGateway     = "client_gateway"
	typeNic               = "nic"
	typeAccessKey         = "access_key"
	typeNetAccessPoint    = "net_access_point"
	typeNetPeering        = "net_peering"
	typeUser              = "user"
	typeUserAccessKey     = "user_access_key"
	typePolicy            = "policy"
	typePolicyLink        = "policy_link"
	typePolicyVersion     = "policy_version"
	typeFlexibleGpu       = "flexible_gpu"
	typeCa                = "ca"
	typeServerCertificate = "server_certificate"
	typeDhcpOption        = "dhcp_option"
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

	client, err := osc.NewClient(profile, options.WithUseragent("frieza/"+FullVersion()))
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
		typeCa,
		typeServerCertificate,
		typeDhcpOption,
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

func (provider *OutscaleOAPI) AuthTest(ctx context.Context) error {
	_, err := provider.readAccountId(ctx)
	return err
}

func (provider *OutscaleOAPI) ReadObjects(ctx context.Context, typeName string) ([]Object, error) {
	switch typeName {
	case typeVm:
		return provider.readVms(ctx)
	case typeLoadBalancer:
		return provider.readLoadBalancers(ctx)
	case typeNatService:
		return provider.readNatServices(ctx)
	case typeSecurityGroup:
		return provider.readSecurityGroups(ctx)
	case typePublicIp:
		return provider.readPublicIps(ctx)
	case typeVolume:
		return provider.readVolumes(ctx)
	case typeKeypair:
		return provider.readKeypairs(ctx)
	case typeRouteTable:
		return provider.readRouteTables(ctx)
	case typeInternetService:
		return provider.readInternetServices(ctx)
	case typeSubnet:
		return provider.readSubnets(ctx)
	case typeNet:
		return provider.readNets(ctx)
	case typeImage:
		return provider.readImages(ctx)
	case typeSnapshot:
		return provider.readSnapshots(ctx)
	case typeVpnConnection:
		return provider.readVpnConnections(ctx)
	case typeVirtualGateway:
		return provider.readVirtualGateways(ctx)
	case typeClientGateway:
		return provider.readClientGateways(ctx)
	case typeNic:
		return provider.readNics(ctx)
	case typeAccessKey:
		return provider.readAccessKeys(ctx)
	case typeNetAccessPoint:
		return provider.readNetAccessPoints(ctx)
	case typeNetPeering:
		return provider.readNetPeerings(ctx)
	case typeUser:
		return provider.readUsers(ctx)
	case typeUserAccessKey:
		return provider.readUserAccessKeys(ctx)
	case typePolicy:
		return provider.readPolicies(ctx)
	case typePolicyLink:
		return provider.readPolicyLinks(ctx)
	case typePolicyVersion:
		return provider.readPolicyVersions(ctx)
	case typeFlexibleGpu:
		return provider.readFlexibleGpus(ctx)
	case typeCa:
		return provider.readCas(ctx)
	case typeServerCertificate:
		return provider.readServerCertificates(ctx)
	case typeDhcpOption:
		return provider.readDhcpOptions(ctx)
	}
	return []Object{}, nil
}

func (provider *OutscaleOAPI) DeleteObjects(ctx context.Context, typeName string, objects []Object) {
	switch typeName {
	case typeVm:
		provider.deleteVms(ctx, objects)
	case typeLoadBalancer:
		provider.deleteLoadBalancers(ctx, objects)
	case typeNatService:
		provider.deleteNatServices(ctx, objects)
	case typeSecurityGroup:
		provider.deleteSecurityGroups(ctx, objects)
	case typePublicIp:
		provider.deletePublicIps(ctx, objects)
	case typeVolume:
		provider.deleteVolumes(ctx, objects)
	case typeKeypair:
		provider.deleteKeypairs(ctx, objects)
	case typeRouteTable:
		provider.deleteRouteTables(ctx, objects)
	case typeInternetService:
		provider.deleteInternetServices(ctx, objects)
	case typeSubnet:
		provider.deleteSubnets(ctx, objects)
	case typeNet:
		provider.deleteNets(ctx, objects)
	case typeImage:
		provider.deleteImages(ctx, objects)
	case typeSnapshot:
		provider.deleteSnapshots(ctx, objects)
	case typeVpnConnection:
		provider.deleteVpnConnections(ctx, objects)
	case typeVirtualGateway:
		provider.deleteVirtualGateways(ctx, objects)
	case typeClientGateway:
		provider.deleteClientGateways(ctx, objects)
	case typeNic:
		provider.deleteNics(ctx, objects)
	case typeAccessKey:
		provider.deleteAccessKeys(ctx, objects)
	case typeNetAccessPoint:
		provider.deleteNetAccessPoints(ctx, objects)
	case typeNetPeering:
		provider.deleteNetPeerings(ctx, objects)
	case typeUser:
		provider.deleteUsers(ctx, objects)
	case typeUserAccessKey:
		provider.deleteUserAccessKeys(ctx, objects)
	case typePolicy:
		provider.deletePolicies(ctx, objects)
	case typePolicyLink:
		provider.deletePolicyLinks(ctx, objects)
	case typePolicyVersion:
		provider.deletePolicyVersions(ctx, objects)
	case typeFlexibleGpu:
		provider.deleteFlexibleGpus(ctx, objects)
	case typeCa:
		provider.deleteCas(ctx, objects)
	case typeServerCertificate:
		provider.deleteServerCertificates(ctx, objects)
	case typeDhcpOption:
		provider.deleteDhcpOptions(ctx, objects)
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
		flexibleGpus:     make(map[string]*osc.FlexibleGpu),
	}
}

func (provider *OutscaleOAPI) readVms(ctx context.Context) ([]Object, error) {
	vms := make([]Object, 0)
	read, err := provider.client.ReadVms(ctx, osc.ReadVmsRequest{
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

func (provider *OutscaleOAPI) forceShutdownVms(ctx context.Context, vms []Object) {
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
	_, err := provider.client.StopVms(ctx, stopOpts)
	if err != nil {
		log.Printf("Error while shutting down vms: %v\n", getErrorInfo(err))
		return
	}
	log.Println("OK")
}

func (provider *OutscaleOAPI) deleteVms(ctx context.Context, vms []Object) {
	if len(vms) == 0 {
		return
	}
	provider.forceShutdownVms(ctx, vms)
	log.Printf("Deleting virtual machines: %s ... ", vms)
	deletionOpts := osc.DeleteVmsRequest{VmIds: vms}
	_, err := provider.client.DeleteVms(ctx, deletionOpts)
	if err != nil {
		log.Printf("Error while deleting vms: %v\n", getErrorInfo(err))
	} else {
		log.Println("OK")
	}
}

func (provider *OutscaleOAPI) readLoadBalancers(ctx context.Context) ([]Object, error) {
	loadBalancers := make([]Object, 0)
	read, err := provider.client.ReadLoadBalancers(
		ctx,
		osc.ReadLoadBalancersRequest{
			Filters: &osc.FiltersLoadBalancer{
				States: &[]osc.LoadBalancerState{ // skipping deleted, deleting
					osc.LoadBalancerStateActive, osc.LoadBalancerStateProvisioning, osc.LoadBalancerStateReconfiguring, osc.LoadBalancerStateReloading, osc.LoadBalancerStateStarting,
				},
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("read load balancers: %w", getErrorInfo(err))
	}
	for _, loadBalancer := range *read.LoadBalancers {
		loadBalancers = append(loadBalancers, loadBalancer.LoadBalancerName)
	}
	return loadBalancers, nil
}

func (provider *OutscaleOAPI) deleteLoadBalancers(ctx context.Context, loadBalancers []Object) {
	if len(loadBalancers) == 0 {
		return
	}
	for _, loadBalancer := range loadBalancers {
		log.Printf("Deleting load balancer %s... ", loadBalancer)
		deletionOpts := osc.DeleteLoadBalancerRequest{LoadBalancerName: loadBalancer}
		_, err := provider.client.DeleteLoadBalancer(ctx, deletionOpts)
		if err != nil {
			log.Printf("Error while deleting load balancer: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readNatServices(ctx context.Context) ([]Object, error) {
	natServices := make([]Object, 0)
	read, err := provider.client.ReadNatServices(
		ctx,
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

func (provider *OutscaleOAPI) deleteNatServices(ctx context.Context, natServices []Object) {
	if len(natServices) == 0 {
		return
	}
	for _, natService := range natServices {
		log.Printf("Deleting nat service %s... ", natService)
		deletionOpts := osc.DeleteNatServiceRequest{NatServiceId: natService}
		_, err := provider.client.DeleteNatService(ctx, deletionOpts)
		if err != nil {
			log.Printf("Error while deleting nat service: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readSecurityGroups(ctx context.Context) ([]Object, error) {
	securityGroups := make([]Object, 0)
	read, err := provider.client.ReadSecurityGroups(
		ctx,
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

func (provider *OutscaleOAPI) deleteSecurityGroupRules(ctx context.Context, securityGroupId string) error {
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
				sgMember.AccountId = new(string)
				sgMember.SecurityGroupName = new(string)
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

		_, err := provider.client.DeleteSecurityGroupRule(ctx, delete)
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
				sgMember.AccountId = new(string)
				sgMember.SecurityGroupName = new(string)
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

		_, err := provider.client.DeleteSecurityGroupRule(ctx, delete)
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

func (provider *OutscaleOAPI) deleteSecurityGroups(ctx context.Context, securityGroups []Object) {
	if len(securityGroups) == 0 {
		return
	}
	for _, sg := range securityGroups {
		if provider.deleteSecurityGroupRules(ctx, sg) != nil {
			continue
		}
		log.Printf("Deleting security group %s... ", sg)
		deletionOpts := osc.DeleteSecurityGroupRequest{SecurityGroupId: &sg}
		_, err := provider.client.DeleteSecurityGroup(ctx, deletionOpts)
		if err != nil {
			log.Printf("Error while deleting security groups: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readPublicIps(ctx context.Context) ([]Object, error) {
	publicIps := make([]Object, 0)
	read, err := provider.client.ReadPublicIps(
		ctx,
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

func (provider *OutscaleOAPI) unlinkPublicIp(ctx context.Context, publicIP *string) error {
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
	_, err := provider.client.UnlinkPublicIp(ctx, unlinkOpts)
	if err != nil {
		log.Printf("Error while unlinking public ip: %v\n", getErrorInfo(err))
		return err
	}
	log.Println("OK")
	return nil
}

func (provider *OutscaleOAPI) deletePublicIps(ctx context.Context, publicIps []Object) {
	if len(publicIps) == 0 {
		return
	}
	for _, publicIP := range publicIps {
		if provider.unlinkPublicIp(ctx, &publicIP) != nil {
			continue
		}
		log.Printf("Deleting public ip %s... ", publicIP)
		deletionOpts := osc.DeletePublicIpRequest{PublicIp: &publicIP}
		_, err := provider.client.DeletePublicIp(ctx, deletionOpts)
		if err != nil {
			log.Printf("Error while deleting public ip: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readVolumes(ctx context.Context) ([]Object, error) {
	volumes := make([]Object, 0)
	read, err := provider.client.ReadVolumes(ctx, osc.ReadVolumesRequest{
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

func (provider *OutscaleOAPI) deleteVolumes(ctx context.Context, volumes []Object) {
	if len(volumes) == 0 {
		return
	}
	for _, volume := range volumes {
		log.Printf("Deleting volume %s... ", volume)
		deletionOpts := osc.DeleteVolumeRequest{VolumeId: volume}
		_, err := provider.client.DeleteVolume(ctx, deletionOpts)
		if err != nil {
			log.Printf("Error while deleting volume: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readKeypairs(ctx context.Context) ([]Object, error) {
	keypairs := make([]Object, 0)
	read, err := provider.client.ReadKeypairs(ctx, osc.ReadKeypairsRequest{})
	if err != nil {
		return nil, fmt.Errorf("read key pairs: %w", getErrorInfo(err))
	}
	for _, keypair := range *read.Keypairs {
		keypairs = append(keypairs, *keypair.KeypairName)
	}
	return keypairs, nil
}

func (provider *OutscaleOAPI) deleteKeypairs(ctx context.Context, keypairs []Object) {
	if len(keypairs) == 0 {
		return
	}
	for _, keypair := range keypairs {
		log.Printf("Deleting keypair %s... ", keypair)
		deletionOpts := osc.DeleteKeypairRequest{KeypairName: &keypair}
		_, err := provider.client.DeleteKeypair(ctx, deletionOpts)
		if err != nil {
			log.Printf("Error while deleting keypair: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readRouteTables(ctx context.Context) ([]Object, error) {
	routeTables := make([]Object, 0)
	read, err := provider.client.ReadRouteTables(
		ctx,
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

func (provider *OutscaleOAPI) unlinkRouteTable(ctx context.Context, routeTableId string) error {
	routeTable := provider.cache.routeTables[routeTableId]
	if routeTable == nil || routeTable.LinkRouteTables == nil {
		return nil
	}
	for _, link := range routeTable.LinkRouteTables {
		if link.Main {
			continue
		}
		linkId := link.LinkRouteTableId
		log.Printf("Unlinking route table %s (link %s)... ", routeTableId, linkId)
		unlinkOps := osc.UnlinkRouteTableRequest{
			LinkRouteTableId: link.LinkRouteTableId,
		}
		_, err := provider.client.UnlinkRouteTable(ctx, unlinkOps)
		if err != nil {
			log.Printf(
				"Error while unlinking route table %s (links %s): %v\n",
				routeTableId,
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

func (provider *OutscaleOAPI) deleteRouteTables(ctx context.Context, routeTables []Object) {
	if len(routeTables) == 0 {
		return
	}
	for _, routeTable := range routeTables {
		if provider.unlinkRouteTable(ctx, routeTable) != nil {
			continue
		}
		log.Printf("Deleting route table %s... ", routeTable)
		deletionOpts := osc.DeleteRouteTableRequest{RouteTableId: routeTable}
		_, err := provider.client.DeleteRouteTable(ctx, deletionOpts)
		if err != nil {
			log.Printf("Error while deleting route table: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readInternetServices(ctx context.Context) ([]Object, error) {
	internetServices := make([]Object, 0)
	read, err := provider.client.ReadInternetServices(
		ctx,
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

func (provider *OutscaleOAPI) unlinkInternetSevice(ctx context.Context, internetServiceId string) error {
	internetService := provider.cache.internetServices[internetServiceId]
	if internetService == nil || internetService.NetId == "" {
		return nil
	}
	log.Printf("Unlinking internet service %s... ", internetServiceId)
	unlinkOps := osc.UnlinkInternetServiceRequest{
		InternetServiceId: internetServiceId,
		NetId:             internetService.NetId,
	}
	_, err := provider.client.UnlinkInternetService(ctx, unlinkOps)
	if err != nil {
		log.Printf("Error while unlinking internet service: %v\n", getErrorInfo(err))
		return err
	} else {
		log.Println("OK")
	}
	return nil
}

func (provider *OutscaleOAPI) deleteInternetServices(ctx context.Context, internetServices []Object) {
	if len(internetServices) == 0 {
		return
	}
	for _, internetService := range internetServices {
		if provider.unlinkInternetSevice(ctx, internetService) != nil {
			continue
		}
		log.Printf("Deleting internet service %s... ", internetService)
		deletionOpts := osc.DeleteInternetServiceRequest{InternetServiceId: internetService}
		_, err := provider.client.DeleteInternetService(ctx, deletionOpts)
		if err != nil {
			log.Printf("Error while deleting internet service: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readSubnets(ctx context.Context) ([]Object, error) {
	subnets := make([]Object, 0)
	read, err := provider.client.ReadSubnets(ctx, osc.ReadSubnetsRequest{})
	if err != nil {
		return nil, fmt.Errorf("read subnets: %w", getErrorInfo(err))
	}
	for _, subnet := range *read.Subnets {
		subnets = append(subnets, subnet.SubnetId)
	}
	return subnets, nil
}

func (provider *OutscaleOAPI) deleteSubnets(ctx context.Context, subnets []Object) {
	if len(subnets) == 0 {
		return
	}
	for _, subnet := range subnets {
		log.Printf("Deleting subnet %s... ", subnet)
		deletionOpts := osc.DeleteSubnetRequest{SubnetId: subnet}
		_, err := provider.client.DeleteSubnet(ctx, deletionOpts)
		if err != nil {
			log.Printf("Error while deleting subnet: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readNets(ctx context.Context) ([]Object, error) {
	nets := make([]Object, 0)
	read, err := provider.client.ReadNets(ctx, osc.ReadNetsRequest{
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

func (provider *OutscaleOAPI) deleteNets(ctx context.Context, nets []Object) {
	if len(nets) == 0 {
		return
	}
	for _, net := range nets {
		log.Printf("Deleting net %s... ", net)
		deletionOpts := osc.DeleteNetRequest{NetId: net}
		_, err := provider.client.DeleteNet(ctx, deletionOpts)
		if err != nil {
			log.Printf("Error while deleting net: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readAccountId(ctx context.Context) (*string, error) {
	if provider.cache.accountId == nil {
		read, err := provider.client.ReadAccounts(
			ctx,
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

func (provider *OutscaleOAPI) readImages(ctx context.Context) ([]Object, error) {
	images := make([]Object, 0)
	accountId, err := provider.readAccountId(ctx)
	if err != nil {
		return images, nil
	}
	var accountIds []string
	accountIds = append(accountIds, *accountId)
	read, err := provider.client.ReadImages(ctx, osc.ReadImagesRequest{
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

func (provider *OutscaleOAPI) deleteImages(ctx context.Context, images []Object) {
	if len(images) == 0 {
		return
	}
	for _, image := range images {
		log.Printf("Deleting image %s... ", image)
		deletionOpts := osc.DeleteImageRequest{ImageId: image}
		_, err := provider.client.DeleteImage(ctx, deletionOpts)
		if err != nil {
			log.Printf("Error while deleting image: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readSnapshots(ctx context.Context) ([]Object, error) {
	snapshots := make([]Object, 0)
	accountId, err := provider.readAccountId(ctx)
	if err != nil {
		return snapshots, nil
	}
	var accountIds []string
	accountIds = append(accountIds, *accountId)
	read, err := provider.client.ReadSnapshots(ctx, osc.ReadSnapshotsRequest{
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

func (provider *OutscaleOAPI) deleteSnapshots(ctx context.Context, snapshots []Object) {
	if len(snapshots) == 0 {
		return
	}
	for _, snapshot := range snapshots {
		log.Printf("Deleting snapshot %s... ", snapshot)
		deletionOpts := osc.DeleteSnapshotRequest{SnapshotId: snapshot}
		_, err := provider.client.DeleteSnapshot(ctx, deletionOpts)
		if err != nil {
			log.Printf("Error while deleting snapshot: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readVpnConnections(ctx context.Context) ([]Object, error) {
	vpnConnections := make([]Object, 0)
	read, err := provider.client.ReadVpnConnections(
		ctx,
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
		vpnConnections = append(vpnConnections, vpnConnection.VpnConnectionId)
	}
	return vpnConnections, nil
}

func (provider *OutscaleOAPI) deleteVpnConnections(ctx context.Context, vpnConnections []Object) {
	if len(vpnConnections) == 0 {
		return
	}
	for _, vpnConnection := range vpnConnections {
		log.Printf("Deleting vpn connection %s... ", vpnConnection)
		deletionOpts := osc.DeleteVpnConnectionRequest{VpnConnectionId: vpnConnection}
		_, err := provider.client.DeleteVpnConnection(ctx, deletionOpts)
		if err != nil {
			log.Printf("Error while deleting vpn connection: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readVirtualGateways(ctx context.Context) ([]Object, error) {
	virtualGateways := make([]Object, 0)
	read, err := provider.client.ReadVirtualGateways(
		ctx,
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
		virtualGateways = append(virtualGateways, virtualGateway.VirtualGatewayId)
	}
	return virtualGateways, nil
}

func (provider *OutscaleOAPI) deleteVirtualGateways(ctx context.Context, virtualGateways []Object) {
	if len(virtualGateways) == 0 {
		return
	}
	for _, virtualGateway := range virtualGateways {
		log.Printf("Deleting virtual gateway %s... ", virtualGateway)
		deletionOpts := osc.DeleteVirtualGatewayRequest{VirtualGatewayId: virtualGateway}
		_, err := provider.client.DeleteVirtualGateway(ctx, deletionOpts)
		if err != nil {
			log.Printf("Error while deleting virtual gateway: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readClientGateways(ctx context.Context) ([]Object, error) {
	clientGateways := make([]Object, 0)
	read, err := provider.client.ReadClientGateways(
		ctx,
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
		clientGateways = append(clientGateways, clientGateway.ClientGatewayId)
	}
	return clientGateways, nil
}

func (provider *OutscaleOAPI) deleteClientGateways(ctx context.Context, clientGateways []Object) {
	if len(clientGateways) == 0 {
		return
	}
	for _, clientGateway := range clientGateways {
		log.Printf("Deleting client gateway %s... ", clientGateway)
		deletionOpts := osc.DeleteClientGatewayRequest{ClientGatewayId: clientGateway}
		_, err := provider.client.DeleteClientGateway(ctx, deletionOpts)
		if err != nil {
			log.Printf("Error while deleting client gateway: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readNics(ctx context.Context) ([]Object, error) {
	nics := make([]Object, 0)
	read, err := provider.client.ReadNics(ctx, osc.ReadNicsRequest{})
	if err != nil {
		return nil, fmt.Errorf("read nics: %w", getErrorInfo(err))
	}
	for i, nic := range *read.Nics {
		nics = append(nics, nic.NicId)
		provider.cache.nics[nic.NicId] = &(*read.Nics)[i]
	}
	return nics, nil
}

func (provider *OutscaleOAPI) unlinkNics(ctx context.Context, nics []Object) {
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
		_, err := provider.client.UnlinkNic(ctx, unlinkOpts)
		if err != nil {
			log.Printf("Error while unlinking nic: %v\n", getErrorInfo(err))
			continue
		}
		log.Println("OK")
	}
}

func (provider *OutscaleOAPI) deleteNics(ctx context.Context, nics []Object) {
	if len(nics) == 0 {
		return
	}
	provider.unlinkNics(ctx, nics)
	for _, nicId := range nics {
		log.Printf("Deleting nic %s... ", nicId)
		deletionOpts := osc.DeleteNicRequest{NicId: nicId}
		_, err := provider.client.DeleteNic(ctx, deletionOpts)
		if err != nil {
			log.Printf("Error while deleting nic: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readAccessKeys(ctx context.Context) ([]Object, error) {
	accessKeys := make([]Object, 0)
	read, err := provider.client.ReadAccessKeys(
		ctx,
		osc.ReadAccessKeysRequest{},
	)
	if err != nil {
		return nil, fmt.Errorf("read ak: %w", getErrorInfo(err))
	}
	for _, accessKey := range *read.AccessKeys {
		if *accessKey.State == "ACTIVE" {
			// skip expired access keys that cannot be deleted
			if accessKey.ExpirationDate != nil && time.Now().After(accessKey.ExpirationDate.Time) {
				continue
			}
			accessKeys = append(accessKeys, *accessKey.AccessKeyId)
		}
	}
	return accessKeys, nil
}

func (provider *OutscaleOAPI) deleteAccessKeys(ctx context.Context, accessKeys []Object) {
	if len(accessKeys) == 0 {
		return
	}
	for _, accessKey := range accessKeys {
		log.Printf("Deleting access key %s... ", accessKey)
		deletionOpts := osc.DeleteAccessKeyRequest{AccessKeyId: accessKey}
		_, err := provider.client.DeleteAccessKey(ctx, deletionOpts)
		if err != nil {
			log.Printf("Error while deleting access key: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readNetAccessPoints(ctx context.Context) ([]Object, error) {
	netAccessPoints := make([]Object, 0)
	read, err := provider.client.ReadNetAccessPoints(
		ctx,
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
		netAccessPoints = append(netAccessPoints, netAccessPoint.NetAccessPointId)
	}
	return netAccessPoints, nil
}

func (provider *OutscaleOAPI) deleteNetAccessPoints(ctx context.Context, netAccessPoints []Object) {
	if len(netAccessPoints) == 0 {
		return
	}
	for _, netAccessPoint := range netAccessPoints {
		log.Printf("Deleting net access point %s... ", netAccessPoint)
		deletionOpts := osc.DeleteNetAccessPointRequest{NetAccessPointId: netAccessPoint}
		_, err := provider.client.DeleteNetAccessPoint(ctx, deletionOpts)
		if err != nil {
			log.Print("Error while deleting net access point: ")
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readNetPeerings(ctx context.Context) ([]Object, error) {
	netPeerings := make([]Object, 0)
	read, err := provider.client.ReadNetPeerings(
		ctx,
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

func (provider *OutscaleOAPI) deleteNetPeerings(ctx context.Context, netPeerings []Object) {
	if len(netPeerings) == 0 {
		return
	}
	for _, netPeering := range netPeerings {
		log.Printf("Deleting net peering %s... ", netPeering)
		deletionOpts := osc.DeleteNetPeeringRequest{NetPeeringId: netPeering}
		_, err := provider.client.DeleteNetPeering(ctx, deletionOpts)
		if err != nil {
			log.Print("Error while deleting net peering: %w", err)
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readUsers(ctx context.Context) ([]Object, error) {
	users := make([]Object, 0)
	read, err := provider.client.ReadUsers(ctx, osc.ReadUsersRequest{})
	if err != nil {
		return nil, fmt.Errorf("read users: %w", getErrorInfo(err))
	}
	for _, user := range *read.Users {
		users = append(users, *user.UserName)
	}
	return users, nil
}

func (provider *OutscaleOAPI) deleteUsers(ctx context.Context, users []Object) {
	if len(users) == 0 {
		return
	}
	for _, user := range users {
		log.Printf("Deleting user %s... ", user)
		deleteOpts := osc.DeleteUserRequest{UserName: user}
		_, err := provider.client.DeleteUser(ctx, deleteOpts)
		if err != nil {
			log.Print("Error while deleting user: %w", err)
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readUserAccessKeys(ctx context.Context) ([]Object, error) {
	accessKeys := make([]Object, 0)

	readUser, err := provider.client.ReadUsers(ctx, osc.ReadUsersRequest{})
	if err != nil {
		return nil, fmt.Errorf("read users: %w", getErrorInfo(err))
	}

	for _, user := range *readUser.Users {
		read, err := provider.client.ReadAccessKeys(
			ctx,
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

func (provider *OutscaleOAPI) deleteUserAccessKeys(ctx context.Context, accessKeys []Object) {
	if len(accessKeys) == 0 {
		return
	}
	for _, accessKey := range accessKeys {
		log.Printf("Deleting user access key %s... ", accessKey)
		parts := strings.SplitN(string(accessKey), ",", 2)
		if len(parts) != 2 {
			log.Printf("Invalid access key format: %s", accessKey)
			continue
		}

		deletionOpts := osc.DeleteAccessKeyRequest{AccessKeyId: parts[1], UserName: &parts[0]}
		_, err := provider.client.DeleteAccessKey(ctx, deletionOpts)
		if err != nil {
			log.Printf("Error while deleting user access key: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readPolicies(ctx context.Context) ([]Object, error) {
	policies := make([]Object, 0)

	read, err := provider.client.ReadPolicies(ctx, osc.ReadPoliciesRequest{})
	if err != nil {
		return nil, fmt.Errorf("read policies: %w", getErrorInfo(err))
	}
	for _, policy := range *read.Policies {
		policies = append(policies, *policy.Orn)
	}
	return policies, nil
}

func (provider *OutscaleOAPI) deletePolicies(ctx context.Context, policies []Object) {
	if len(policies) == 0 {
		return
	}
	for _, policy := range policies {
		log.Printf("Deleting policy %s... ", policy)
		deleteOpts := osc.DeletePolicyRequest{PolicyOrn: policy}
		_, err := provider.client.DeletePolicy(ctx, deleteOpts)
		if err != nil {
			log.Print("Error while deleting policy: %w", err)
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readPolicyLinks(ctx context.Context) ([]Object, error) {
	policyLinks := make([]Object, 0)

	read, err := provider.client.ReadPolicies(ctx, osc.ReadPoliciesRequest{})
	if err != nil {
		return nil, fmt.Errorf("read policies: %w", getErrorInfo(err))
	}
	for _, policy := range *read.Policies {
		read, err := provider.client.ReadEntitiesLinkedToPolicy(
			ctx,
			osc.ReadEntitiesLinkedToPolicyRequest{
				EntitiesType: &[]osc.ReadEntitiesLinkedToPolicyRequestEntitiesType{"USER", "GROUP"},
				PolicyOrn:    *policy.Orn,
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

func (provider *OutscaleOAPI) deletePolicyLinks(ctx context.Context, policyLinks []Object) {
	if len(policyLinks) == 0 {
		return
	}

	for _, policylink := range policyLinks {
		log.Printf("Deleting policy link %s... ", policylink)
		parts := strings.SplitN(string(policylink), ",", 3)
		if len(parts) != 3 {
			log.Printf("Invalid policy link format: %s", policylink)
			continue
		}
		linkType := parts[0]
		policyOrn := parts[1]
		linkName := parts[2]

		switch linkType {
		case "USER":
			deleteOpts := osc.UnlinkPolicyRequest{
				PolicyOrn: policyOrn,
				UserName:  linkName,
			}
			_, err := provider.client.UnlinkPolicy(ctx, deleteOpts)
			if err != nil {
				log.Print("Error while unlinking policy: %w", err)
			}

		case "GROUP":
			deleteOpts := osc.UnlinkManagedPolicyFromUserGroupRequest{
				PolicyOrn:     policyOrn,
				UserGroupName: linkName,
			}
			_, err := provider.client.UnlinkManagedPolicyFromUserGroup(
				ctx,
				deleteOpts,
			)
			if err != nil {
				log.Print("Error while unlinking policy: %w", err)
			}
		}
	}
}

func (provider *OutscaleOAPI) readPolicyVersions(ctx context.Context) ([]Object, error) {
	policyVersions := make([]Object, 0)

	read, err := provider.client.ReadPolicies(ctx, osc.ReadPoliciesRequest{})
	if err != nil {
		return nil, fmt.Errorf("read policies: %w", getErrorInfo(err))
	}
	for _, policy := range *read.Policies {
		read, err := provider.client.ReadPolicyVersions(
			ctx,
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

func (provider *OutscaleOAPI) deletePolicyVersions(ctx context.Context, policyVersions []Object) {
	if len(policyVersions) == 0 {
		return
	}

	for _, policyVersion := range policyVersions {
		log.Printf("Deleting policy version %s... ", policyVersion)
		parts := strings.SplitN(string(policyVersion), ",", 2)
		if len(parts) != 2 {
			log.Printf("Invalid policy version format: %s", policyVersion)
			continue
		}

		deleteOpts := osc.DeletePolicyVersionRequest{
			PolicyOrn: parts[0],
			VersionId: parts[1],
		}
		_, err := provider.client.DeletePolicyVersion(ctx, deleteOpts)
		if err != nil {
			log.Print("Error while deleting policy version: %w", err)
		}
	}
}

func (provider *OutscaleOAPI) readFlexibleGpus(ctx context.Context) ([]Object, error) {
	flexibleGpus := make([]Object, 0)

	read, err := provider.client.ReadFlexibleGpus(ctx, osc.ReadFlexibleGpusRequest{})
	if err != nil {
		return nil, fmt.Errorf("read flexible gpus: %w", getErrorInfo(err))
	}
	for i, gpu := range *read.FlexibleGpus {
		flexibleGpus = append(flexibleGpus, *gpu.FlexibleGpuId)
		provider.cache.flexibleGpus[*gpu.FlexibleGpuId] = &(*read.FlexibleGpus)[i]
	}
	return flexibleGpus, nil
}

func (provider *OutscaleOAPI) unlinkFlexibleGpus(ctx context.Context, flexibleGpus []Object) {
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
		_, err := provider.client.UnlinkFlexibleGpu(ctx, unlinkOpts)
		if err != nil {
			log.Printf("Error while unlinking flexible gpu: %v\n", getErrorInfo(err))
			continue
		}
		log.Println("OK")
	}
}

func (provider *OutscaleOAPI) deleteFlexibleGpus(ctx context.Context, flexibleGpus []Object) {
	if len(flexibleGpus) == 0 {
		return
	}
	provider.unlinkFlexibleGpus(ctx, flexibleGpus)
	for _, gpu := range flexibleGpus {
		log.Printf("Releasing flexible gpu %s... ", gpu)
		deleteOpts := osc.DeleteFlexibleGpuRequest{FlexibleGpuId: gpu}
		_, err := provider.client.DeleteFlexibleGpu(ctx, deleteOpts)
		if err != nil {
			log.Print("Error while deleting flexible gpu: %w", err)
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readCas(ctx context.Context) ([]Object, error) {
	cas := make([]Object, 0)
	read, err := provider.client.ReadCas(ctx, osc.ReadCasRequest{})
	if err != nil {
		return nil, fmt.Errorf("read cas: %w", getErrorInfo(err))
	}
	for _, ca := range *read.Cas {
		cas = append(cas, *ca.CaId)
	}
	return cas, nil
}

func (provider *OutscaleOAPI) deleteCas(ctx context.Context, cas []Object) {
	if len(cas) == 0 {
		return
	}
	for _, ca := range cas {
		log.Printf("Deleting CA %s... ", ca)
		deleteOpts := osc.DeleteCaRequest{CaId: ca}
		_, err := provider.client.DeleteCa(ctx, deleteOpts)
		if err != nil {
			log.Printf("Error while deleting CA: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readServerCertificates(ctx context.Context) ([]Object, error) {
	serverCertificates := make([]Object, 0)
	read, err := provider.client.ReadServerCertificates(ctx, osc.ReadServerCertificatesRequest{})
	if err != nil {
		return nil, fmt.Errorf("read server certificates: %w", getErrorInfo(err))
	}
	for _, cert := range *read.ServerCertificates {
		serverCertificates = append(serverCertificates, *cert.Name)
	}
	return serverCertificates, nil
}

func (provider *OutscaleOAPI) deleteServerCertificates(ctx context.Context, serverCertificates []Object) {
	if len(serverCertificates) == 0 {
		return
	}
	for _, cert := range serverCertificates {
		log.Printf("Deleting server certificate %s... ", cert)
		deleteOpts := osc.DeleteServerCertificateRequest{Name: cert}
		_, err := provider.client.DeleteServerCertificate(ctx, deleteOpts)
		if err != nil {
			log.Printf("Error while deleting server certificate: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOAPI) readDhcpOptions(ctx context.Context) ([]Object, error) {
	dhcpOptions := make([]Object, 0)
	read, err := provider.client.ReadDhcpOptions(ctx, osc.ReadDhcpOptionsRequest{})
	if err != nil {
		return nil, fmt.Errorf("read dhcp options: %w", getErrorInfo(err))
	}
	for _, option := range *read.DhcpOptionsSets {
		dhcpOptions = append(dhcpOptions, *option.DhcpOptionsSetId)
	}
	return dhcpOptions, nil
}

func (provider *OutscaleOAPI) deleteDhcpOptions(ctx context.Context, dhcpOptions []Object) {
	if len(dhcpOptions) == 0 {
		return
	}
	for _, option := range dhcpOptions {
		log.Printf("Deleting DHCP option %s... ", option)
		deleteOpts := osc.DeleteDhcpOptionsRequest{DhcpOptionsSetId: option}
		_, err := provider.client.DeleteDhcpOptions(ctx, deleteOpts)
		if err != nil {
			log.Printf("Error while deleting DHCP option: %v\n", getErrorInfo(err))
		} else {
			log.Println("OK")
		}
	}
}
