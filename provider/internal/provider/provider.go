package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type FilegenProvider struct {
	version string
}

var _ provider.Provider = &FilegenProvider{}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &FilegenProvider{
			version: version,
		}
	}
}

func (p *FilegenProvider) Metadata(
	ctx context.Context,
	req provider.MetadataRequest,
	resp *provider.MetadataResponse) {
	resp.TypeName = "filegen"
	resp.Version = p.version
}

func (p *FilegenProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewFilledFileResource,
	}
}

func (p *FilegenProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func (p *FilegenProvider) Schema(
	ctx context.Context,
	req provider.SchemaRequest,
	resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *FilegenProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse) {

}
