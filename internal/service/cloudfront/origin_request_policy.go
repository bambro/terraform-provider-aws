package cloudfront

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func ResourceOriginRequestPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceOriginRequestPolicyCreate,
		Read:   resourceOriginRequestPolicyRead,
		Update: resourceOriginRequestPolicyUpdate,
		Delete: resourceOriginRequestPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cookies_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cookie_behavior": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"none", "whitelist", "all"}, false),
						},
						"cookies": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"items": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"etag": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"headers_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"header_behavior": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"none", "whitelist", "allViewer", "allViewerAndWhitelistCloudFront"}, false),
						},
						"headers": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"items": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query_strings_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query_string_behavior": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"none", "whitelist", "allExcept", "all"}, false),
						},
						"query_strings": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"items": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceOriginRequestPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).CloudFrontConn

	request := &cloudfront.CreateOriginRequestPolicyInput{
		OriginRequestPolicyConfig: expandCloudFrontOriginRequestPolicyConfig(d),
	}

	resp, err := conn.CreateOriginRequestPolicy(request)

	if err != nil {
		return err
	}

	d.SetId(aws.StringValue(resp.OriginRequestPolicy.Id))

	return resourceOriginRequestPolicyRead(d, meta)
}

func resourceOriginRequestPolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).CloudFrontConn
	request := &cloudfront.GetOriginRequestPolicyInput{
		Id: aws.String(d.Id()),
	}

	resp, err := conn.GetOriginRequestPolicy(request)
	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, "ResourceNotFoundException") {
		log.Printf("[WARN] CloudFront Origin Request Policy (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading CloudFront Origin Request Policy (%s): %w", d.Id(), err)
	}

	if resp == nil || resp.OriginRequestPolicy == nil || resp.OriginRequestPolicy.OriginRequestPolicyConfig == nil {
		return fmt.Errorf("error reading CloudFront Origin Request Policy (%s): empty response", d.Id())
	}

	d.Set("etag", resp.ETag)

	originRequestPolicy := resp.OriginRequestPolicy.OriginRequestPolicyConfig
	d.Set("comment", originRequestPolicy.Comment)
	d.Set("name", originRequestPolicy.Name)
	d.Set("cookies_config", flattenCloudFrontOriginRequestPolicyCookiesConfig(originRequestPolicy.CookiesConfig))
	d.Set("headers_config", flattenCloudFrontOriginRequestPolicyHeadersConfig(originRequestPolicy.HeadersConfig))
	d.Set("query_strings_config", flattenCloudFrontOriginRequestPolicyQueryStringsConfig(originRequestPolicy.QueryStringsConfig))

	return nil
}

func resourceOriginRequestPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).CloudFrontConn

	request := &cloudfront.UpdateOriginRequestPolicyInput{
		OriginRequestPolicyConfig: expandCloudFrontOriginRequestPolicyConfig(d),
		Id:                        aws.String(d.Id()),
		IfMatch:                   aws.String(d.Get("etag").(string)),
	}

	_, err := conn.UpdateOriginRequestPolicy(request)
	if err != nil {
		return err
	}

	return resourceOriginRequestPolicyRead(d, meta)
}

func resourceOriginRequestPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).CloudFrontConn

	request := &cloudfront.DeleteOriginRequestPolicyInput{
		Id:      aws.String(d.Id()),
		IfMatch: aws.String(d.Get("etag").(string)),
	}

	_, err := conn.DeleteOriginRequestPolicy(request)
	if err != nil {
		if tfawserr.ErrMessageContains(err, cloudfront.ErrCodeNoSuchOriginRequestPolicy, "") {
			return nil
		}
		return err
	}

	return nil
}
