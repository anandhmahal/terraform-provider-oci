// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	oci_file_storage "github.com/oracle/oci-go-sdk/filestorage"
)

const (
	MountTargetRequiredOnlyResource = MountTargetResourceDependencies + `
resource "oci_file_storage_mount_target" "test_mount_target" {
	#Required
	availability_domain = "${var.mount_target_availability_domain}"
	compartment_id = "${var.compartment_id}"
	subnet_id = "${oci_core_subnet.test_subnet.id}"
}
`

	MountTargetResourceConfig = MountTargetResourceDependencies + `
resource "oci_file_storage_mount_target" "test_mount_target" {
	#Required
	availability_domain = "${var.mount_target_availability_domain}"
	compartment_id = "${var.compartment_id}"
	subnet_id = "${oci_core_subnet.test_subnet.id}"

	#Optional
	display_name = "${var.mount_target_display_name}"
	hostname_label = "${var.mount_target_hostname_label}"
	ip_address = "${var.mount_target_ip_address}"
}
`
	MountTargetPropertyVariables = `
variable "mount_target_availability_domain" { default = "kIdk:PHX-AD-1" }
variable "mount_target_display_name" { default = "mount-target-5" }
variable "mount_target_hostname_label" { default = "hostnameLabel" }
variable "mount_target_ip_address" { default = "10.0.1.5" } # Subnet CIDR = 10.0.0.0/16. This IP needs to be in the allowable range.

`
	MountTargetResourceDependencies = SubnetPropertyVariables + SubnetResourceConfig
)

