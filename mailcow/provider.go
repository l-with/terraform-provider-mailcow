package hashicups

import (
	"os"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/l-with/mailcow-go"
)

var stderr = os.Stderr

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	configured bool
	client     *mailcow.Client
}
