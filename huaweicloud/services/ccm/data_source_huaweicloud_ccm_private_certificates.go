// Generated by PMS #256
package ccm

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourcePrivateCertificates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateCertificatesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the private certificate.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the private certificate status.`,
			},
			"sort_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sorting attribute.`,
			},
			"sort_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sorting sequence.`,
			},
			"certificates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The certificate details.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"distinguished_name": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The certificate name configuration.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"locality": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The country or region name.`,
									},
									"organization": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The organization name.`,
									},
									"organizational_unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The organization unit.`,
									},
									"common_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The common certificate name.`,
									},
									"country": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The country code.`,
									},
									"state": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The state or city name.`,
									},
								},
							},
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the private certificate.`,
						},
						"issuer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the parent CA.`,
						},
						"key_algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The key algorithm.`,
						},
						"signature_algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The signature algorithm.`,
						},
						"gen_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The certificate generation method.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the certificate, in RFC3339 format.`,
						},
						"issuer_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the parent CA certificate.`,
						},
						"start_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The private certificate validity start from, in RFC3339 format.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The certificate status.`,
						},
						"expired_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The certificate expiration time, in RFC3339 format.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The enterprise project ID.`,
						},
					},
				},
			},
		},
	}
}

type PrivateCertificatesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newPrivateCertificatesDSWrapper(d *schema.ResourceData, meta interface{}) *PrivateCertificatesDSWrapper {
	return &PrivateCertificatesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourcePrivateCertificatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newPrivateCertificatesDSWrapper(d, meta)
	listCertificateRst, err := wrapper.ListCertificate()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listCertificateToSchema(listCertificateRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CCM GET /v1/private-certificates
func (w *PrivateCertificatesDSWrapper) ListCertificate() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "ccm")
	if err != nil {
		return nil, err
	}

	uri := "/v1/private-certificates"
	params := map[string]any{
		"name":     w.Get("name"),
		"status":   w.Get("status"),
		"sort_key": w.Get("sort_key"),
		"sort_dir": w.Get("sort_dir"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("certificates", "offset", "limit", 1000).
		Request().
		Result()
}

func (w *PrivateCertificatesDSWrapper) listCertificateToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("certificates", schemas.SliceToList(body.Get("certificates"),
			func(certificates gjson.Result) any {
				return map[string]any{
					"distinguished_name": schemas.SliceToList(certificates.Get("distinguished_name"),
						func(distinguishedName gjson.Result) any {
							return map[string]any{
								"locality":            distinguishedName.Get("locality").Value(),
								"organization":        distinguishedName.Get("organization").Value(),
								"organizational_unit": distinguishedName.Get("organizational_unit").Value(),
								"common_name":         distinguishedName.Get("common_name").Value(),
								"country":             distinguishedName.Get("country").Value(),
								"state":               distinguishedName.Get("state").Value(),
							}
						},
					),
					"id":                    certificates.Get("certificate_id").Value(),
					"issuer_id":             certificates.Get("issuer_id").Value(),
					"key_algorithm":         certificates.Get("key_algorithm").Value(),
					"signature_algorithm":   certificates.Get("signature_algorithm").Value(),
					"gen_mode":              certificates.Get("gen_mode").Value(),
					"created_at":            w.setPrivateCertificateCreatedAt(certificates),
					"issuer_name":           certificates.Get("issuer_name").Value(),
					"start_at":              w.setPrivateCertificateStartAt(certificates),
					"status":                certificates.Get("status").Value(),
					"expired_at":            w.setPrivateCertificateExpiredAt(certificates),
					"enterprise_project_id": certificates.Get("enterprise_project_id").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (*PrivateCertificatesDSWrapper) setPrivateCertificateCreatedAt(data gjson.Result) string {
	return utils.FormatTimeStampRFC3339(data.Get("create_time").Int()/1000, false)
}

func (*PrivateCertificatesDSWrapper) setPrivateCertificateStartAt(data gjson.Result) string {
	return utils.FormatTimeStampRFC3339(data.Get("not_before").Int()/1000, false)
}

func (*PrivateCertificatesDSWrapper) setPrivateCertificateExpiredAt(data gjson.Result) string {
	return utils.FormatTimeStampRFC3339(data.Get("not_after").Int()/1000, false)
}