func TestFileStorageMountTargetResource_basic(t *testing.T) {
	provider := testAccProvider
	config := testProviderConfig()

	compartmentId := getRequiredEnvSetting("compartment_id_for_create")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)
	compartmentId2 := getRequiredEnvSetting("compartment_id_for_update")
	compartmentIdVariableStr2 := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId2)

	resourceName := "oci_file_storage_mount_target.test_mount_target"
	datasourceName := "data.oci_file_storage_mount_targets.test_mount_targets"

	var resId, resId2 string

	resource.Test(t, resource.TestCase{
		Providers: map[string]terraform.ResourceProvider{
			"oci": provider,
		},
		Steps: []resource.TestStep{
			// verify create
			{
				ImportState:       true,
				ImportStateVerify: true,
				Config:            config + MountTargetPropertyVariables + compartmentIdVariableStr + MountTargetRequiredOnlyResource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "availability_domain", "kIdk:PHX-AD-1"),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						return err
					},
				),
			},

			// delete before next create
			{
				Config: config + compartmentIdVariableStr + MountTargetResourceDependencies,
			},
			// verify create with optionals
			{
				Config: config + MountTargetPropertyVariables + compartmentIdVariableStr + MountTargetResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "availability_domain", "kIdk:PHX-AD-1"),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "display_name", "mount-target-5"),
					resource.TestCheckResourceAttr(resourceName, "hostname_label", "hostnameLabel"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "ip_address", "10.0.1.5"),
					resource.TestCheckResourceAttr(resourceName, "private_ip_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip_ids.0"),
					resource.TestCheckResourceAttr(resourceName, "state", string(oci_file_storage.MountTargetLifecycleStateActive)),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "time_created"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						return err
					},
				),
			},

			// verify updates to updatable parameters
			{
				Config: config + `
variable "mount_target_availability_domain" { default = "kIdk:PHX-AD-1" }
variable "mount_target_display_name" { default = "displayName2" } # changing this value to test updates
variable "mount_target_hostname_label" { default = "hostnameLabel" }
variable "mount_target_ip_address" { default = "10.0.1.5" } # Subnet CIDR = 10.0.0.0/16. This IP needs to be in the allowable range.

                ` + compartmentIdVariableStr + MountTargetResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "availability_domain", "kIdk:PHX-AD-1"),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "display_name", "displayName2"),
					resource.TestCheckResourceAttr(resourceName, "hostname_label", "hostnameLabel"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "ip_address", "10.0.1.5"),
					resource.TestCheckResourceAttr(resourceName, "private_ip_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip_ids.0"),
					resource.TestCheckResourceAttr(resourceName, "state", string(oci_file_storage.MountTargetLifecycleStateActive)),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "time_created"),

					func(s *terraform.State) (err error) {
						resId2, err = fromInstanceState(s, resourceName, "id")
						if resId != resId2 {
							return fmt.Errorf("Resource recreated when it was supposed to be updated.")
						}
						return err
					},
				),
			},
			// verify updates to Force New parameters.
			{
				Config: config + `
variable "mount_target_availability_domain" { default = "kIdk:PHX-AD-1" }  # Subnet setup available in this AD only. So not changing.
variable "mount_target_display_name" { default = "displayName2" }
variable "mount_target_hostname_label" { default = "hostnameLabel2" }
variable "mount_target_ip_address" { default = "10.0.1.6" } # Subnet CIDR = 10.0.0.0/16. This IP needs to be in the allowable range.

                ` + compartmentIdVariableStr2 + MountTargetResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "availability_domain", "kIdk:PHX-AD-1"),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId2),
					resource.TestCheckResourceAttr(resourceName, "display_name", "displayName2"),
					resource.TestCheckResourceAttr(resourceName, "hostname_label", "hostnameLabel2"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "ip_address", "10.0.1.6"),
					resource.TestCheckResourceAttr(resourceName, "private_ip_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip_ids.0"),
					resource.TestCheckResourceAttr(resourceName, "state", string(oci_file_storage.MountTargetLifecycleStateActive)),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "time_created"),

					func(s *terraform.State) (err error) {
						resId2, err = fromInstanceState(s, resourceName, "id")
						if resId == resId2 {
							return fmt.Errorf("Resource was expected to be recreated but it wasn't.")
						}
						return err
					},
				),
			},
			// verify datasource
			{
				Config: config + `
variable "mount_target_availability_domain" { default = "kIdk:PHX-AD-1" }
variable "mount_target_display_name" { default = "displayName2" }
variable "mount_target_hostname_label" { default = "hostnameLabel2" }
variable "mount_target_ip_address" { default = "10.0.1.5" } # Subnet CIDR = 10.0.0.0/16. This IP needs to be in the allowable range.

data "oci_file_storage_mount_targets" "test_mount_targets" {
	#Required
	availability_domain = "${var.mount_target_availability_domain}"
	compartment_id = "${var.compartment_id}"

	#Optional
	display_name = "${var.mount_target_display_name}"
	id = "${oci_file_storage_mount_target.test_mount_target.id}"
	state = "ACTIVE"

    filter {
    	name = "id"
    	values = ["${oci_file_storage_mount_target.test_mount_target.id}"]
    }
}
                ` + compartmentIdVariableStr2 + MountTargetResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "availability_domain", "kIdk:PHX-AD-1"),
					resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId2),
					resource.TestCheckResourceAttr(datasourceName, "mount_targets.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "mount_targets.0.availability_domain", "kIdk:PHX-AD-1"),
					resource.TestCheckResourceAttr(datasourceName, "mount_targets.0.compartment_id", compartmentId2),
					resource.TestCheckResourceAttr(datasourceName, "mount_targets.0.display_name", "displayName2"),
					resource.TestCheckResourceAttrSet(datasourceName, "mount_targets.0.id"),
					resource.TestCheckResourceAttr(datasourceName, "mount_targets.0.private_ip_ids.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceName, "mount_targets.0.private_ip_ids.0"),
					resource.TestCheckResourceAttr(datasourceName, "mount_targets.0.state", string(oci_file_storage.MountTargetLifecycleStateActive)),
					resource.TestCheckResourceAttrSet(datasourceName, "mount_targets.0.subnet_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "mount_targets.0.export_set_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "mount_targets.0.time_created"),
				),
			},
		},
	})
}
