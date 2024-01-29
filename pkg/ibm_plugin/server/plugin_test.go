//go:build ibm

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

package ibm

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"testing"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/require"

	fake "github.com/NetSys/invisinets/pkg/fake/controller/rpc"
	sdk "github.com/NetSys/invisinets/pkg/ibm_plugin/sdk"
	"github.com/NetSys/invisinets/pkg/invisinetspb"
	utils "github.com/NetSys/invisinets/pkg/utils"
)

var testResGroupName = flag.String("sg", "pywren", "Name of the user's security group")
var testResourceIDUSEast1 string
var testResourceIDUSEast2 string
var testResourceIDEUDE1 string
var testResourceIDUSSouth1 string

func TestMain(m *testing.M) {
	flag.Parse()
	testResourceIDUSEast1 = "/ResourceGroupName/" + *testResGroupName + "/Zone/" + testZoneUSEast1 + "/ResourceID/" + testInstanceNameUSEast1
	testResourceIDUSEast2 = "/ResourceGroupName/" + *testResGroupName + "/Zone/" + testZoneUSEast2 + "/ResourceID/" + testInstanceNameUSEast2
	testResourceIDEUDE1 = "/ResourceGroupName/" + *testResGroupName + "/Zone/" + testZoneEUDE1 + "/ResourceID/" + testInstanceNameEUDE1
	testResourceIDUSSouth1 = "/ResourceGroupName/" + *testResGroupName + "/Zone/" + testZoneUSSouth1 + "/ResourceID/" + testInstanceNameUSSouth1
	exitCode := m.Run()
	os.Exit(exitCode)
}

const (
	testUSEastRegion        = "us-east"
	testUSSouthRegion        = "us-south"
	testEURegion        = "eu-de"
	testZoneUSEast1     = testUSEastRegion + "-1"
	testZoneUSEast2     = testUSEastRegion + "-2"
	testZoneUSSouth1     = testUSSouthRegion + "-1"
	testZoneEUDE1       = testEURegion + "-1"
	testInstanceNameUSEast1 = "invisinets-vm-east-1"
	testInstanceNameUSEast2 = "invisinets-vm-east-2"
	testInstanceNameUSSouth1 = "invisinets-vm-south-1"
	testInstanceNameEUDE1 = "invisinets-vm-de-1"

	testImageUSEast = "r014-0acbdcb5-a68f-4a52-98ea-4da4fe89bacb" // us-east Ubuntu 22.04
	testImageEUDE   = "r010-f68ef7b3-1c5e-4ef7-8040-7ae0f5bf04fd" // eu-de Ubuntu 22.04
	testImageUSSouth   = "r006-01deb923-46f6-44c3-8fdc-99d8493d2464" // us-south Ubuntu 22.04
	testProfile     = "bx2-2x8"
	testNamespace   = "inv-namespace"
)

// permit list example
var testPermitList []*invisinetspb.PermitListRule = []*invisinetspb.PermitListRule{
	//TCP protocol rules
	{
		Direction: invisinetspb.Direction_INBOUND,
		SrcPort:   443,
		DstPort:   443,
		Protocol:  6,
		Targets:   []string{"10.0.0.0/18"},
	},
	{
		Direction: invisinetspb.Direction_OUTBOUND,
		SrcPort:   8080,
		DstPort:   8080,
		Protocol:  6,
		Targets:   []string{"10.0.128.12", "10.0.128.13"},
	},
	//All protocol rules
	{
		Direction: invisinetspb.Direction_INBOUND,
		SrcPort:   -1,
		DstPort:   -1,
		Protocol:  -1,
		Targets:   []string{"10.0.64.0/22", "10.0.64.0/24"},
	},
	{
		Direction: invisinetspb.Direction_OUTBOUND,
		SrcPort:   -1,
		DstPort:   -1,
		Protocol:  -1,
		Targets:   []string{"10.0.64.1"},
	},
}

// TODO @praveingk: Change the tests to use fake IBM handlers

