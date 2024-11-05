// Generated by PMS #400
package gaussdb

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceGaussdbMysqlEngineVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussdbMysqlEngineVersionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"database_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the DB engine.`,
			},
			"datastores": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the DB version list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the DB version ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the DB version number.`,
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the compatible open-source DB version.`,
						},
						"kernel_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the DB version.`,
						},
					},
				},
			},
		},
	}
}

type MysqlEngineVersionsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newMysqlEngineVersionsDSWrapper(d *schema.ResourceData, meta interface{}) *MysqlEngineVersionsDSWrapper {
	return &MysqlEngineVersionsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceGaussdbMysqlEngineVersionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newMysqlEngineVersionsDSWrapper(d, meta)
	shoGauMySqlEngVerRst, err := wrapper.ShowGaussMySqlEngineVersion()
	if err != nil {
		return diag.FromErr(err)
	}

	err = wrapper.showGaussMySqlEngineVersionToSchema(shoGauMySqlEngVerRst)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	return nil
}

// @API GAUSSDB GET /v3/{project_id}/datastores/{database_name}
func (w *MysqlEngineVersionsDSWrapper) ShowGaussMySqlEngineVersion() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "gaussdb")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/datastores/{database_name}"
	uri = strings.ReplaceAll(uri, "{database_name}", w.Get("database_name").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *MysqlEngineVersionsDSWrapper) showGaussMySqlEngineVersionToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("datastores", schemas.SliceToList(body.Get("datastores"),
			func(datastores gjson.Result) any {
				return map[string]any{
					"id":             datastores.Get("id").Value(),
					"name":           datastores.Get("name").Value(),
					"version":        datastores.Get("version").Value(),
					"kernel_version": datastores.Get("kernel_version").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
