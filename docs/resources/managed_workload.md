---
layout: "illumio-core"
page_title: "illumio-core_workload Resource - terraform-provider-illumio-core"
sidebar_current: "docs-illumio-core-resoruce-workload"
subcategory: ""
description: |-
  Manages Illumio Managed Workloads
---

# illumio-core_workload (Resource)

The `illumio-core_managed_workload` resource can be used to read, update, and unpair managed Workloads with paired VENs in the Illumio Policy Compute Engine.  

~> managed workload resources cannot be created using Terraform. Managed workloads must be imported after a VEN is installed and paired with the PCE.  

## Importing  

After pairing a VEN with the PCE, the managed workload object can be imported with the following command:  

```sh
$ terraform import illumio-core_managed_workload.example ""
```

After import, configuration changes can be planned and applied as normal.  

Example Usage
------------

```hcl
resource "illumio-core_managed_workload" "aws_host1" {
    name = "aws.sc.db0"
    description = "Shopping Cart app database 01"

    external_data_set = "S3"
    external_data_reference = "s3://hosts/shopping_cart"

    service_principal_name = "Kerberos SPN"
    service_provider = "aws"

    data_center = "North Virginia"
    data_center_zone = "us-east-1"

    enforcement_mode = "visibility_only"

    labels {
      href = illumio-core_label.role_db.href
    }

    labels {
      href = illumio-core_label.app_shopping_cart.href
    }

    labels {
      href = illumio-core_label.env_dev.href
    }

    labels {
      href = illumio-core_label.loc_aws.href
    }
}  
```

## Schema

### Optional

