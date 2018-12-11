package main

import (
	"fmt"

	"github.com/atlassian/go-vtm"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePool() *schema.Resource {
	return &schema.Resource{
		Create: resourcePoolCreate,
		Read:   resourcePoolRead,
		Update: resourcePoolUpdate,
		Delete: resourcePoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"keepalive": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"bandwidth_class": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"connection_max_connect_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  4,
			},

			"connection_max_connections_per_node": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"connection_max_queue_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"connection_max_reply_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},

			"connection_queue_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10,
			},

			"dns_autoscale_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"dns_autoscale_hostnames": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"dns_autoscale_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  80,
			},

			"failure_pool": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"load_balancing_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "round_robin",
			},

			"load_balancing_priority_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"load_balancing_priority_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"max_connection_attempts": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"max_idle_connections_pernode": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  50,
			},

			"max_timed_out_connection_attempts": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},

			"monitors": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"node_close_with_rst": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"node_connection_attempts": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
			},

			"node_delete_behavior": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "immediate",
			},

			"node_drain_to_delete_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"node": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node": {
							Type:     schema.TypeString,
							Required: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"state": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "active",
						},
						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
					},
				},
			},
			"note": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"passive_monitoring": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"persistence_class": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"tcp_nagle": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"transparent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"ssl_encrypt": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"udp_accept_from": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "dest_only",
			},

			"udp_accept_from_mask": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func resourcePoolCreate(d *schema.ResourceData, meta interface{}) error {
	err := resourcePoolSet(d, meta)
	if err != nil {
		return err
	}

	return resourcePoolRead(d, meta)
}

func resourcePoolRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerConfig).client

	r, resp, err := c.GetPool(d.Get("name").(string))
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			// The resource doesn't exist anymore
			d.SetId("")

			return nil
		}

		return fmt.Errorf("Error reading resource: %s", err)
	}

	d.Set("keepalive", bool(*r.HTTP.Keepalive))
	d.Set("bandwidth_class", string(*r.Basic.BandwidthClass))
	d.Set("connection_max_connect_time", int(*r.Connection.MaxConnectTime))
	d.Set("connection_max_connections_per_node", int(*r.Connection.MaxConnectionsPerNode))
	d.Set("connection_max_queue_size", int(*r.Connection.MaxQueueSize))
	d.Set("connection_max_reply_time", int(*r.Connection.MaxReplyTime))
	d.Set("connection_queue_timeout", int(*r.Connection.QueueTimeout))

	d.Set("dns_autoscale_enabled", bool(*r.DNSAutoscale.Enabled))
	d.Set("dns_autoscale_hostnames", []string(*r.DNSAutoscale.Hostnames))
	d.Set("dns_autoscale_port", int(*r.DNSAutoscale.Port))

	d.Set("failure_pool", string(*r.Basic.FailurePool))
	d.Set("load_balancing_algorithm", string(*r.LoadBalancing.Algorithm))
	d.Set("load_balancing_priority_enabled", bool(*r.LoadBalancing.PriorityEnabled))
	d.Set("load_balancing_priority_nodes", int(*r.LoadBalancing.PriorityNodes))
	d.Set("max_connection_attempts", int(*r.Basic.MaxConnectionAttempts))
	d.Set("max_idle_connections_pernode", int(*r.Basic.MaxIdleConnectionsPerNode))
	d.Set("max_timed_out_connection_attempts", int(*r.Basic.MaxTimedOutConnectionAttempts))
	d.Set("monitors", []string(*r.Basic.Monitors))
	d.Set("node_close_with_rst", bool(*r.Basic.NodeCloseWithRST))
	d.Set("node_connection_attempts", int(*r.Basic.NodeConnectionAttempts))
	d.Set("node_delete_behavior", string(*r.Basic.NodeDeleteBehavior))
	d.Set("node_drain_to_delete_timeout", int(*r.Basic.NodeDrainToDeleteTimeout))
	d.Set("note", string(*r.Basic.Note))
	d.Set("passive_monitoring", bool(*r.Basic.PassiveMonitoring))
	d.Set("persistence_class", string(*r.Basic.PersistenceClass))
	d.Set("tcp_nagle", bool(*r.TCP.Nagle))
	d.Set("udp_accept_from", string(*r.UDP.AcceptFrom))
	d.Set("udp_accept_from_mask_class", string(*r.UDP.AcceptFromMask))
	d.Set("transparent", bool(*r.Basic.Transparent))
	d.Set("ssl_encrypt", bool(*r.SSL.Enable))

	nodesList := make([]map[string]interface{}, 0, len(*r.Basic.NodesTable))

	if r.DNSAutoscale.Enabled != nil && !*r.DNSAutoscale.Enabled {
		for _, node := range *r.Basic.NodesTable {
			nodeTerraform := make(map[string]interface{})
			if node.Node != nil {
				nodeTerraform["node"] = string(*node.Node)
			}
			if node.Weight != nil {
				nodeTerraform["weight"] = int(*node.Weight)
			}
			if node.State != nil {
				nodeTerraform["state"] = string(*node.State)
			}
			if node.Priority != nil {
				nodeTerraform["priority"] = int(*node.Priority)
			}
			nodesList = append(nodesList, nodeTerraform)
		}
	}

	d.Set("node", nodesList)

	return nil
}

func resourcePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	err := resourcePoolSet(d, meta)
	if err != nil {
		return err
	}

	return resourcePoolRead(d, meta)
}

func resourcePoolDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerConfig).client
	r := stingray.NewPool(d.Id())

	_, err := c.Delete(r)
	return err
}

func resourcePoolSet(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerConfig).client
	r := stingray.NewPool(d.Get("name").(string))

	setBool(&r.HTTP.Keepalive, d, "keepalive")
	setString(&r.Basic.BandwidthClass, d, "bandwidth_class")
	setInt(&r.Connection.MaxConnectTime, d, "connection_max_connect_time")
	setInt(&r.Connection.MaxConnectionsPerNode, d, "connection_max_connections_per_node")
	setInt(&r.Connection.MaxQueueSize, d, "connection_max_queue_size")
	setInt(&r.Connection.MaxReplyTime, d, "connection_max_reply_time")
	setInt(&r.Connection.QueueTimeout, d, "connection_queue_timeout")
	setBool(&r.DNSAutoscale.Enabled, d, "dns_autoscale_enabled")
	setStringSet(&r.DNSAutoscale.Hostnames, d, "dns_autoscale_hostnames")
	setInt(&r.DNSAutoscale.Port, d, "dns_autoscale_port")
	setString(&r.Basic.FailurePool, d, "failure_pool")
	setString(&r.LoadBalancing.Algorithm, d, "load_balancing_algorithm")
	setBool(&r.LoadBalancing.PriorityEnabled, d, "load_balancing_priority_enabled")
	setInt(&r.LoadBalancing.PriorityNodes, d, "load_balancing_priority_nodes")
	setInt(&r.Basic.MaxConnectionAttempts, d, "max_connection_attempts")
	setInt(&r.Basic.MaxIdleConnectionsPerNode, d, "max_idle_connections_pernode")
	setInt(&r.Basic.MaxTimedOutConnectionAttempts, d, "max_timed_out_connection_attempts")
	setStringSet(&r.Basic.Monitors, d, "monitors")
	setBool(&r.Basic.NodeCloseWithRST, d, "node_close_with_rst")
	setInt(&r.Basic.NodeConnectionAttempts, d, "node_connection_attempts")
	setString(&r.Basic.NodeDeleteBehavior, d, "node_delete_behavior")
	setInt(&r.Basic.NodeDrainToDeleteTimeout, d, "node_drain_to_delete_timeout")
	setString(&r.Basic.Note, d, "note")
	setBool(&r.Basic.PassiveMonitoring, d, "passive_monitoring")
	setString(&r.Basic.PersistenceClass, d, "persistence_class")
	setBool(&r.TCP.Nagle, d, "tcp_nagle")
	setBool(&r.Basic.Transparent, d, "transparent")
	setBool(&r.SSL.Enable, d, "ssl_encrypt")
	setString(&r.UDP.AcceptFrom, d, "udp_accept_from")
	setString(&r.UDP.AcceptFromMask, d, "udp_accept_from_mask")

	*r.Basic.NodesTable = make([]stingray.Node, 0, 0)

	if v, ok := d.GetOk("node"); ok && r.DNSAutoscale.Enabled != nil && !*r.DNSAutoscale.Enabled {
		nodeList := v.(*schema.Set).List()

		for _, node := range nodeList {
			terraformNode := node.(map[string]interface{})

			VtmNode := &stingray.Node{}
			VtmNode.Node = stingray.String(terraformNode["node"].(string))
			if terraformNode["weight"].(int) != 0 {
				VtmNode.Weight = stingray.Int(terraformNode["weight"].(int))
			}
			VtmNode.Priority = stingray.Int(terraformNode["priority"].(int))
			VtmNode.State = stingray.String(terraformNode["state"].(string))

			*r.Basic.NodesTable = append(*r.Basic.NodesTable, *VtmNode)
		}
	}

	_, err := c.Set(r)
	if err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))

	return nil
}
