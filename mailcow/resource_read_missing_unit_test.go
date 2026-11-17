package mailcow

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/l-with/terraform-provider-mailcow/api"
)

// clientAnswering is an APIClient whose every GET is answered with body.
func clientAnswering(t *testing.T, body string) *APIClient {
	t.Helper()
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	config := api.NewConfiguration()
	config.Host = strings.TrimPrefix(server.URL, "https://")
	config.Scheme = "https"
	config.HTTPClient = server.Client()
	return &APIClient{client: api.NewAPIClient(config)}
}

// mailcow answers an id it holds nothing for with an empty object, not an
// error. Read must take that as gone — removing the resource, so an import
// fails with "Cannot import non-existent remote object" and a refresh plans it
// anew — rather than as an object with empty fields, which a plan "replaces".
func TestReadOfAMissingObjectRemovesIt(t *testing.T) {
	cases := []struct {
		name string
		res  *schema.Resource
		read func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics
		id   string
		body string
		gone bool
	}{
		{"missing mailbox", resourceMailbox(), resourceMailboxRead, "ghost@example.com", `{}`, true},
		{"existing mailbox", resourceMailbox(), resourceMailboxRead, "ann@example.com",
			`{"local_part":"ann","domain":"example.com","name":"Ann","quota":0,"active":1}`, false},
		{"existing mailbox, mixed-case id", resourceMailbox(), resourceMailboxRead, "Ann@Example.com",
			`{"local_part":"ann","domain":"example.com","name":"Ann","quota":0,"active":1}`, false},
		{"missing dkim", resourceDkim(), resourceDkimRead, "example.com", `{}`, true},
		{"existing dkim", resourceDkim(), resourceDkimRead, "example.com",
			`{"dkim_selector":"dkim","length":"2048","pubkey":"MIIB"}`, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, c.res.Schema, map[string]interface{}{})
			d.SetId(c.id)
			if diags := c.read(context.Background(), d, clientAnswering(t, c.body)); diags.HasError() {
				t.Fatalf("read: %v", diags)
			}
			if gone := d.Id() == ""; gone != c.gone {
				t.Errorf("removed from the state = %v, want %v", gone, c.gone)
			}
		})
	}
}

// mailcow looks a mailbox up regardless of case and answers in lower case, so
// the data source must find an existing mailbox asked for with a mixed-case
// address and still report a missing one.
func TestDataSourceMailboxAddressCase(t *testing.T) {
	cases := []struct {
		name    string
		address string
		body    string
		found   bool
	}{
		{"existing mailbox, mixed-case address", "Ann@Example.com",
			`{"local_part":"ann","domain":"example.com","name":"Ann","quota":0,"active":1,"attributes":{}}`, true},
		{"missing mailbox", "ghost@example.com", `{}`, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := dataSourceMailbox()
			d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{"address": c.address})
			diags := dataSourceMailboxRead(context.Background(), d, clientAnswering(t, c.body))
			if found := !diags.HasError(); found != c.found {
				t.Errorf("found = %v, want %v: %v", found, c.found, diags)
			}
		})
	}
}