- `name` (String) Name of the Workload. The name should be up to 255 characters
- `data_center` (String) Data center for Workload. The data_center should be up to 255 characters
- `data_center_zone` (String) Data center Zone for Workload. The data_center_zone should be up to 255 characters
- `description` (String) The long description of the workload
- `enforcement_mode` (String) Enforcement mode of workload(s) to return. Allowed values for enforcement modes are "idle","visibility_only","full", and "selective". Default value: "visibility_only"
- `external_data_reference` (String) Unique identifier for the workload in the external data source
- `external_data_set` (String) The data source from which a resource originates
- `ignored_interface_names` (List of String) Workload interface names to ignore (e.g. `eth0`). Ignored interfaces will not be included in policy configuration provided by the PCE.
- `labels` (Block Set) Assigned labels for workload (see [below for nested schema](#nestedblock--labels))
- `service_principal_name` (String) The Kerberos Service Principal Name (SPN). The SPN should be between 1 to 255 characters
- `service_provider` (String) Service provider for Workload. The service_provider should be up to 255 characters

### Read-Only

- `hostname` (String) The hostname of this workload. Set by the VEN.
- `public_ip` (String) The public IP address of the server. Set by the VEN.
- `online` (Boolean) True if the workload is online, otherwise false. Set by the VEN.
- `agent_to_pce_certificate_authentication_id` (String) PKI Certificate identifier to be used by the PCE for authenticating the VEN.
- `blocked_connection_action` (String) Blocked Connection Action for Workload
- `caps` (List of String) CAPS for Workload
- `container_cluster` (List of Object) Container Cluster for Workload (see [below for nested schema](#nestedatt--container_cluster))
- `containers_inherit_host_policy` (Boolean) This workload will apply the policy it receives both to itself and the containers hosted by it
- `created_at` (String) Timestamp when this label group was first created
- `created_by` (Map of String) User who created this label group
- `deleted` (Boolean) This indicates that the workload has been deleted
- `deleted_at` (String) Timestamp when this label group was last deleted
- `deleted_by` (Map of String) User who deleted this label group
- `detected_vulnerabilities` (List of Object) Detected Vulnerabilities for Workload (see [below for nested schema](#nestedatt--detected_vulnerabilities))
- `distinguished_name` (String) X.509 Subject distinguished name.
- `firewall_coexistence` (List of Object) Firewall coexistence mode for Workload (see [below for nested schema](#nestedatt--firewall_coexistence))
- `href` (String) URI of the Workload
- `ike_authentication_certificate` (Map of String) IKE authentication certificate for certificate-based Secure Connect and Machine Auth
- `interfaces` (Block Set) Workload network interfaces (see [below for nested schema](#nestedblock--interfaces))
- `os_detail` (String) Additional OS details - just displayed to end-user. Set by the VEN.
- `os_id` (String) OS identifier for Workload. Set by the VEN.
- `selectively_enforced_services` (List of Object) Selectively Enforced Services for Workload (see [below for nested schema](#nestedatt--selectively_enforced_services))
- `services` (List of Object) Service report for Workload (see [below for nested schema](#nestedatt--services))
- `updated_at` (String) Timestamp when this label group was last updated
- `updated_by` (Map of String) User who last updated this label group
- `ven` (Map of String) VENS for Workload
- `visibility_level` (String) Visibility Level of workload(s) to return
- `vulnerabilities_summary` (List of Object) Vulnerabilities summary associated with the workload (see [below for nested schema](#nestedatt--vulnerabilities_summary))

<a id="nestedblock--labels"></a>
### Nested Schema for `labels`

Required:

- `href` (String) URI of label

<a id="nestedatt--container_cluster"></a>
### Nested Schema for `container_cluster`

Read-Only:

- `href` (String) URI of container cluster
- `name` (String) Name of container cluster

<a id="nestedatt--detected_vulnerabilities"></a>
### Nested Schema for `detected_vulnerabilities`

Read-Only:

- `ip_address` (String) The IP address of the host where the vulnerability is found
- `port` (Number) The port that is associated with the vulnerability
- `port_exposure` (Number) The exposure of the port based on the current policy
- `port_wide_exposure` (List of Object) High end of an IP range(see [below for nested schema](#nestedobjatt--detected_vulnerabilities--port_wide_exposure))
- `proto` (Number) The protocol that is associated with the vulnerability
- `vulnerability` (Set of Object) Vulnerability for Workload (see [below for nested schema](#nestedobjatt--detected_vulnerabilities--vulnerability))
- `vulnerability_report` (Set of Object) Vulnerability Report for Workload(see [below for nested schema](#nestedobjatt--detected_vulnerabilities--vulnerability_report))
- `workload` (List of Object) URI of Workload (see [below for nested schema](#nestedobjatt--detected_vulnerabilities--workload))

<a id="nestedblock--interfaces"></a>
### Nested Schema for `interfaces`

Read-only:

- `name` (String) Interface name. Can be up to 255 characters
- `address` (String) Interface IP address. Must be in IPv4 or IPv6 format
- `cidr_block` (Number) Interface CIDR block bits
- `default_gateway_address` (String) Interface Default Gateway IP address. Must be in IPv4 or IPv6 format
- `friendly_name` (String) User-friendly interface name. Can be up to 255 characters
- `link_state` (String) Interface link state. Allowed values are "up", "down", and "unknown"
- `loopback` (Boolean) Whether or not the interface represents a loopback address on the workload
- `network` (Map of String) Interface Network object HREFs
- `network_detection_mode` (String) Interface Network Detection Mode

<a id="nestedobjatt--detected_vulnerabilities--port_wide_exposure"></a>
### Nested Schema for `detected_vulnerabilities.port_wide_exposure`

Read-Only:

- `any` (Boolean) The boolean value representing if at least one port is exposed to internet (any rule) on the workload
- `ip_list` (Boolean) The boolean value representing if at least one port is exposed to ip_list(s) on the workload

<a id="nestedobjatt--detected_vulnerabilities--vulnerability"></a>
### Nested Schema for `detected_vulnerabilities.vulnerability`

Read-Only:

- `href` (String) The URI of the workload to which this vulnerability belongs to
- `name` (String) The title/name of the vulnerability
- `score` (Number) The normalized score of the vulnerability within the range of 0 to 100


<a id="nestedobjatt--detected_vulnerabilities--vulnerability_report"></a>
### Nested Schema for `detected_vulnerabilities.vulnerability_report`

Read-Only:

- `href` (String) The URI of the report to which this vulnerability belongs to

<a id="nestedobjatt--detected_vulnerabilities--workload"></a>
### Nested Schema for `detected_vulnerabilities.workload`

Read-Only:

- `href` (String) The URI of the workload to which this vulnerability belongs to

<a id="nestedatt--firewall_coexistence"></a>
### Nested Schema for `firewall_coexistence`

Read-Only:

- `illumio_primary` (Boolean) Illumio is the primary firewall if set to true

<a id="nestedatt--selectively_enforced_services"></a>
### Nested Schema for `selectively_enforced_services`

Read-Only:

- `href` (String) URI of Selectively Enforced Services
- `port` (Integer) Port number, or the starting port of a range. If unspecified, this will apply to all ports for the given protocol. Minimum and maximum value for port is 0 and 65535 respectively.
- `to_port` (Integer) Upper end of port range; this field should not be included if specifying an individual port. Minimum and maximum value for to_port is 0 and 65535 respectively.
- `proto` (Integer) Transport protocol of Selectively Enforced Services

<a id="nestedatt--services"></a>
### Nested Schema for `services`

Read-Only:

- `created_at` (String) Timestamp when this service was first created
- `open_service_ports` (List of Object) A list of open ports (see [below for nested schema] (#nestedobjatt--services--open_service_ports))
- `uptime_seconds` (Number) How long since the last reboot of this box - used as a timestamp for this

<a id="nestedobjatt--services--open_service_ports"></a>
### Nested Schema for `services.open_service_ports`

Read-Only:

- `address` (String) The local address this service is bound to
- `package` (String) The RPM/DEB package that the program is part of
- `port` (Number) The local port this service is bound to
- `process_name` (String) The process name (including the full path)
- `protocol` (Number) Transport protocol for open service ports
- `user` (String) The user account that the process is running under
- `win_service_name` (String) Name of the Windows service

<a id="nestedatt--vulnerabilities_summary"></a>
### Nested Schema for `vulnerabilities_summary`

Read-Only:

- `max_vulnerability_score` (Number) The maximum of all the vulnerability scores associated with the detected_vulnerabilities on the workload
- `num_vulnerabilities` (Number) Number of vulnerabilities associated with the workload
- `vulnerability_exposure_score` (Number) The aggregated vulnerability exposure score of the workload across all the vulnerable ports.
- `vulnerability_score` (Number) The aggregated vulnerability score of the workload across all the vulnerable ports.
- `vulnerable_port_exposure` (Number) The aggregated vulnerability port exposure score of the workload across all the vulnerable ports
- `vulnerable_port_wide_exposure` (List of Object) High end of an IP range (see [below for nested schema](#nestedobjatt--vulnerabilities_summary--vulnerable_port_wide_exposure))

<a id="nestedobjatt--vulnerabilities_summary--vulnerable_port_wide_exposure"></a>
### Nested Schema for `vulnerabilities_summary.vulnerable_port_wide_exposure`

Read-Only:

- `any` (Boolean) The boolean value representing if at least one port is exposed to internet (any rule) on the workload
- `ip_list` (Boolean) The boolean value representing if at least one port is exposed to ip_list(s) on the workload
