// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package provider

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

// This applies the differences between the regular schema and the one
// we supply for default resources, and returns the schema for a default resource
func ConvertToDefaultVcnResourceSchema(resourceSchema *schema.Resource) *schema.Resource {
	if resourceSchema == nil {
		return nil
	}

	resourceSchema.Importer = &schema.ResourceImporter{
		State: ImportDefaultVcnResource,
	}

	resourceSchema.Schema["manage_default_resource_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	}

	delete(resourceSchema.Schema, "compartment_id")
	delete(resourceSchema.Schema, "vcn_id")

	return resourceSchema
}

func ImportDefaultVcnResource(d *schema.ResourceData, value interface{}) ([]*schema.ResourceData, error) {
	err := d.Set("manage_default_resource_id", d.Id())
	return []*schema.ResourceData{d}, err
}

func validateNotEmptyString() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}
		if len(v) == 0 {
			es = append(es, fmt.Errorf("%s cannot be an empty string", k))
		}
		return
	}
}
