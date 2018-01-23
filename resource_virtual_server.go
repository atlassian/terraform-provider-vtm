package main

import (
	"fmt"

	"github.com/atlassian/go-vtm"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVirtualServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceVirtualServerCreate,
		Read:   resourceVirtualServerRead,
		Update: resourceVirtualServerUpdate,
		Delete: resourceVirtualServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"add_x_forwarded_for": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"add_x_forwarded_proto": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},

			"connection_errors_error_file": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Default",
			},

			"connection_keepalive_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10,
			},

			"connection_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},

			"connect_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10,
			},

			// NOTE: Default for the Stingray API is false
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"gzip_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"gzip_compress_level": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"gzip_max_size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10000000,
			},
			"gzip_min_size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1000,
			},

			"gzip_include_mime": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"http_location_rewrite": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "if_host_matches",
			},

			// NOTE: Default for the Stingray API is true
			"listen_on_any": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"listen_on_traffic_ips": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"log_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"log_filename": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "%zeushome%/zxtm/log/%v.log",
			},

			"log_format": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "%h %l %u %t \"%r\" %s %b \"%{Referer}i\" \"%{User-agent}i\"",
			},

			"log_server_connection_failures": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"note": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"pool": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"protocol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "http",
			},

			"recent_connections_save_all": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"protection_class": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"request_rules": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"response_rules": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"completion_rules": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"ssl_add_http_headers": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ssl_trust_magic": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"ssl_decrypt": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"ssl_server_cert_default": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"ssl_support_ssl3": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "use_default",
			},

			"ssl_support_tls1": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "use_default",
			},

			"ssl_support_tls1_1": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "use_default",
			},

			"ssl_support_tls1_2": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "use_default",
			},

			"ssl_server_cert_host_mapping": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"host": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set: hashServerCertHostMapping,
			},

			"syslog_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"syslog_format": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "%h %l %u %t \"%r\" %s %b \"%{Referer}i\" \"%{User-agent}i\"",
			},
			"syslog_ip_end_point": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"syslog_msg_len_limit": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1024,
			},

			"web_cache_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"web_cache_max_time": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  600,
			},
			"web_cache_refresh_time": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},
			"web_cache_control_out": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"web_cache_error_page_time": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
			"max_client_buffer": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  65536,
			},
			"max_server_buffer": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  65536,
			},
			"ssl_ciphers": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func resourceVirtualServerCreate(d *schema.ResourceData, meta interface{}) error {
	err := resourceVirtualServerSet(d, meta)
	if err != nil {
		return err
	}

	return resourceVirtualServerRead(d, meta)
}

func resourceVirtualServerRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerConfig).client

	r, resp, err := c.GetVirtualServer(d.Get("name").(string))
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			// The resource doesn't exist anymore
			d.SetId("")

			return nil
		}

		return fmt.Errorf("Error reading resource: %s", err)
	}

	d.Set("connection_errors_error_file", string(*r.ConnectionErrors.ErrorFile))
	d.Set("connection_keepalive_timeout", int(*r.Connection.KeepaliveTimeout))
	d.Set("connection_timeout", int(*r.Connection.Timeout))
	d.Set("protection_class", string(*r.Basic.ProtectionClass))
	d.Set("add_x_forwarded_for", bool(*r.Basic.AddXForwardedFor))
	d.Set("add_x_forwarded_proto", bool(*r.Basic.AddXForwardedProto))
	d.Set("connect_timeout", int(*r.Basic.ConnectTimeout))
	d.Set("enabled", bool(*r.Basic.Enabled))
	d.Set("gzip_enabled", bool(*r.Gzip.Enabled))
	d.Set("gzip_compress_level", int(*r.Gzip.CompressLevel))
	d.Set("gzip_max_size", int(*r.Gzip.MaxSize))
	d.Set("gzip_min_size", int(*r.Gzip.MinSize))
	d.Set("gzip_include_mime", []string(*r.Gzip.IncludeMIME))
	d.Set("http_location_rewrite", string(*r.HTTP.LocationRewrite))
	d.Set("listen_on_any", bool(*r.Basic.ListenOnAny))
	d.Set("listen_on_traffic_ips", []string(*r.Basic.ListenOnTrafficIPs))
	d.Set("log_enabled", bool(*r.Log.Enabled))
	d.Set("log_filename", string(*r.Log.Filename))
	d.Set("log_format", string(*r.Log.Format))
	d.Set("log_server_connection_failures", bool(*r.Log.ServerConnectionFailures))
	d.Set("note", string(*r.Basic.Note))
	d.Set("pool", string(*r.Basic.Pool))
	d.Set("port", int(*r.Basic.Port))
	d.Set("protocol", string(*r.Basic.Protocol))
	d.Set("recent_connections_save_all", bool(*r.RecentConnections.SaveAll))
	d.Set("request_rules", []string(*r.Basic.RequestRules))
	d.Set("response_rules", []string(*r.Basic.ResponseRules))
	d.Set("completion_rules", []string(*r.Basic.CompletionRules))
	d.Set("ssl_add_http_headers", bool(*r.SSL.AddHTTPHeaders))
	d.Set("ssl_trust_magic", bool(*r.SSL.TrustMagic))
	d.Set("ssl_decrypt", bool(*r.Basic.SSLDecrypt))
	d.Set("ssl_server_cert_default", string(*r.SSL.ServerCertDefault))
	d.Set("ssl_server_cert_host_mapping", flattenServerCertHostMappingTable(*r.SSL.ServerCertHostMapping))
	d.Set("syslog_enabled", bool(*r.Syslog.Enabled))
	d.Set("syslog_format", string(*r.Syslog.Format))
	d.Set("syslog_ip_end_point", string(*r.Syslog.IPEndPoint))
	d.Set("syslog_msg_len_limit", int(*r.Syslog.MsgLenLimit))
	d.Set("web_cache_enabled", bool(*r.WebCache.Enabled))
	d.Set("web_cache_max_time", int(*r.WebCache.MaxTime))
	d.Set("web_cache_refresh_time", int(*r.WebCache.RefreshTime))
	d.Set("web_cache_control_out", string(*r.WebCache.ControlOut))
	d.Set("web_cache_error_page_time", int(*r.WebCache.ErrorPageTime))
	d.Set("max_client_buffer", int(*r.Connection.MaxClientBuffer))
	d.Set("max_server_buffer", int(*r.Connection.MaxServerBuffer))
	d.Set("ssl_ciphers", string(*r.SSL.SSLCiphers))
	d.Set("ssl_support_ssl3", string(*r.SSL.SSLSupportSSL3))
	d.Set("ssl_support_tls1", string(*r.SSL.SSLSupportTLS1))
	d.Set("ssl_support_tls1_1", string(*r.SSL.SSLSupportTLS11))
	d.Set("ssl_support_tls1_2", string(*r.SSL.SSLSupportTLS12))

	return nil
}

func resourceVirtualServerUpdate(d *schema.ResourceData, meta interface{}) error {
	err := resourceVirtualServerSet(d, meta)
	if err != nil {
		return err
	}

	return resourceVirtualServerRead(d, meta)
}

func resourceVirtualServerDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerConfig).client
	r := stingray.NewVirtualServer(d.Id())

	_, err := c.Delete(r)
	if err != nil {
		return err
	}

	return nil
}

