package ibm

import (
	"context"
	"fmt"
	"net"
	"os"

	logger "github.com/NetSys/invisinets/pkg/logger"

	sdk "github.com/NetSys/invisinets/pkg/ibm_plugin/sdk"
	"github.com/NetSys/invisinets/pkg/invisinetspb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ibmPluginServer struct {
	invisinetspb.UnimplementedCloudPluginServer
	cloudClient        *sdk.IBMCloudClient
	frontendServerAddr string
}

func (s *ibmPluginServer) setupCloudClient(region string) error {
	client, err := sdk.NewIbmCloudClient(region)
	if err != nil {
		logger.Log.Println("Failed to set up IBM clients with error:", err)
		return err
	}
	s.cloudClient = client
	return nil
}

// Creates the specified resource. Currently only supports instance creation.
// Default instance profile is 2CPU, 8GB RAM, unless specified.
// Default instance name will be auto-generated unless specified.
func (s *ibmPluginServer) CreateResource(c context.Context, resourceDesc *invisinetspb.ResourceDescription) (*invisinetspb.BasicResponse, error) {
	var vpcID string
	var subnetID string

	vmFields, err := getInstanceData(resourceDesc)
	if err != nil {
		return nil, err
	}

	region, err := sdk.Zone2Region(vmFields.Zone)
	if err != nil {
		logger.Log.Println("Invalid region:", region)
		return nil, err
	}
	err = s.setupCloudClient(region)
	if err != nil {
		return nil, err
	}
	/* TODO: Future support in multiple deployments and multiple vpcs
	in single region will require adding deployment ID as a tag
	*/
	vpcIDs, err := s.cloudClient.GetInvisinetsTaggedResources(sdk.VPC, nil,
		sdk.ResourceQuery{Region: region})
	if err != nil {
		return nil, err
	}

	// use existing invisinets VPC or create a new one
	if len(vpcIDs) != 0 {
		// currently assuming a single VPC per region
		vpcID = vpcIDs[0]
		logger.Log.Printf("Reusing invisinets VPC with ID: %v in region %v", vpcID, region)
	} else {
		conn, err := grpc.Dial(s.frontendServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		client := invisinetspb.NewControllerClient(conn)
		response, err := client.FindUnusedAddressSpace(context.Background(), &invisinetspb.Empty{})
		if err != nil {
			return nil, err
		}
		// create a vpc with a subnet in each zone
		vpc, err := s.cloudClient.CreateVpc("", response.Address)
		if err != nil {
			return nil, err
		}
		vpcID = *vpc.ID
	}

	// look for an invisinets subnet that's tagged with the above VPC ID.
	requiredTags := []string{vpcID}
	subnetsIDs, err := s.cloudClient.GetInvisinetsTaggedResources(sdk.SUBNET, requiredTags,
		sdk.ResourceQuery{Zone: vmFields.Zone})
	if err != nil {
		return nil, err
	}
	if len(subnetsIDs) != 0 {
		// at least one invisinets subnet that fits the query was found. Use a random one.
		subnetID = subnetsIDs[0]
	} else {
		// No invisinets subnets were found matching vpc and zone, currently er
		return nil, fmt.Errorf("invisinets subnet wasn't found")
	}

	vm, err := s.cloudClient.CreateVM(vpcID, subnetID,
		vmFields.Zone, vmFields.Name, vmFields.Profile)
	if err != nil {
		return nil, err
	}
	return &invisinetspb.BasicResponse{Success: true, Message: "successfully created VM",
		UpdatedResource: &invisinetspb.ResourceID{Id: *vm.ID}}, nil
}

// returns a list of address spaces used by either user's or invisinets' sunbets,
// for each invisinets vpc.
func (s *ibmPluginServer) GetUsedAddressSpaces(ctx context.Context, deployment *invisinetspb.InvisinetsDeployment) (*invisinetspb.AddressSpaceList, error) {
	var invisinetsAddressSpaces []string
	err := s.setupCloudClient("")
	if err != nil {
		return nil, err
	}
	// get all VPCs in the deployment.
	// TODO future multi deployment support will require sending deployment id as tag, currently using static tag.
	deploymentVpcIDs, err := s.cloudClient.GetInvisinetsTaggedResources(sdk.VPC, nil, sdk.ResourceQuery{})
	if err != nil {
		return nil, err
	}
	// for each vpc, collect the address space of all subnets, including users'.
	for _, vpcID := range deploymentVpcIDs {
		subnets, err := s.cloudClient.GetSubnetsInVPC(vpcID)
		if err != nil {
			return nil, err
		}
		for _, subnet := range subnets {
			invisinetsAddressSpaces = append(invisinetsAddressSpaces, *subnet.Ipv4CIDRBlock)
		}
	}
	return &invisinetspb.AddressSpaceList{AddressSpaces: invisinetsAddressSpaces}, nil
}

/*
returns security rules of security groups associated with the specified instance.
Invisinets Assumptions:
- single security group per instance.
- ignoring possible user security groups.
*/
func (s *ibmPluginServer) GetPermitList(ctx context.Context, resourceID *invisinetspb.ResourceID) (*invisinetspb.PermitList, error) {
	permitList := &invisinetspb.PermitList{
		AssociatedResource: resourceID.Id,
		Rules:              []*invisinetspb.PermitListRule{},
	}
	resourceIDInfo, err := getResourceIDInfo(resourceID)
	if err != nil {
		return nil, err
	}
	region, vmID := resourceIDInfo.Region, resourceIDInfo.ResourceID
	if !sdk.IsRegionValid(region) {
		return nil, fmt.Errorf("region %v isn't valid", region)
	}
	err = s.setupCloudClient(region)
	if err != nil {
		return nil, err
	}

	securityGroups, err := s.cloudClient.GetSecurityGroupsOfVM(vmID)
	if err != nil {
		return nil, err
	}
	rulesHashValues := map[int64]bool{}
	for _, sgID := range securityGroups {
		// Currently adding rules for all security groups attached to VM.
		ibmRules, err := s.cloudClient.GetSecurityRulesOfSG(sgID)
		if err != nil {
			return nil, err
		}
		invisinetsRules, err := sgRules2InvisinetsRules(ibmRules)
		if err != nil {
			return nil, err
		}
		// Avoid adding duplicate rules.
		for _, rule := range invisinetsRules {
			// exclude unique field "Id" from hash calculation.
			ruleHashValue, err := getStructHash(*rule, []string{"Id"})
			if err != nil {
				return nil, err
			}
			if _, ruleExists := rulesHashValues[int64(ruleHashValue)]; !ruleExists {
				permitList.Rules = append(permitList.Rules, rule)
				rulesHashValues[int64(ruleHashValue)] = true
			}
		}
	}
	return permitList, nil
}

// starts up the plugin server and stores the frontend server address.
func Setup(port int, frontendAddress string) {
	pluginServerAddress := "localhost"
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%d", pluginServerAddress, port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	ibmServer := ibmPluginServer{
		cloudClient:        &sdk.IBMCloudClient{},
		frontendServerAddr: fmt.Sprintf("%v:%v", frontendAddress, port),
	}
	invisinetspb.RegisterCloudPluginServer(grpcServer, &ibmServer)
	fmt.Printf("Starting plugin server on: %v:%v\n", pluginServerAddress, port)
	fmt.Println("Received frontend Server address:", frontendAddress)
	err = grpcServer.Serve(lis)
	if err != nil {
		fmt.Println(err.Error())
	}
}