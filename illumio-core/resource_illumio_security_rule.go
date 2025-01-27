// Copyright 2021 Illumio, Inc. All Rights Reserved.

package illumiocore

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Jeffail/gabs/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/illumio/terraform-provider-illumio-core/models"
)

var (
	validSRResLabelAsValues = []string{"workloads", "virtual_services"}
	validSRIngSerProtos     = []string{"6", "17"}
	validSRConsumerActors   = []string{"ams", "container_host"}
	validSRProducerActors   = []string{"ams"}
)

func resourceIllumioSecurityRule() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceIllumioSecurityRuleRead,
		CreateContext: resourceIllumioSecurityRuleCreate,
		UpdateContext: resourceIllumioSecurityRuleUpdate,
		DeleteContext: resourceIllumioSecurityRuleDelete,
		SchemaVersion: version,
		Description:   "Manages Illumio Security Rule",
		Schema:        securityRuleResourceSchemaMap(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func securityRuleResourceSchemaMap() map[string]*schema.Schema {
	securityRuleSchema := securityRuleResourceBaseSchemaMap()
	securityRuleSchema["rule_set_href"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "URI of Rule set, in which security rule will be added",
	}
	return securityRuleSchema
}

func securityRuleResourceBaseSchemaMap() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "URI of Security Rule",
		},
		"enabled": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Enabled flag. Determines whether the rule will be enabled in ruleset or not",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description of Security Rule",
		},
		"external_data_set": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			Description:      "External data set identifier",
		},
		"external_data_reference": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			Description:      "External data reference identifier",
		},
		"ingress_services": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Collection of Ingress Service. If resolve_label_as.providers list includes \"workloads\" then ingress_service is required. Only one of the {\"href\"} or {\"proto\", \"port\", \"to_port\"} parameter combination is allowed",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"proto": {
						Type:             schema.TypeString,
						Optional:         true,
						Description:      "Protocol number. Allowed values are 6 (TCP) and 17 (UDP)",
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validSRIngSerProtos, true)),
					},
					"port": {
						Type:             schema.TypeString,
						Optional:         true,
						Description:      "Port number used with protocol or starting port when specifying a range. Allowed range is 0-65535",
						ValidateDiagFunc: isStringAPortNumber(),
					},
					"to_port": {
						Type:             schema.TypeString,
						Optional:         true,
						Description:      "Upper end of port range. Allowed range is 0-65535",
						ValidateDiagFunc: isStringAPortNumber(),
					},
					"href": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: isServiceHref,
						Description:      "URI of Service",
					},
				},
			},
		},
		"resolve_labels_as": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "resolve label as for Security rule",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"providers": {
						Type:        schema.TypeSet,
						Required:    true,
						MinItems:    1,
						Description: "providers for resolve_labels_as. Allowed values are \"workloads\", \"virtual_services\"",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"consumers": {
						Type:        schema.TypeSet,
						Required:    true,
						MinItems:    1,
						Description: "consumers for resolve_labels_as. Allowed values are \"workloads\", \"virtual_services\"",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"sec_connect": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Determines whether a secure connection is established. Default value: false",
		},
		"stateless": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Determines whether packet filtering is stateless for the rule. Default value: false",
		},
		"machine_auth": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Determines whether machine authentication is enabled. Default value: false",
		},
		"providers": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "providers for Security Rule. Only one actor can be specified in one providers block",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"actors": {
						Type:         schema.TypeString,
						Optional:     true,
						Description:  "actors for providers. Allowed value is \"ams\"",
						ValidateFunc: validation.StringInSlice(validSRProducerActors, false),
					},
					"label": {
						Type:        schema.TypeSet,
						Optional:    true,
						MaxItems:    1,
						Description: "Href of Label",
						Elem:        hrefSchemaRequired("Label", isLabelHref),
					},
					"label_group": {
						Type:        schema.TypeSet,
						Optional:    true,
						MaxItems:    1,
						Description: "Href of Label Group",
						Elem:        hrefSchemaRequired("Label Group", isLabelGroupHref),
					},
					"workload": {
						Type:        schema.TypeSet,
						Optional:    true,
						MaxItems:    1,
						Description: "Href of Workload",
						Elem:        hrefSchemaRequired("Workload", isWorkloadHref),
					},
					"virtual_service": {
						Type:        schema.TypeSet,
						Optional:    true,
						MaxItems:    1,
						Description: "Href of Virtual Service",
						Elem:        hrefSchemaRequired("Virtual Service", isVirtualServiceHref),
					},
					"virtual_server": {
						Type:        schema.TypeSet,
						Optional:    true,
						MaxItems:    1,
						Description: "Href of Virtual Server",
						Elem: hrefSchemaRequired("Virtual Server", validation.ToDiagFunc(
							validation.StringIsNotEmpty,
						)),
					},
					"ip_list": {
						Type:        schema.TypeSet,
						Optional:    true,
						MaxItems:    1,
						Description: "Href of IP List",
						Elem:        hrefSchemaRequired("IP List", isIPListHref),
					},
				},
			},
		},
		"consumers": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "Consumers for Security Rule. Only one actor can be specified in one consumer block",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"actors": {
						Type:         schema.TypeString,
						Optional:     true,
						Description:  "actors for consumers parameter. Allowed values are \"ams\" and \"container_host\"",
						ValidateFunc: validation.StringInSlice(validSRConsumerActors, false),
					},
					"label": {
						Type:        schema.TypeSet,
						Optional:    true,
						MaxItems:    1,
						Description: "Href of Label",
						Elem:        hrefSchemaRequired("Label", isLabelHref),
					},
					"label_group": {
						Type:        schema.TypeSet,
						Optional:    true,
						MaxItems:    1,
						Description: "Href of Label Group",
						Elem:        hrefSchemaRequired("Label Group", isLabelGroupHref),
					},
					"workload": {
						Type:        schema.TypeSet,
						Optional:    true,
						MaxItems:    1,
						Description: "Href of Workload",
						Elem:        hrefSchemaRequired("Workload", isWorkloadHref),
					},
					"virtual_service": {
						Type:        schema.TypeSet,
						Optional:    true,
						MaxItems:    1,
						Description: "Href of Virtual Service",
						Elem:        hrefSchemaRequired("Virtual Service", isVirtualServiceHref),
					},
					"ip_list": {
						Type:        schema.TypeSet,
						Optional:    true,
						MaxItems:    1,
						Description: "Href of IP List",
						Elem:        hrefSchemaRequired("IP List", isIPListHref),
					},
				},
			},
		},
		"unscoped_consumers": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "If false (the default), the created Rule will be an intra-scope rule. If true, it will be extra-scope. Default value: false",
		},
		"update_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type of update",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Timestamp when this security rule was first created",
		},
		"updated_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Timestamp when this security rule was last updated",
		},
		"deleted_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Timestamp when this security rule was deleted",
		},
		"created_by": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "User who created this security rule",
		},
		"updated_by": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "User who last updated this security rule",
		},
		"deleted_by": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "User who deleted this security rule",
		},
	}
}

func resourceIllumioSecurityRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pConfig, _ := m.(Config)
	illumioClient := pConfig.IllumioClient

	hrefRuleSet := d.Get("rule_set_href").(string)

	secRule, diags := expandIllumioSecurityRule(d)

	if diags.HasError() {
		return *diags
	}

	_, data, err := illumioClient.Create(fmt.Sprintf("%s/sec_rules", hrefRuleSet), secRule)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	pConfig.StoreHref("rule_sets", hrefRuleSet)

	d.SetId(data.S("href").Data().(string))

	return resourceIllumioSecurityRuleRead(ctx, d, m)
}

func expandIllumioSecurityRule(d *schema.ResourceData) (*models.SecurityRule, *diag.Diagnostics) {
	var diags diag.Diagnostics
	secRule := &models.SecurityRule{
		Enabled:               d.Get("enabled").(bool),
		Description:           d.Get("description").(string),
		ExternalDataSet:       d.Get("external_data_set").(string),
		ExternalDataReference: d.Get("external_data_reference").(string),
		SecConnect:            d.Get("sec_connect").(bool),
		Stateless:             d.Get("stateless").(bool),
		MachineAuth:           d.Get("machine_auth").(bool),
		UnscopedConsumers:     d.Get("unscoped_consumers").(bool),
	}

	if secRule.HasConflicts() {
		diags = append(diags, diag.Errorf("[illumio-core_security_rule] Only one of [\"sec_connect\", \"machine_auth\", \"stateless\"] can be set to true")...)
	}

	resLabelAs, errs := expandIllumioSecurityRuleResolveLabelsAs(d.Get("resolve_labels_as").([]interface{})[0])
	diags = append(diags, errs...)
	secRule.ResolveLabelsAs = resLabelAs

	ingServs, errs := expandIllumioSecurityRuleIngressService(
		d.Get("ingress_services").(*schema.Set).List(),
		secRule.ResolveLabelsAs.ProviderIsVirtualService(),
	)
	diags = append(diags, errs...)
	secRule.IngressServices = ingServs

	povs, errs := expandIllumioSecurityRuleProviders(d.Get("providers").(*schema.Set).List())
	diags = append(diags, errs...)
	secRule.Providers = povs

	cons, errs := expandIllumioSecurityRuleConsumers(d.Get("consumers").(*schema.Set).List())
	diags = append(diags, errs...)
	secRule.Consumers = cons

	return secRule, &diags
}