// go test --tags=ibm -run TestCreateResourceNewVPC -sg=<security group name>
func TestCreateResourceNewVPC(t *testing.T) {
	// Notes for tester:
	// to change region set the values below according to constants above, e.g.:
	// - test arguments for EU-DE-1:
	// image, zone, instanceName, resourceID := testImageEUDE, testZoneEUDE1, testInstanceNameEUDE1, testResourceIDEUDE1
	// - test arguments for us-east-2:
	// image, zone, instanceName, resourceID := testImageUSEast, testZoneUSEast2, testInstanceNameUSEast2, testResourceIDUSEast2
	// - test arguments for us-south-1:
	// image, zone, instanceName, resourceID := testImageUSSouth, testZoneUSSouth1, testInstanceNameUSSouth1, testResourceIDUSSouth1
	image, zone, instanceName, resourceID := testImageUSEast, testZoneUSEast1, testInstanceNameUSEast1, testResourceIDUSEast1
	
	_, fakeControllerServerAddr, err := fake.SetupFakeOrchestratorRPCServer(utils.IBM)
	if err != nil {
		t.Fatal(err)
	}
	imageIdentity := vpcv1.ImageIdentityByID{ID: &image}
	zoneIdentity := vpcv1.ZoneIdentityByName{Name: &zone}
	myTestProfile := string(testProfile)

	testPrototype := &vpcv1.InstancePrototypeInstanceByImage{
		Image:   &imageIdentity,
		Zone:    &zoneIdentity,
		Name:    core.StringPtr(instanceName),
		Profile: &vpcv1.InstanceProfileIdentityByName{Name: &myTestProfile},
	}

	s := &ibmPluginServer{
		frontendServerAddr: fakeControllerServerAddr,
		cloudClient:        make(map[string]*sdk.CloudClient)}

	description, err := json.Marshal(vpcv1.CreateInstanceOptions{InstancePrototype: vpcv1.InstancePrototypeIntf(testPrototype)})
	require.NoError(t, err)

	resource := &invisinetspb.ResourceDescription{Id: resourceID, Description: description, Namespace: testNamespace}
	resp, err := s.CreateResource(context.Background(), resource)
	if err != nil {
		println(err)
	}
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// This func tests creating a new VM in an existing region, ergo to properly test:
// 1. Have an invisinets VPC deployed beforehand.
// 2. create the new VM in the same region as the deployed VPC.
// go test --tags=ibm -run TestCreateResourceExistingVPC -sg=<security group name>
func TestCreateResourceExistingVPC(t *testing.T) {
	// Notes for tester:
	// to change region set the values below according to constants above, e.g.:
	// - test arguments for EU-DE-1:
	//   image, zone, instanceName, resourceID := testImageEUDE, testZoneEUDE1, testInstanceNameEU1, testResourceIDEUDE1
	// - test arguments for US-EAST-1:
	// image, zone, instanceName, resourceID := testImageUSEast, testZoneUSEast1, testInstanceNameUS1, testResourceIDUSEast1
	// - test arguments for US-EAST-2:
	image, zone, instanceName, resourceID := testImageUSEast, testZoneUSEast2, testInstanceNameUSEast2, testResourceIDUSEast2

	_, fakeControllerServerAddr, err := fake.SetupFakeOrchestratorRPCServer(utils.IBM)
	if err != nil {
		t.Fatal(err)
	}
	imageIdentity := vpcv1.ImageIdentityByID{ID: core.StringPtr(image)}
	zoneIdentity := vpcv1.ZoneIdentityByName{Name: core.StringPtr(zone)}
	myTestProfile := string(testProfile)

	testPrototype := &vpcv1.InstancePrototypeInstanceByImage{
		Image:   &imageIdentity,
		Zone:    &zoneIdentity,
		Name:    core.StringPtr(instanceName),
		Profile: &vpcv1.InstanceProfileIdentityByName{Name: &myTestProfile},
	}

	s := &ibmPluginServer{
		frontendServerAddr: fakeControllerServerAddr,
		cloudClient:        make(map[string]*sdk.CloudClient)}
	description, err := json.Marshal(vpcv1.CreateInstanceOptions{InstancePrototype: vpcv1.InstancePrototypeIntf(testPrototype)})
	require.NoError(t, err)

	resource := &invisinetspb.ResourceDescription{Id: resourceID, Description: description, Namespace: testNamespace}
	resp, err := s.CreateResource(context.Background(), resource)
	if err != nil {
		println(err)
	}
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// usage: go test --tags=ibm -run TestGetPermitList -sg=<security group name>
func TestGetPermitList(t *testing.T) {
	resourceID := testResourceIDUSEast1  // replace as needed with other IDs, e.g. testResourceIDEUDE1

	s := &ibmPluginServer{cloudClient: make(map[string]*sdk.CloudClient)}

	resp, err := s.GetPermitList(context.Background(), &invisinetspb.ResourceID{Id: resourceID,
		Namespace: testNamespace})
	require.NoError(t, err)
	require.NotNil(t, resp)

	b, err := json.MarshalIndent(resp, "", "  ")
	require.NoError(t, err)
	// Note: direction:0(inbound) will not be printed.
	utils.Log.Printf("Permit rules of instance %v are:\n%v", testInstanceNameUSEast1, string(b))
}

// usage: go test --tags=ibm -run TestAddPermitListRules -sg=<security group name>
func TestAddPermitListRules(t *testing.T) {
	resourceID := testResourceIDUSEast1  // replace as needed with other IDs, e.g. testResourceIDEUDE1

	permitList := &invisinetspb.PermitList{
		AssociatedResource: resourceID,
		Rules:              testPermitList,
		Namespace:          testNamespace,
	}

	s := &ibmPluginServer{cloudClient: make(map[string]*sdk.CloudClient)}

	resp, err := s.AddPermitListRules(context.Background(), permitList)
	require.NoError(t, err)
	require.NotNil(t, resp)

	utils.Log.Printf("Response: %v", resp)
}

// usage: go test --tags=ibm -run TestDeletePermitListRule -sg=<security group name>
func TestDeletePermitListRules(t *testing.T) {
	resourceID := testResourceIDUSEast1 // replace as needed with other IDs, e.g. testResourceIDUSSouth1

	permitList := &invisinetspb.PermitList{
		AssociatedResource: resourceID,
		Rules:              testPermitList,
		Namespace:          testNamespace,
	}

	s := &ibmPluginServer{cloudClient: make(map[string]*sdk.CloudClient)}

	resp, err := s.DeletePermitListRules(context.Background(), permitList)
	require.NoError(t, err)
	require.NotNil(t, resp)

	utils.Log.Printf("Response: %v", resp)
}

// usage: go test --tags=ibm -run TestGetUsedAddressSpaces -sg=<security group name>
// this function logs subnets' address spaces from all invisinets' VPCs.
func TestGetUsedAddressSpaces(t *testing.T) {
	// GetUsedAddressSpaces() is independent of any region, since it returns
	// address spaces in global scope, so any test resource ID will do.
	deployment := &invisinetspb.InvisinetsDeployment{Id: testResourceIDUSEast1}

	s := &ibmPluginServer{cloudClient: make(map[string]*sdk.CloudClient)}

	usedAddressSpace, err := s.GetUsedAddressSpaces(context.Background(), deployment)
	require.NoError(t, err)
	require.NotEmpty(t, usedAddressSpace)
}