func resourceVirtualServerSet(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerConfig).client
	r := stingray.NewVirtualServer(d.Get("name").(string))

	setString(&r.ConnectionErrors.ErrorFile, d, "connection_errors_error_file")
	setInt(&r.Connection.KeepaliveTimeout, d, "connection_keepalive_timeout")
	setInt(&r.Connection.Timeout, d, "connection_timeout")
	setInt(&r.Basic.ConnectTimeout, d, "connect_timeout")
	setBool(&r.Basic.Enabled, d, "enabled")
	setString(&r.Basic.ProtectionClass, d, "protection_class")
	setBool(&r.Basic.AddXForwardedFor, d, "add_x_forwarded_for")
	setBool(&r.Basic.AddXForwardedProto, d, "add_x_forwarded_proto")
	// NOTE: Set default for gzip_include_mime
	//
	// Default does not work for sets (only for primitive types),
	// so we must use GetOK. This means we will get the default
	// value if the parameter is unset OR set to the empty value.
	if _, ok := d.GetOk("gzip_include_mime"); ok {
		setStringSet(&r.Gzip.IncludeMIME, d, "gzip_include_mime")
	} else {
		r.Gzip.IncludeMIME = &[]string{"text/html", "text/plain"}
	}
	setBool(&r.Gzip.Enabled, d, "gzip_enabled")
	setInt(&r.Gzip.CompressLevel, d, "gzip_compress_level")
	setInt(&r.Gzip.MaxSize, d, "gzip_max_size")
	setInt(&r.Gzip.MinSize, d, "gzip_min_size")
	setString(&r.HTTP.LocationRewrite, d, "http_location_rewrite")
	setBool(&r.Basic.ListenOnAny, d, "listen_on_any")
	setStringSet(&r.Basic.ListenOnTrafficIPs, d, "listen_on_traffic_ips")
	setBool(&r.Log.Enabled, d, "log_enabled")
	setString(&r.Log.Filename, d, "log_filename")
	setString(&r.Log.Format, d, "log_format")
	setBool(&r.Log.ServerConnectionFailures, d, "log_server_connection_failures")
	setString(&r.Basic.Note, d, "note")
	setString(&r.Basic.Pool, d, "pool")
	setInt(&r.Basic.Port, d, "port")
	setString(&r.Basic.Protocol, d, "protocol")
	setBool(&r.RecentConnections.SaveAll, d, "recent_connections_save_all")
	setStringList(&r.Basic.RequestRules, d, "request_rules")
	setStringList(&r.Basic.ResponseRules, d, "response_rules")
	setStringList(&r.Basic.CompletionRules, d, "completion_rules")
	setBool(&r.SSL.AddHTTPHeaders, d, "ssl_add_http_headers")
	setBool(&r.SSL.TrustMagic, d, "ssl_trust_magic")
	setBool(&r.Basic.SSLDecrypt, d, "ssl_decrypt")
	setString(&r.SSL.ServerCertDefault, d, "ssl_server_cert_default")
	setServerCertHostMappingTable(&r.SSL.ServerCertHostMapping, d, "ssl_server_cert_host_mapping")
	setBool(&r.Syslog.Enabled, d, "syslog_enabled")
	setString(&r.Syslog.Format, d, "syslog_format")
	setString(&r.Syslog.IPEndPoint, d, "syslog_ip_end_point")
	setInt(&r.Syslog.MsgLenLimit, d, "syslog_msg_len_limit")
	setBool(&r.WebCache.Enabled, d, "web_cache_enabled")
	setInt(&r.WebCache.MaxTime, d, "web_cache_max_time")
	setInt(&r.WebCache.ErrorPageTime, d, "web_cache_error_page_time")
	setInt(&r.WebCache.RefreshTime, d, "web_cache_refresh_time")
	setString(&r.WebCache.ControlOut, d, "web_cache_control_out")
	setInt(&r.Connection.MaxClientBuffer, d, "max_client_buffer")
	setInt(&r.Connection.MaxServerBuffer, d, "max_server_buffer")
	setString(&r.SSL.SSLCiphers, d, "ssl_ciphers")
	setString(&r.SSL.SSLSupportSSL3, d, "ssl_support_ssl3")
	setString(&r.SSL.SSLSupportTLS1, d, "ssl_support_tls1")
	setString(&r.SSL.SSLSupportTLS11, d, "ssl_support_tls1_1")
	setString(&r.SSL.SSLSupportTLS12, d, "ssl_support_tls1_2")

	_, err := c.Set(r)
	if err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))

	return nil
}

func setServerCertHostMappingTable(target **stingray.ServerCertHostMappingTable, d *schema.ResourceData, key string) {
	if _, ok := d.GetOk(key); ok {
		table := d.Get(key).(*schema.Set).List()
		*target, _ = expandServerCertHostMappingTable(table)
	}
}

func expandServerCertHostMappingTable(configured []interface{}) (*stingray.ServerCertHostMappingTable, error) {
	table := make(stingray.ServerCertHostMappingTable, 0, len(configured))

	for _, raw := range configured {
		data := raw.(map[string]interface{})

		s := stingray.ServerCertHostMapping{
			Certificate: stingray.String(data["certificate"].(string)),
			Host:        stingray.String(data["host"].(string)),
		}

		table = append(table, s)
	}

	return &table, nil
}

func flattenServerCertHostMappingTable(list stingray.ServerCertHostMappingTable) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))

	for _, i := range list {
		s := map[string]interface{}{
			"certificate": *i.Certificate,
			"host":        *i.Host,
		}
		result = append(result, s)
	}

	return result
}

func hashServerCertHostMapping(v interface{}) int {
	m := v.(map[string]interface{})
	return hashcode.String(m["host"].(string))
}