func expandIllumioSecurityRuleResolveLabelsAs(o interface{}) (*models.SecurityRuleResolveLabelAs, diag.Diagnostics) {
	var diags diag.Diagnostics
	resLabelsAs := o.(map[string]interface{})

	rProvs := getStringList(resLabelsAs["providers"].(*schema.Set).List())
	rCons := getStringList(resLabelsAs["consumers"].(*schema.Set).List())

	if !validateList(rProvs, validSRResLabelAsValues) {
		diags = append(diags, diag.Errorf(`[illumio-core_security_rule] Invalid value for resolve_value_as.providers, allowed values are "workloads" and "virtual_services"`)...)
	}

	if !validateList(rCons, validSRResLabelAsValues) {
		diags = append(diags, diag.Errorf(`[illumio-core_security_rule] Invalid value for resolve_value_as.consumers, allowed values are "workloads" and "virtual_services"`)...)
	}

	v := &models.SecurityRuleResolveLabelAs{
		Providers: rProvs,
		Consumers: rCons,
	}

	return v, diags
}

func expandIllumioSecurityRuleIngressService(inServices []interface{}, setEmpty bool) ([]map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	iss := []map[string]interface{}{}

	// Throw error if virtual_services is the only value set in resolve_label_as.provider and ingress_service's resource is non empty
	if setEmpty && len(inServices) > 0 {
		diags = append(diags, diag.Errorf("[illumio-core_security_rule] If the only value in the providers of resolve_label_as block is \"virtual_services\", then setting ingress_services is not allowed")...)
	}

	if !setEmpty {
		if len(inServices) == 0 {
			diags = append(diags, diag.Errorf("[illumio-core_security_rule] At least one ingress_service must be specified if providers of resolve_label_as block has \"workloads\"")...)
		}
		for _, service := range inServices {
			s := service.(map[string]interface{})

			m := make(map[string]interface{})

			if isIngressServiceSchemaValid(s, &diags) {
				if s["href"].(string) != "" {
					m["href"] = s["href"].(string)
				}

				if v, ok := getInt(s["proto"]); ok {
					m["proto"] = v
					if vPort, ok := getInt(s["port"]); ok {
						m["port"] = vPort
						if vToPort, ok := getInt(s["to_port"]); ok {
							if vToPort <= vPort {
								diags = append(diags, diag.Errorf(" [illumio-core_security_rule] Value of to_port can't be less or equal to value of port inside ingress_services")...)
							} else {
								m["to_port"] = vToPort
							}
						}
					}
				}
			}

			iss = append(iss, m)
		}
	}

	return iss, diags
}

