/*
Copyright 2023 The Invisinets Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package azure_plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	invisinetspb "github.com/NetSys/invisinets/pkg/invisinetspb"
	utils "github.com/NetSys/invisinets/pkg/utils"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const maxPriority = 4096

const (
	invisinetsPrefix = "invisinets"
)

type ResourceIDInfo struct {
	SubscriptionID    string
	ResourceGroupName string
	ResourceName      string
}

type azurePluginServer struct {
	invisinetspb.UnimplementedCloudPluginServer
	azureHandler           AzureSDKHandler
	orchestratorServerAddr string
}

const (
	vpnLocation                = "westus" // TODO @seankimkdy: should this be configurable/dynamic?
	gatewaySubnetName          = "GatewaySubnet"
	gatewaySubnetAddressPrefix = "192.168.255.0/27"
)

func (s *azurePluginServer) setupAzureHandler(resourceIdInfo ResourceIDInfo) error {
	cred, err := s.azureHandler.GetAzureCredentials()
	if err != nil {
		utils.Log.Printf("An error occured while getting azure credentials:%+v", err)
		return err
	}
	s.azureHandler.SetSubIdAndResourceGroup(resourceIdInfo.SubscriptionID, resourceIdInfo.ResourceGroupName)
	err = s.azureHandler.InitializeClients(cred)
	if err != nil {
		utils.Log.Printf("An error occured while initializing azure clients: %+v", err)
		return err
	}

	return nil
}

// GetPermitList returns the permit list for the given resource by getting the NSG rules
// associated with the resource and filtering out the Invisinets rules
func (s *azurePluginServer) GetPermitList(ctx context.Context, req *invisinetspb.GetPermitListRequest) (*invisinetspb.GetPermitListResponse, error) {
	resourceId := req.Resource
	resourceIdInfo, err := getResourceIDInfo(resourceId)
	if err != nil {
		utils.Log.Printf("An error occured while getting resource ID info: %+v", err)
		return nil, err
	}
	err = s.setupAzureHandler(resourceIdInfo)
	if err != nil {
		return nil, err
	}

	// make sure the resource is in the right namespace
	err = s.getAndCheckResourceNamespace(ctx, resourceId, req.Namespace)
	if err != nil {
		return nil, err
	}

	// get the nsg associated with the resource
	nsg, err := s.getNSGFromResource(ctx, resourceId)
	if err != nil {
		utils.Log.Printf("An error occured while getting NSG for resource %s: %+v", resourceId, err)
		return nil, err
	}

	// initialize a list of permit list rules
	rules := []*invisinetspb.PermitListRule{}

	// get the NSG rules
	for _, rule := range nsg.Properties.SecurityRules {
		if !strings.HasPrefix(*rule.Name, denyAllNsgRulePrefix) && strings.HasPrefix(*rule.Name, invisinetsPrefix) {
			plRule, err := s.azureHandler.GetPermitListRuleFromNSGRule(rule)
			if err != nil {
				utils.Log.Printf("An error occured while getting Invisinets rule from NSG rule: %+v", err)
				return nil, err
			}
			plRule.Name = getRuleNameFromNSGRuleName(plRule.Name)
			rules = append(rules, plRule)
		}
	}
	return &invisinetspb.GetPermitListResponse{Rules: rules}, nil
}

// AddPermitListRules does the mapping from Invisinets to Azure by creating/updating NSG for the given resource.
// It creates an NSG rule for each permit list rule and applies this NSG to the associated resource (VM)'s NIC (if it doesn't exist).
// It returns a BasicResponse that includes the nsg ID if successful and an error if it fails.
func (s *azurePluginServer) AddPermitListRules(ctx context.Context, req *invisinetspb.AddPermitListRulesRequest) (*invisinetspb.AddPermitListRulesResponse, error) {
	resourceID := req.GetResource()
	resourceIdInfo, err := getResourceIDInfo(resourceID)
	if err != nil {
		utils.Log.Printf("An error occured while getting resource ID info: %+v", err)
		return nil, err
	}
	err = s.setupAzureHandler(resourceIdInfo)
	if err != nil {
		return nil, err
	}

	// make sure the resource is in the right namespace
	err = s.getAndCheckResourceNamespace(ctx, resourceID, req.Namespace)
	if err != nil {
		return nil, err
	}

	// get the nic associated with the resource
	nic, err := s.azureHandler.GetResourceNIC(ctx, resourceID)
	if err != nil {
		utils.Log.Printf("An error occured while getting NIC for resource %s: %+v", resourceID, err)
		return nil, err
	}

	// get the NSG associated with the resource
	nsg, err := s.getNSG(ctx, nic, resourceID)

	if err != nil {
		utils.Log.Printf("An error occured while getting NSG for resource %s: %+v", resourceID, err)
		return nil, err
	}

	var existingRulePriorities map[string]int32 = make(map[string]int32)
	var reservedPrioritiesInbound map[int32]bool = make(map[int32]bool)
	var reservedPrioritiesOutbound map[int32]bool = make(map[int32]bool)
	err = s.setupMaps(reservedPrioritiesInbound, reservedPrioritiesOutbound, existingRulePriorities, nsg)
	if err != nil {
		utils.Log.Printf("An error occured during setup: %+v", err)
		return nil, err
	}
	var outboundPriority int32 = 100
	var inboundPriority int32 = 100

	resourceAddress := *nic.Properties.IPConfigurations[0].Properties.PrivateIPAddress

	// Get used address spaces of all clouds
	controllerConn, err := grpc.Dial(s.orchestratorServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("unable to establish connection with orchestrator: %w", err)
	}
	defer controllerConn.Close()
	controllerClient := invisinetspb.NewControllerClient(controllerConn)
	getUsedAddressSpacesResp, err := controllerClient.GetUsedAddressSpaces(context.Background(), &invisinetspb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("unable to get used address spaces: %w", err)
	}

	// get the vnet to be able to get both the address space as well as the peering when needed
	resourceVnet, err := s.azureHandler.GetVNet(ctx, getVnetName(*nic.Location, req.Namespace))
	if err != nil {
		utils.Log.Printf("An error occured while getting resource vnet:%+v", err)
		return nil, err
	}

	// Get subnet address space
	var subnetAddressPrefix string
	for _, subnet := range resourceVnet.Properties.Subnets {
		if *subnet.Name == "default" {
			if subnet.Properties.AddressPrefix != nil {
				// Check to avoid nil pointer dereference
				subnetAddressPrefix = *subnet.Properties.AddressPrefix
			}
			break
		}
	}
	if subnetAddressPrefix == "" {
		return nil, fmt.Errorf("unable to get subnet address prefix")
	}

	invisinetsVnetsMap, err := s.azureHandler.GetVNetsAddressSpaces(ctx, getInvisinetsNamespacePrefix(req.Namespace))

	if err != nil {
		utils.Log.Printf("An error occured while getting invisinets vnets address spaces:%+v", err)
		return nil, err
	}

	// Add the rules to the NSG
	for _, rule := range req.GetRules() {
		err = utils.CheckAndConnectClouds(utils.AZURE, subnetAddressPrefix, req.Namespace, ctx, rule, getUsedAddressSpacesResp.AddressSpaceMappings, controllerClient)
		if err != nil {
			return nil, fmt.Errorf("unable to check and connect clouds: %w", err)
		}

		// TODO @seankimkdy: merge this process with the checking address spaces across all clouds to avoid duplicate checking of Azure address spaces
		err := s.checkAndCreatePeering(ctx, resourceVnet, rule, invisinetsVnetsMap, req.Namespace)
		if err != nil {
			utils.Log.Printf("An error occured while checking network peering:%+v", err)
			return nil, err
		}

		// To avoid conflicted priorities, we need to check whether the priority is already used by other rules
		// if the priority is already used, we need to find the next available priority
		priority, ok := existingRulePriorities[getNSGRuleName(rule.Name)]
		if !ok {
			if rule.Direction == invisinetspb.Direction_INBOUND {
				priority = getPriority(reservedPrioritiesInbound, inboundPriority, maxPriority)
				inboundPriority = priority + 1
			} else if rule.Direction == invisinetspb.Direction_OUTBOUND {
				priority = getPriority(reservedPrioritiesOutbound, outboundPriority, maxPriority)
				outboundPriority = priority + 1
			}
		}

		// Create the NSG rule
		securityRule, err := s.azureHandler.CreateSecurityRule(ctx, rule, *nsg.Name, getNSGRuleName(rule.Name), resourceAddress, priority)
		if err != nil {
			utils.Log.Printf("An error occured while creating security rule:%+v", err)
			return nil, err
		}
		utils.Log.Printf("Successfully created network security rule: %s", *securityRule.ID)
	}

	return &invisinetspb.AddPermitListRulesResponse{}, nil
}

// DeletePermitListRules does the mapping from Invisinets to Azure by deleting NSG rules for the given resource.
func (s *azurePluginServer) DeletePermitListRules(c context.Context, req *invisinetspb.DeletePermitListRulesRequest) (*invisinetspb.DeletePermitListRulesResponse, error) {
	resourceID := req.GetResource()
	resourceIdInfo, err := getResourceIDInfo(resourceID)
	if err != nil {
		utils.Log.Printf("An error occured while getting resource ID info: %+v", err)
		return nil, err
	}
	err = s.setupAzureHandler(resourceIdInfo)
	if err != nil {
		return nil, err
	}

	// make sure the resource is in the right namespace
	err = s.getAndCheckResourceNamespace(c, resourceID, req.Namespace)
	if err != nil {
		return nil, err
	}

	nsg, err := s.getNSGFromResource(c, resourceID)
	if err != nil {
		utils.Log.Printf("An error occured while getting NSG for resource %s: %+v", resourceID, err)
		return nil, err
	}

	for _, rule := range req.GetRuleNames() {
		err := s.azureHandler.DeleteSecurityRule(c, *nsg.Name, getNSGRuleName(rule))
		if err != nil {
			utils.Log.Printf("An error occured while deleting security rule:%+v", err)
			return nil, err
		}
		utils.Log.Printf("Successfully deleted network security rule: %s", rule)
	}

	return &invisinetspb.DeletePermitListRulesResponse{}, nil
}

// CreateResource does the mapping from Invisinets to Azure to create an invisinets enabled resource
// which means the resource should be added to a valid invisinets network, the attachement to an invisinets network
// is determined by the resource's location.
func (s *azurePluginServer) CreateResource(ctx context.Context, resourceDesc *invisinetspb.ResourceDescription) (*invisinetspb.CreateResourceResponse, error) {
	invisinetsVm, err := getVmFromResourceDesc(resourceDesc.Description)
	if err != nil {
		utils.Log.Printf("Resource description is invalid:%+v", err)
		return nil, err
	}

	resourceIdInfo, err := getResourceIDInfo(resourceDesc.Id)
	if err != nil {
		utils.Log.Printf("An error occured while getting resource id info:%+v", err)
		return nil, err
	}

	err = s.setupAzureHandler(resourceIdInfo)
	if err != nil {
		return nil, err
	}

	vnetName := getVnetName(*invisinetsVm.Location, resourceDesc.Namespace)
	invisinetsVnet, err := s.azureHandler.GetInvisinetsVnet(ctx, vnetName, *invisinetsVm.Location, resourceDesc.Namespace, s.orchestratorServerAddr)
	if err != nil {
		utils.Log.Printf("An error occured while getting invisinets vnet:%+v", err)
		return nil, err
	}

	nic, err := s.azureHandler.CreateNetworkInterface(ctx, *invisinetsVnet.Properties.Subnets[0].ID, *invisinetsVm.Location, getInvisinetsResourceName("nic"))
	if err != nil {
		utils.Log.Printf("An error occured while creating network interface:%+v", err)
		return nil, err
	}

	invisinetsVm.Properties.NetworkProfile = &armcompute.NetworkProfile{
		NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
			{
				ID: nic.ID,
			},
		},
	}

	invisinetsVm, err = s.azureHandler.CreateVirtualMachine(ctx, *invisinetsVm, resourceIdInfo.ResourceName)
	if err != nil {
		utils.Log.Printf("An error occured while creating the virtual machine:%+v", err)
		return nil, err
	}

	nic, err = s.azureHandler.GetResourceNIC(ctx, *invisinetsVm.ID)
	if err != nil {
		utils.Log.Printf("An error occured while getting the network interface:%+v", err)
		return nil, err
	}

	// Create VPN gateway vnet if not already created
	// The vnet is created even if there's no multicloud connections at the moment for ease of connection in the future.
	// Note that vnets are free, so this is not a problem.
	vpnGwVnetName := getVpnGatewayVnetName(resourceDesc.Namespace)
	_, err = s.azureHandler.GetVirtualNetwork(ctx, vpnGwVnetName)
	if err != nil {
		if isErrorNotFound(err) {
			virtualNetworkParameters := armnetwork.VirtualNetwork{
				Location: to.Ptr(vpnLocation),
				Properties: &armnetwork.VirtualNetworkPropertiesFormat{
					AddressSpace: &armnetwork.AddressSpace{
						AddressPrefixes: []*string{to.Ptr(gatewaySubnetAddressPrefix)},
					},
					Subnets: []*armnetwork.Subnet{
						{
							Name: to.Ptr(gatewaySubnetName),
							Properties: &armnetwork.SubnetPropertiesFormat{
								AddressPrefix: to.Ptr(gatewaySubnetAddressPrefix),
							},
						},
					},
				},
			}
			_, err = s.azureHandler.CreateVirtualNetwork(ctx, getVpnGatewayVnetName(resourceDesc.Namespace), virtualNetworkParameters)
			if err != nil {
				return nil, fmt.Errorf("unable to create VPN gateway vnet: %w", err)
			}
		} else {
			return nil, fmt.Errorf("unable to get VPN gateway vnet: %w", err)
		}
	}

	// Create peering VPN gateway vnet and VM vnet. If the VPN gateway already exists, then establish a VPN gateway transit relationship where the vnet can use the gatewayVnet's VPN gateway.
	// - This peering is created even if there's no multicloud connections at the moment for ease of connection in the future.
	// - Peerings are only charge based on amount of data transferred, so this will not incur extra charge until the VPN gateway is created.
	// - VPN gateway transit relationship cannot be established before the VPN gateway creation.
	// - If the VPN gateway hasn't been created, then the gateway transit relationship will be established on VPN gateway creation.
	_, err = s.azureHandler.GetVirtualNetworkPeering(ctx, vnetName, vpnGwVnetName)
	var peeringExists bool
	if err != nil {
		if isErrorNotFound(err) {
			peeringExists = false
		} else {
			return nil, fmt.Errorf("unable to get vnet peering between VM vnet and VPN gateway vnet: %w", err)
		}
	} else {
		peeringExists = true
	}
	// Only add peering if it doesn't exist
	if !peeringExists {
		vpnGwName := getVpnGatewayName(resourceDesc.Namespace)
		_, err = s.azureHandler.GetVirtualNetworkGateway(ctx, vpnGwName)
		if err != nil {
			if isErrorNotFound(err) {
				// Create regular peering which will be augmented with gateway transit relationship later on VPN gateway creation
				err = s.azureHandler.CreateVnetPeering(ctx, vnetName, vpnGwVnetName)
				if err != nil {
					return nil, fmt.Errorf("unable to create vnet peerings between VM vnet and VPN gateway vnet: %w", err)
				}
			} else {
				return nil, fmt.Errorf("unable to get VPN gateway: %w", err)
			}
		} else {
			// Create peering with gateway transit relationship if VPN gateway already exists
			err = s.azureHandler.CreateOrUpdateVnetPeeringRemoteGateway(ctx, vnetName, vpnGwVnetName, nil, nil)
			if err != nil {
				return nil, fmt.Errorf("unable to create vnet peerings (with gateway transit) between VM vnet and VPN gateway vnet: %w", err)
			}
		}
	}

	return &invisinetspb.CreateResourceResponse{Name: *invisinetsVm.Name, Uri: *invisinetsVm.ID, Ip: *nic.Properties.IPConfigurations[0].Properties.PrivateIPAddress}, nil
}

// GetUsedAddressSpaces returns the address spaces used by invisinets which are the address spaces of the invisinets vnets
func (s *azurePluginServer) GetUsedAddressSpaces(ctx context.Context, req *invisinetspb.GetUsedAddressSpacesRequest) (*invisinetspb.GetUsedAddressSpacesResponse, error) {
	resp := &invisinetspb.GetUsedAddressSpacesResponse{}
	resp.AddressSpaceMappings = make([]*invisinetspb.AddressSpaceMapping, len(req.Deployments))
	for i, deployment := range req.Deployments {
		resp.AddressSpaceMappings[i] = &invisinetspb.AddressSpaceMapping{
			Cloud:     utils.AZURE,
			Namespace: deployment.Namespace,
		}
		resourceIdInfo, err := getResourceIDInfo(deployment.Id)
		if err != nil {
			utils.Log.Printf("An error occured while getting resource ID info: %+v", err)
			return nil, err
		}
		err = s.setupAzureHandler(resourceIdInfo)
		if err != nil {
			return nil, err
		}

		addressSpaces, err := s.azureHandler.GetVNetsAddressSpaces(ctx, getInvisinetsNamespacePrefix(deployment.Namespace))
		if err != nil {
			utils.Log.Printf("An error occured while getting address spaces:%+v", err)
			return nil, err
		}
		for _, address := range addressSpaces {
			resp.AddressSpaceMappings[i].AddressSpaces = append(resp.AddressSpaceMappings[i].AddressSpaces, address)
		}
	}
	return resp, nil
}

func (s *azurePluginServer) GetUsedAsns(ctx context.Context, req *invisinetspb.GetUsedAsnsRequest) (*invisinetspb.GetUsedAsnsResponse, error) {
	resp := &invisinetspb.GetUsedAsnsResponse{}
	for _, deployment := range req.Deployments {
		resourceIdInfo, err := getResourceIDInfo(deployment.Id)
		if err != nil {
			utils.Log.Printf("An error occured while getting resource ID info: %+v", err)
			return nil, err
		}
		err = s.setupAzureHandler(resourceIdInfo)
		if err != nil {
			return nil, err
		}

		virtualNetworkGatewayName := getVpnGatewayName(deployment.Namespace)
		virtualNetworkGateway, err := s.azureHandler.GetVirtualNetworkGateway(ctx, virtualNetworkGatewayName)
		if err != nil {
			if isErrorNotFound(err) {
				continue
			} else {
				return nil, fmt.Errorf("unable to get virtual network gateway: %w", err)
			}
		}
		resp.Asns = append(resp.Asns, uint32(*virtualNetworkGateway.Properties.BgpSettings.Asn))
	}
	return resp, nil
}

func (s *azurePluginServer) GetUsedBgpPeeringIpAddresses(ctx context.Context, req *invisinetspb.GetUsedBgpPeeringIpAddressesRequest) (*invisinetspb.GetUsedBgpPeeringIpAddressesResponse, error) {
	resp := &invisinetspb.GetUsedBgpPeeringIpAddressesResponse{}
	for _, deployment := range req.Deployments {
		resourceIdInfo, err := getResourceIDInfo(deployment.Id)
		if err != nil {
			utils.Log.Printf("An error occured while getting resource ID info: %+v", err)
			return nil, err
		}
		err = s.setupAzureHandler(resourceIdInfo)
		if err != nil {
			return nil, err
		}

		virtualNetworkGatewayName := getVpnGatewayName(deployment.Namespace)
		virtualNetworkGateway, err := s.azureHandler.GetVirtualNetworkGateway(ctx, virtualNetworkGatewayName)
		if err != nil {
			if isErrorNotFound(err) {
				continue
			} else {
				return nil, fmt.Errorf("unable to get virtual network gateway: %w", err)
			}
		}
		for _, bgpPeeringAddress := range virtualNetworkGateway.Properties.BgpSettings.BgpPeeringAddresses {
			resp.IpAddresses = append(resp.IpAddresses, *bgpPeeringAddress.CustomBgpIPAddresses[0])
		}
	}
	return resp, nil
}

// getNSG returns the network security group object given the resource NIC
func (s *azurePluginServer) getNSG(ctx context.Context, nic *armnetwork.Interface, resourceID string) (*armnetwork.SecurityGroup, error) {
	var nsg *armnetwork.SecurityGroup
	if nic.Properties.NetworkSecurityGroup != nil {
		nsg = nic.Properties.NetworkSecurityGroup

		// nic.Properties.NetworkSecurityGroup returns an nsg obj with only the ID and other fields are nil
		// so this way we need to get the nsg object from the ID using nsgClient
		nsgID := *nsg.ID
		nsgName, err := s.azureHandler.GetLastSegment(nsgID)
		if err != nil {
			utils.Log.Printf("An error occured while getting NSG name for resource %s: %+v", resourceID, err)
			return nil, err
		}

		nsg, err = s.azureHandler.GetSecurityGroup(ctx, nsgName)
		if err != nil {
			utils.Log.Printf("An error occured while getting NSG for resource %s: %+v", resourceID, err)
			return nil, err
		}
	} else {
		// TODO @nnomier: should we handle this in another way?
		return nil, fmt.Errorf("resource %s does not have a default network security group", resourceID)
	}
	return nsg, nil
}

// getNSGFromResource gets the NSG associated with the given resource
// by getting the NIC associated with the resource and then getting the NSG associated with the NIC
func (s *azurePluginServer) getNSGFromResource(c context.Context, resourceID string) (*armnetwork.SecurityGroup, error) {
	// get the nic associated with the resource
	nic, err := s.azureHandler.GetResourceNIC(c, resourceID)
	if err != nil {
		utils.Log.Printf("An error occured while getting NIC for resource %s: %+v", resourceID, err)
		return nil, err
	}

	// avoid nil pointer dereference error
	if nic.Properties.NetworkSecurityGroup == nil {
		return nil, fmt.Errorf("resource %s does not have a network security group", resourceID)
	}

	nsgID := *nic.Properties.NetworkSecurityGroup.ID
	nsgName, err := s.azureHandler.GetLastSegment(nsgID)
	if err != nil {
		utils.Log.Printf("An error occured while getting NSG name for resource %s: %+v", resourceID, err)
		return nil, err
	}

	nsg, err := s.azureHandler.GetSecurityGroup(c, nsgName)
	if err != nil {
		utils.Log.Printf("An error occured while getting NSG for resource %s: %+v", resourceID, err)
		return nil, err
	}

	return nsg, nil
}

// Extract the Vnet name from the subnet ID
func getVnetFromSubnetId(subnetId string) string {
	parts := strings.Split(subnetId, "/")
	return parts[8] // TODO @smcclure20: do this in a less brittle way
}

// Check if the resource is in a vnet for the given namespace
func (s *azurePluginServer) getAndCheckResourceNamespace(c context.Context, resourceID string, namespace string) error {
	if namespace == "" {
		return fmt.Errorf("namespace cannot be empty")
	}

	// get the vnet associated with the resource
	nic, err := s.azureHandler.GetResourceNIC(c, resourceID)
	if err != nil {
		utils.Log.Printf("An error occured while getting nic for resource %s: %+v", resourceID, err)
		return err
	}
	vnet := getVnetFromSubnetId(*nic.Properties.IPConfigurations[0].Properties.Subnet.ID)

	if !strings.HasPrefix(vnet, getInvisinetsNamespacePrefix(namespace)) {
		return fmt.Errorf("resource %s is not in the namespace %s", resourceID, namespace)
	}

	return nil
}

// setupMaps fills the reservedPrioritiesInbound and reservedPrioritiesOutbound maps with the priorities of the existing rules in the NSG
// This is done to avoid priorities conflicts when creating new rules
// Existing rules map is filled to ensure that rules that just need their contents updated do not get recreated with new priorities
func (s *azurePluginServer) setupMaps(reservedPrioritiesInbound map[int32]bool, reservedPrioritiesOutbound map[int32]bool, existingRulePriorities map[string]int32, nsg *armnetwork.SecurityGroup) error {
	for _, rule := range nsg.Properties.SecurityRules {
		if *rule.Properties.Direction == armnetwork.SecurityRuleDirectionInbound {
			reservedPrioritiesInbound[*rule.Properties.Priority] = true
		} else if *rule.Properties.Direction == armnetwork.SecurityRuleDirectionOutbound {
			reservedPrioritiesOutbound[*rule.Properties.Priority] = true
		}

		// skip rules that are not created by Invisinets, because some rules are added by default and have
		// different fields such as port ranges which is not supported by Invisinets at the moment
		if !strings.HasPrefix(*rule.Name, invisinetsPrefix) {
			continue
		}
		existingRulePriorities[*rule.Name] = *rule.Properties.Priority
	}
	return nil
}

// getPriority returns the next available priority that is not used by other rules
func getPriority(reservedPriorities map[int32]bool, start int32, end int32) int32 {
	var i int32
	for i = start; i < end; i++ {
		if !reservedPriorities[i] {
			reservedPriorities[i] = true
			break
		}
	}
	return i
}

// getVmFromResourceDesc gets the armcompute.VirtualMachine object
// from the given resource description which should be a valid resource payload for a VM
func getVmFromResourceDesc(resourceDesc []byte) (*armcompute.VirtualMachine, error) {
	vm := &armcompute.VirtualMachine{}
	err := json.Unmarshal(resourceDesc, vm)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal resource description:%+v", err)
	}

	// Some validations on the VM
	if vm.Location == nil || vm.Properties == nil {
		return nil, fmt.Errorf("resource description is missing location or properties")
	}

	// Reject VMs that already have network interfaces
	if vm.Properties.NetworkProfile != nil && vm.Properties.NetworkProfile.NetworkInterfaces != nil {
		return nil, fmt.Errorf("resource description cannot contain network interface")
	}

	return vm, nil
}

// getInvisinetsResourceName returns a name for the Invisinets resource
func getInvisinetsResourceName(resourceType string) string {
	// TODO @nnomier: change based on invisinets naming convention
	return invisinetsPrefix + "-" + resourceType + "-" + uuid.New().String()
}

// getNSGRuleName returns a name for the Invisinets rule
func getNSGRuleName(ruleName string) string {
	return invisinetsPrefix + "-" + ruleName
}

func getRuleNameFromNSGRuleName(ruleName string) string {
	return strings.TrimPrefix(ruleName, invisinetsPrefix+"-")
}

// getResourceIDInfo parses the resourceID to extract subscriptionID and resourceGroupName (and VM name if needed)
// and returns a ResourceIDInfo object filled with the extracted values
// a valid resourceID should be in the format of '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/...'
func getResourceIDInfo(resourceID string) (ResourceIDInfo, error) {
	parts := strings.Split(resourceID, "/")
	if len(parts) < 5 {
		return ResourceIDInfo{}, fmt.Errorf("invalid resource ID format: expected at least 5 parts in the format of '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/...', got %d", len(parts))
	}

	if parts[0] != "" || parts[1] != "subscriptions" || parts[3] != "resourceGroups" {
		return ResourceIDInfo{}, fmt.Errorf("invalid resource ID format: expected '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/...', got '%s'", resourceID)
	}

	info := ResourceIDInfo{
		SubscriptionID:    parts[2],
		ResourceGroupName: parts[4],
	}

	info.ResourceName = parts[len(parts)-1]

	return info, nil
}

// checkAndCreatePeering checks whether the given rule has a tag that is in the address space of any of the invisinets vnets
// and if requires a peering or not
func (s *azurePluginServer) checkAndCreatePeering(ctx context.Context, resourceVnet *armnetwork.VirtualNetwork, rule *invisinetspb.PermitListRule, invisinetsVnetsMap map[string]string, namespace string) error {
	for _, target := range rule.Targets {
		isTagInResourceAddressSpace, err := utils.IsPermitListRuleTagInAddressSpace(target, *resourceVnet.Properties.AddressSpace.AddressPrefixes[0])
		if err != nil {
			return err
		}
		if isTagInResourceAddressSpace {
			continue
		}

		// if the tag is not in the resource address space, then check on the other invisinets vnets
		// if it matches one of them, then a peering is required (if it doesn't exist already)
		for vnetLocation, vnetAddressSpace := range invisinetsVnetsMap {
			isTagInVnetAddressSpace, err := utils.IsPermitListRuleTagInAddressSpace(target, vnetAddressSpace)
			if err != nil {
				return err
			}
			if isTagInVnetAddressSpace {
				peeringExists := false
				for _, peeredVnet := range resourceVnet.Properties.VirtualNetworkPeerings {
					if strings.HasSuffix(*peeredVnet.Properties.RemoteVirtualNetwork.ID, getVnetName(vnetLocation, namespace)) {
						peeringExists = true
						break
					}
				}

				if !peeringExists {
					err := s.azureHandler.CreateVnetPeering(ctx, getVnetName(vnetLocation, namespace), getVnetName(*resourceVnet.Location, namespace))
					if err != nil {
						return err
					}
				}
				break // No need to continue checking other vnets
			}
		}

		// TODO @nnomier: if the tag is not in any of the invisinets vnets (peeringExists = false), this might mean it's remote (another cloud),
		// so the multicloud setup could be checked/achieved here
	}
	return nil
}

func getInvisinetsNamespacePrefix(namespace string) string {
	return invisinetsPrefix + "-" + namespace
}

// getVnetName returns the name of the invisinets vnet in the given location
// since an invisients vnet is unique per location
func getVnetName(location string, namespace string) string {
	return getInvisinetsNamespacePrefix(namespace) + "-" + location + "-vnet"
}

func getVpnGatewayVnetName(namespace string) string {
	return getVpnGatewayName(namespace) + "-vnet"
}

func getVpnGatewayName(namespace string) string {
	return getInvisinetsNamespacePrefix(namespace) + "-vpn-gw"
}

func getVPNGatewayIPAddressName(namespace string, idx int) string {
	return getVpnGatewayName(namespace) + "-ip-" + strconv.Itoa(idx)
}

func getLocalNetworkGatewayName(namespace string, cloud string, idx int) string {
	return getInvisinetsNamespacePrefix(namespace) + "-" + cloud + "-local-gw-" + strconv.Itoa(idx)
}

func getVirtualNetworkGatewayConnectionName(namespace string, cloud string, idx int) string {
	return getInvisinetsNamespacePrefix(namespace) + "-" + cloud + "-conn-" + strconv.Itoa(idx)
}

func (s *azurePluginServer) CreateVpnGateway(ctx context.Context, req *invisinetspb.CreateVpnGatewayRequest) (*invisinetspb.CreateVpnGatewayResponse, error) {
	resourceId := req.Deployment.Id
	namespace := req.Deployment.Namespace
	resourceIdInfo, err := getResourceIDInfo(resourceId)
	if err != nil {
		return nil, fmt.Errorf("unable to get resource ID info: %w", err)
	}
	err = s.setupAzureHandler(resourceIdInfo)
	if err != nil {
		return nil, fmt.Errorf("unable to setup azure handler: %w", err)
	}

	vpnNumConnections := utils.GetNumVpnConnections(req.Cloud, utils.AZURE)
	publicIPAddresses := make([]*armnetwork.PublicIPAddress, vpnNumConnections)
	virtualNetworkGatewayName := getVpnGatewayName(namespace)
	virtualNetworkGateway, err := s.azureHandler.GetVirtualNetworkGateway(ctx, virtualNetworkGatewayName)
	var asn uint32
	if err != nil {
		if isErrorNotFound(err) {
			// Create two public IP addresses (need a second for active-active mode)
			publicIPAddressParameters := armnetwork.PublicIPAddress{
				Location: to.Ptr(vpnLocation),
				Properties: &armnetwork.PublicIPAddressPropertiesFormat{
					PublicIPAddressVersion:   to.Ptr(armnetwork.IPVersionIPv4),
					PublicIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodStatic),
				},
				SKU: &armnetwork.PublicIPAddressSKU{
					Name: to.Ptr(armnetwork.PublicIPAddressSKUNameStandard),
				},
			}
			for i := 0; i < vpnNumConnections; i++ {
				vpnGatewayIPAddressName := getVPNGatewayIPAddressName(namespace, i)
				publicIPAddress, err := s.azureHandler.GetPublicIPAddress(ctx, vpnGatewayIPAddressName)
				if err != nil {
					if isErrorNotFound(err) {
						publicIPAddress, err = s.azureHandler.CreatePublicIPAddress(ctx, vpnGatewayIPAddressName, publicIPAddressParameters)
						if err != nil {
							return nil, fmt.Errorf("unable to create public IP address: %w", err)
						}
					} else {
						return nil, fmt.Errorf("unable to get public IP address: %w", err)
					}
				}
				publicIPAddresses[i] = publicIPAddress
			}

			// Get VPN gateway subnet
			gatewayVnetName := getVpnGatewayVnetName(namespace)
			vpnGwSubnet, err := s.azureHandler.GetSubnet(ctx, gatewayVnetName, gatewaySubnetName)
			if err != nil {
				return nil, fmt.Errorf("unable to get VPN gateway subnet: %w", err)
			}

			conn, err := grpc.Dial(s.orchestratorServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return nil, fmt.Errorf("unable to establish connection with orchestrator: %w", err)
			}
			defer conn.Close()
			client := invisinetspb.NewControllerClient(conn)
			findUnusedAsnResp, err := client.FindUnusedAsn(ctx, &invisinetspb.FindUnusedAsnRequest{})
			if err != nil {
				return nil, fmt.Errorf("unable to find unused address space: %w", err)
			}
			asn = findUnusedAsnResp.Asn

			// Create VPN gateway
			virtualNetworkGatewayParameters := armnetwork.VirtualNetworkGateway{
				Location: to.Ptr(vpnLocation),
				Properties: &armnetwork.VirtualNetworkGatewayPropertiesFormat{
					Active: to.Ptr(true),
					BgpSettings: &armnetwork.BgpSettings{
						Asn: to.Ptr(int64(asn)),
					},
					EnableBgp:              to.Ptr(true),
					EnablePrivateIPAddress: to.Ptr(false),
					GatewayType:            to.Ptr(armnetwork.VirtualNetworkGatewayTypeVPN),
					IPConfigurations: []*armnetwork.VirtualNetworkGatewayIPConfiguration{
						{
							Name:       to.Ptr("default"),
							Properties: &armnetwork.VirtualNetworkGatewayIPConfigurationPropertiesFormat{},
						},
					},
					SKU: &armnetwork.VirtualNetworkGatewaySKU{
						Name: to.Ptr(armnetwork.VirtualNetworkGatewaySKUNameVPNGw1),
						Tier: to.Ptr(armnetwork.VirtualNetworkGatewaySKUTierVPNGw1),
					},
					VPNGatewayGeneration: to.Ptr(armnetwork.VPNGatewayGenerationGeneration1),
					VPNType:              to.Ptr(armnetwork.VPNTypeRouteBased),
				},
			}
			virtualNetworkGatewayParameters.Properties.IPConfigurations = make([]*armnetwork.VirtualNetworkGatewayIPConfiguration, vpnNumConnections)
			ipConfigurationNames := []string{"default", "activeActive"} // TODO @seankimkdy: come up with better naming convention ... ? (these are Azure defaults so they may rely on them actually)
			for i := 0; i < vpnNumConnections; i++ {
				virtualNetworkGatewayParameters.Properties.IPConfigurations[i] = &armnetwork.VirtualNetworkGatewayIPConfiguration{
					Name: to.Ptr(ipConfigurationNames[i]),
					Properties: &armnetwork.VirtualNetworkGatewayIPConfigurationPropertiesFormat{
						PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
						PublicIPAddress: &armnetwork.SubResource{
							ID: publicIPAddresses[i].ID,
						},
						Subnet: &armnetwork.SubResource{
							ID: vpnGwSubnet.ID,
						},
					},
				}
			}
			virtualNetworkGateway, err = s.azureHandler.CreateOrUpdateVirtualNetworkGateway(ctx, virtualNetworkGatewayName, virtualNetworkGatewayParameters)
			if err != nil {
				return nil, fmt.Errorf("unable to create virtual network gateway: %w", err)
			}

			// Add BGP IP addresses
			virtualNetworkGateway.Properties.BgpSettings.BgpPeeringAddresses = make([]*armnetwork.IPConfigurationBgpPeeringAddress, vpnNumConnections)
			for i := 0; i < vpnNumConnections; i++ {
				virtualNetworkGateway.Properties.BgpSettings.BgpPeeringAddresses[i] = &armnetwork.IPConfigurationBgpPeeringAddress{
					CustomBgpIPAddresses: []*string{to.Ptr(req.BgpPeeringIpAddresses[i])},
					IPConfigurationID:    virtualNetworkGateway.Properties.IPConfigurations[i].ID,
				}
			}
			_, err = s.azureHandler.CreateOrUpdateVirtualNetworkGateway(ctx, virtualNetworkGatewayName, *virtualNetworkGateway)
			if err != nil {
				return nil, fmt.Errorf("unable to update virtual network gateway with BGP IP addresses: %w", err)
			}

			// Update existing peerings with gateway transit relationship
			gatewayVnetPeerings, err := s.azureHandler.ListVirtualNetworkPeerings(ctx, gatewayVnetName)
			if err != nil {
				return nil, fmt.Errorf("unable to get peerings of virtual gateway vnet: %w", err)
			}
			for _, gatewayVnetToVnetPeering := range gatewayVnetPeerings {
				vnetResourceIDInfo, err := getResourceIDInfo(*gatewayVnetToVnetPeering.Properties.RemoteVirtualNetwork.ID)
				if err != nil {
					return nil, fmt.Errorf("unable to parse vnet resource ID from the gateway vnet to vnet peering: %w", err)
				}
				vnetName := vnetResourceIDInfo.ResourceName
				vnetToGatewayVnetPeering, err := s.azureHandler.GetVirtualNetworkPeering(ctx, vnetName, getPeeringName(vnetName, gatewayVnetName))
				if err != nil {
					return nil, fmt.Errorf("unable to get vnet to gateway vnet peering: %w", err)
				}
				err = s.azureHandler.CreateOrUpdateVnetPeeringRemoteGateway(ctx, vnetName, gatewayVnetName, vnetToGatewayVnetPeering, gatewayVnetToVnetPeering)
				if err != nil {
					return nil, fmt.Errorf("unable to update peerings between vnet and gateway vnet for VPN gateway transit: %w", err)
				}
			}
		} else {
			return nil, fmt.Errorf("unable to get virtual network gateway: %w", err)
		}
	} else {
		// Retrieve VPN gateway ASN and IP addresses
		asn = uint32(*virtualNetworkGateway.Properties.BgpSettings.Asn)
		for i, ipConfiguration := range virtualNetworkGateway.Properties.IPConfigurations {
			publicIPAddressIdInfo, err := getResourceIDInfo(*ipConfiguration.Properties.PublicIPAddress.ID)
			if err != nil {
				return nil, fmt.Errorf("unable to get public IP address ID info: %w", err)
			}
			publicIPAddress, err := s.azureHandler.GetPublicIPAddress(ctx, publicIPAddressIdInfo.ResourceName)
			if err != nil {
				return nil, fmt.Errorf("unable to get public IP address: %w", err)
			}
			publicIPAddresses[i] = publicIPAddress
		}
	}

	resp := &invisinetspb.CreateVpnGatewayResponse{Asn: asn}
	resp.GatewayIpAddresses = make([]string, vpnNumConnections)
	for i := 0; i < vpnNumConnections; i++ {
		resp.GatewayIpAddresses[i] = *publicIPAddresses[i].Properties.IPAddress
	}
	return resp, nil
}

func (s *azurePluginServer) CreateVpnConnections(ctx context.Context, req *invisinetspb.CreateVpnConnectionsRequest) (*invisinetspb.BasicResponse, error) {
	resourceIdInfo, err := getResourceIDInfo(req.Deployment.Id)
	if err != nil {
		return nil, fmt.Errorf("unable to get resource ID info: %w", err)
	}
	err = s.setupAzureHandler(resourceIdInfo)
	if err != nil {
		return nil, fmt.Errorf("unable to setup azure handler: %w", err)
	}

	vpnNumConnections := utils.GetNumVpnConnections(req.Cloud, utils.AZURE)
	localNetworkGateways := make([]*armnetwork.LocalNetworkGateway, vpnNumConnections)
	for i := 0; i < vpnNumConnections; i++ {
		localNetworkGatewayName := getLocalNetworkGatewayName(req.Deployment.Namespace, req.Cloud, i)
		localNetworkGateway, err := s.azureHandler.GetLocalNetworkGateway(ctx, localNetworkGatewayName)
		if err != nil {
			if isErrorNotFound(err) {
				localNetworkGatewayParameters := armnetwork.LocalNetworkGateway{
					Properties: &armnetwork.LocalNetworkGatewayPropertiesFormat{
						BgpSettings: &armnetwork.BgpSettings{
							Asn:               to.Ptr(int64(req.Asn)),
							BgpPeeringAddress: to.Ptr(req.BgpIpAddresses[i]),
							PeerWeight:        to.Ptr(int32(0)),
						},
						GatewayIPAddress: to.Ptr(req.GatewayIpAddresses[i]),
					},
					Location: to.Ptr(vpnLocation),
				}
				localNetworkGateway, err = s.azureHandler.CreateLocalNetworkGateway(ctx, localNetworkGatewayName, localNetworkGatewayParameters)
				if err != nil {
					return nil, fmt.Errorf("unable to create local network gateway: %w", err)
				}
			} else {
				return nil, fmt.Errorf("unable to get local network gateway: %w", err)
			}
		}
		localNetworkGateways[i] = localNetworkGateway
	}

	virtualNetworkGateway, err := s.azureHandler.GetVirtualNetworkGateway(ctx, getVpnGatewayName(req.Deployment.Namespace))
	if err != nil {
		return nil, fmt.Errorf("unable to get virtual network gateway: %w", err)
	}
	for i := 0; i < vpnNumConnections; i++ {
		virtualNetworkGatewayconnectionName := getVirtualNetworkGatewayConnectionName(req.Deployment.Namespace, req.Cloud, i)
		_, err := s.azureHandler.GetVirtualNetworkGatewayConnection(ctx, virtualNetworkGatewayconnectionName)
		if err != nil {
			if isErrorNotFound(err) {
				// Checks if a virtual network gateway connection already exists. Even though CreateOrUpdate is a PUT (i.e. idempotent),
				// a new random shared key is generated upon every call to this method from the controller server. Therefore, we don't
				// want to update the shared key since some other cloud plugins (e.g. GCP) will not update the shared key due to POST
				// semantics (i.e. GCP will not update the shared key).
				virtualNetworkGatewayConnectionParameters := &armnetwork.VirtualNetworkGatewayConnection{
					Properties: &armnetwork.VirtualNetworkGatewayConnectionPropertiesFormat{
						ConnectionType:                 to.Ptr(armnetwork.VirtualNetworkGatewayConnectionTypeIPsec),
						VirtualNetworkGateway1:         virtualNetworkGateway,
						ConnectionMode:                 to.Ptr(armnetwork.VirtualNetworkGatewayConnectionModeDefault),
						ConnectionProtocol:             to.Ptr(armnetwork.VirtualNetworkGatewayConnectionProtocolIKEv2),
						DpdTimeoutSeconds:              to.Ptr(int32(45)),
						EnableBgp:                      to.Ptr(true),
						IPSecPolicies:                  []*armnetwork.IPSecPolicy{},
						LocalNetworkGateway2:           localNetworkGateways[i],
						RoutingWeight:                  to.Ptr(int32(0)),
						SharedKey:                      to.Ptr(req.SharedKey),
						TrafficSelectorPolicies:        []*armnetwork.TrafficSelectorPolicy{},
						UseLocalAzureIPAddress:         to.Ptr(false),
						UsePolicyBasedTrafficSelectors: to.Ptr(false),
					},
					Location: to.Ptr(vpnLocation),
				}
				_, err := s.azureHandler.CreateVirtualNetworkGatewayConnection(ctx, virtualNetworkGatewayconnectionName, *virtualNetworkGatewayConnectionParameters)
				if err != nil {
					return nil, fmt.Errorf("unable to create virtual network gateway connection: %w", err)
				}
			} else {
				return nil, fmt.Errorf("unable to get virtual network gateway connection: %w", err)
			}
		}
	}

	return &invisinetspb.BasicResponse{Success: true}, nil
}

func Setup(port int, orchestratorServerAddr string) *azurePluginServer {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	azureServer := &azurePluginServer{
		azureHandler:           &azureSDKHandler{},
		orchestratorServerAddr: orchestratorServerAddr,
	}
	invisinetspb.RegisterCloudPluginServer(grpcServer, azureServer)
	fmt.Println("Starting server on port :", port)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			fmt.Println(err.Error())
		}
	}()
	return azureServer
}
