package mailcow

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// fakeMailcow answers the four mailbox calls the resource makes and records
// every password it is sent, so a test can tell what went out and when.
type fakeMailcow struct {
	mu        sync.Mutex
	mailboxes map[string]map[string]interface{}
	passwords []string // in the order they were sent, create and edit alike
}

func (f *fakeMailcow) handler(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	ok := func(msg string) { fmt.Fprintf(w, `[{"type":"success","log":[],"msg":[%q]}]`, msg) }

	if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/v1/get/mailbox/") {
		id, _ := url.PathUnescape(strings.TrimPrefix(r.URL.Path, "/api/v1/get/mailbox/"))
		box, found := f.mailboxes[id]
		if !found {
			box = map[string]interface{}{}
		}
		_ = json.NewEncoder(w).Encode(box)
		return
	}

	raw, _ := io.ReadAll(r.Body)
	var body map[string]interface{}
	_ = json.Unmarshal(raw, &body)
	switch r.URL.Path {
	case "/api/v1/add/mailbox":
		address := fmt.Sprintf("%v@%v", body["local_part"], body["domain"])
		f.passwords = append(f.passwords, fmt.Sprint(body["password"]))
		f.mailboxes[address] = map[string]interface{}{
			"username": address, "domain": body["domain"], "local_part": body["local_part"],
			"name": body["name"], "active": 1, "authsource": "mailcow", "quota": 0,
			"attributes": map[string]interface{}{
				"force_pw_update": "1", "tls_enforce_in": "0", "tls_enforce_out": "0", "sogo_access": "1",
				"imap_access": "1", "pop3_access": "1", "smtp_access": "1", "sieve_access": "1",
			},
		}
		ok("mailbox_added")
	case "/api/v1/edit/mailbox":
		if attr, isMap := body["attr"].(map[string]interface{}); isMap {
			if password, sent := attr["password"]; sent {
				f.passwords = append(f.passwords, fmt.Sprint(password))
			}
		}
		ok("mailbox_modified")
	case "/api/v1/delete/mailbox":
		if items, isList := body["items"].([]interface{}); isList {
			for _, item := range items {
				delete(f.mailboxes, fmt.Sprint(item))
			}
		}
		ok("mailbox_removed")
	default:
		http.Error(w, "unexpected "+r.URL.Path, http.StatusNotFound)
	}
}

func (f *fakeMailcow) sent() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.passwords...)
}

func testMailboxPasswordWoConfig(host, password string, version int) string {
	return fmt.Sprintf(`
provider "mailcow" {
  host_name = %q
  api_key   = "not-a-key"
  insecure  = true
}

resource "mailcow_mailbox" "wo" {
  domain              = "example.org"
  local_part          = "wo"
  full_name           = "Write Only"
  password_wo         = %q
  password_wo_version = %d
}
`, host, password, version)
}

// noValueInState fails when the secret appears anywhere in the state — the
// point of a write-only argument.
func noValueInState(secret string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			for key, value := range rs.Primary.Attributes {
				if strings.Contains(value, secret) {
					return fmt.Errorf("%s.%s holds the write-only password", name, key)
				}
			}
		}
		return nil
	}
}

// TestMailboxPasswordWriteOnlyApply runs Terraform against fakeMailcow: the
// password is sent on create, sent again only when the version changes, and
// never stored.
func TestMailboxPasswordWriteOnlyApply(t *testing.T) {
	fake := &fakeMailcow{mailboxes: map[string]map[string]interface{}{}}
	server := httptest.NewTLSServer(http.HandlerFunc(fake.handler))
	defer server.Close()
	host := strings.TrimPrefix(server.URL, "https://")

	expectSent := func(want ...string) resource.TestCheckFunc {
		return func(*terraform.State) error {
			got := fake.sent()
			if strings.Join(got, ",") != strings.Join(want, ",") {
				return fmt.Errorf("passwords sent: %v, want %v", got, want)
			}
			return nil
		}
	}

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testMailboxPasswordWoConfig(host, "first-secret-1", 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_mailbox.wo", "address", "wo@example.org"),
					resource.TestCheckNoResourceAttr("mailcow_mailbox.wo", "password_wo"),
					noValueInState("first-secret-1"),
					expectSent("first-secret-1"),
				),
			},
			{
				// A new value under the same version is not noticed: nothing is sent.
				Config: testMailboxPasswordWoConfig(host, "ignored-secret", 1),
				Check: resource.ComposeTestCheckFunc(
					noValueInState("ignored-secret"),
					expectSent("first-secret-1"),
				),
			},
			{
				Config: testMailboxPasswordWoConfig(host, "second-secret-2", 2),
				Check: resource.ComposeTestCheckFunc(
					noValueInState("second-secret-2"),
					expectSent("first-secret-1", "second-secret-2"),
				),
			},
		},
	})
}