// Validates schema of the security_rule.ingress_service parameter.
//
// Verifes if required fileds are defined or not.
func isIngressServiceSchemaValid(s map[string]interface{}, diags *diag.Diagnostics) bool {
	hrefOk := s["href"].(string) != ""
	protoOk := s["proto"].(string) != ""
	portOk := s["port"].(string) != ""
	toPortOk := s["to_port"].(string) != ""

	switch {
	case !hrefOk && !protoOk:
		*diags = append(*diags, diag.Errorf("[illumio-core_security_rule] Atleast one of [href, proto] is required inside ingress_services")...)

	case hrefOk && protoOk:
		*diags = append(*diags, diag.Errorf("[illumio-core_security_rule] Exactly one of [href, proto] is allowed inside ingress_services")...)

	case hrefOk:
		if portOk || toPortOk { // If port or to_port are defined with href, return error
			*diags = append(*diags, diag.Errorf("[illumio-core_security_rule] port/proto is not allowed with href inside ingress_services")...)
			return false
		}
		return true

	case protoOk:
		if !portOk && toPortOk { // If to_port is defined without port, return error
			*diags = append(*diags, diag.Errorf("[illumio-core_security_rule] port is required with to_port inside ingress_services")...)
			return false
		}
		return true
	}

	return false
}

func expandIllumioSecurityRuleProviders(providers []interface{}) ([]*models.SecurityRuleProvider, diag.Diagnostics) {
	provs := []*models.SecurityRuleProvider{}

	for _, provider := range providers {
		p := provider.(map[string]interface{})
		prov := &models.SecurityRuleProvider{
			Actors:         p["actors"].(string),
			Label:          getHrefObj(p["label"]),
			LabelGroup:     getHrefObj(p["label_group"]),
			Workload:       getHrefObj(p["workload"]),
			VirtualService: getHrefObj(p["virtual_service"]),
			VirtualServer:  getHrefObj(p["virtual_server"]),
			IPList:         getHrefObj(p["ip_list"]),
		}
		if !prov.HasOneActor() {
			return nil, diag.Errorf("[illumio-core_security_rule] Provider block can have only one rule actor")
		}

		provs = append(provs, prov)
	}
	return provs, diag.Diagnostics{}
}

func expandIllumioSecurityRuleConsumers(consumers []interface{}) ([]*models.SecurityRuleConsumer, diag.Diagnostics) {
	cons := []*models.SecurityRuleConsumer{}

	for _, consumer := range consumers {
		p := consumer.(map[string]interface{})

		con := &models.SecurityRuleConsumer{
			Actors:         p["actors"].(string),
			Label:          getHrefObj(p["label"]),
			LabelGroup:     getHrefObj(p["label_group"]),
			Workload:       getHrefObj(p["workload"]),
			VirtualService: getHrefObj(p["virtual_service"]),
			IPList:         getHrefObj(p["ip_list"]),
		}

		if !con.HasOneActor() {
			return nil, diag.Errorf("[illumio-core_security_rule] Consumer block can have only one rule actor")
		}
		cons = append(cons, con)
	}

	return cons, diag.Diagnostics{}
}

func resourceIllumioSecurityRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	pConfig, _ := m.(Config)
	illumioClient := pConfig.IllumioClient

	_, data, err := illumioClient.Get(d.Id(), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// extract the parent HREF and set it
	href := data.S("href").Data().(string)
	d.Set("rule_set_href", getParentHref(href))

	srKeys := []string{
		"href",
		"enabled",
		"description",
		"external_data_set",
		"external_data_reference",
		"sec_connect",
		"stateless",
		"machine_auth",
		"unscoped_consumers",
		"update_type",
		"created_at",
		"updated_at",
		"deleted_at",
		"created_by",
		"updated_by",
		"deleted_by",
	}

	for _, key := range srKeys {
		if data.Exists(key) {
			d.Set(key, data.S(key).Data())
		} else {
			d.Set(key, nil)
		}
	}

	isKey := "ingress_services"
	if data.Exists(isKey) {
		d.Set(isKey, extractResourceSecurityRuleIngressService(data.S(isKey)))
	} else {
		d.Set(isKey, nil)
	}

	rlaKey := "resolve_labels_as"
	if data.Exists(rlaKey) {
		d.Set(rlaKey, extractSecurityRuleResolveLabelAs(data.S(rlaKey)))
	} else {
		d.Set(rlaKey, nil)
	}

	prkey := "providers"
	if data.Exists(prkey) {
		d.Set(prkey, extractResourceRuleActors(data.S(prkey)))
	}

	cnKeys := "consumers"
	if data.Exists(cnKeys) {
		d.Set(cnKeys, extractResourceRuleActors(data.S(cnKeys)))
	}

	return diagnostics
}

func extractSecurityRuleResolveLabelAs(data *gabs.Container) []interface{} {
	m := map[string][]interface{}{
		"providers": data.S("providers").Data().([]interface{}),
		"consumers": data.S("consumers").Data().([]interface{}),
	}

	return []interface{}{m}
}

func extractResourceSecurityRuleIngressService(data *gabs.Container) []map[string]interface{} {
	isKeys := []string{
		"proto",
		"port",
		"to_port",
	}

	iss := []map[string]interface{}{}
	for _, ingSerData := range data.Children() {
		is := map[string]interface{}{}

		for k, v := range ingSerData.ChildrenMap() {
			if k == "href" {
				is[k] = v.Data().(string)
			} else if contains(isKeys, k) {
				is[k] = strconv.Itoa(int(v.Data().(float64)))
			}

			iss = append(iss, is)
		}
	}

	return iss
}

func resourceIllumioSecurityRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pConfig, _ := m.(Config)
	illumioClient := pConfig.IllumioClient

	var diags diag.Diagnostics

	resLabelAs, errs := expandIllumioSecurityRuleResolveLabelsAs(d.Get("resolve_labels_as").([]interface{})[0])
	diags = append(diags, errs...)

	ingServs, errs := expandIllumioSecurityRuleIngressService(
		d.Get("ingress_services").(*schema.Set).List(),
		resLabelAs.ProviderIsVirtualService(),
	)
	diags = append(diags, errs...)

	povs, errs := expandIllumioSecurityRuleProviders(d.Get("providers").(*schema.Set).List())
	diags = append(diags, errs...)

	cons, errs := expandIllumioSecurityRuleConsumers(d.Get("consumers").(*schema.Set).List())
	diags = append(diags, errs...)

	secRule := &models.SecurityRule{
		Enabled:               d.Get("enabled").(bool),
		Description:           d.Get("description").(string),
		ExternalDataSet:       d.Get("external_data_set").(string),
		ExternalDataReference: d.Get("external_data_reference").(string),
		SecConnect:            d.Get("sec_connect").(bool),
		Stateless:             d.Get("stateless").(bool),
		MachineAuth:           d.Get("machine_auth").(bool),
		UnscopedConsumers:     d.Get("unscoped_consumers").(bool),
		ResolveLabelsAs:       resLabelAs,
		IngressServices:       ingServs,
		Providers:             povs,
		Consumers:             cons,
	}

	if secRule.HasConflicts() {
		diags = append(diags, diag.Errorf("[illumio-core_security_rule] Exactly one of [\"sec_connect\", \"machine_auth\", \"stateless\"] can be set to true")...)
	}

	if diags.HasError() {
		return diags
	}

	_, err := illumioClient.Update(d.Id(), secRule)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleSetHREF := d.Get("rule_set_href").(string)

	pConfig.StoreHref("rule_sets", ruleSetHREF)

	return resourceIllumioSecurityRuleRead(ctx, d, m)
}

func resourceIllumioSecurityRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	pConfig, _ := m.(Config)
	illumioClient := pConfig.IllumioClient

	ruleSetHREF := d.Get("rule_set_href").(string)

	_, err := illumioClient.Delete(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	pConfig.StoreHref("rule_sets", ruleSetHREF)

	d.SetId("")
	return diagnostics
}
