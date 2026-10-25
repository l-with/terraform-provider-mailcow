package mailcow

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
)

// TestMailboxSchemaIsValid runs the SDK's own checks, which refuse a write-only
// attribute that is also Computed, ForceNew or has a Default.
func TestMailboxSchemaIsValid(t *testing.T) {
	if err := resourceMailbox().InternalValidate(nil, true); err != nil {
		t.Fatalf("mailbox schema is not valid: %s", err)
	}
}

// TestMailboxPasswordSchema pins the shape of the two password arguments:
// password_wo never reaches the state, and exactly one of the two is given.
func TestMailboxPasswordSchema(t *testing.T) {
	s := resourceMailbox().Schema
	if !s["password_wo"].WriteOnly || !s["password_wo"].Sensitive {
		t.Error("password_wo must be write-only and sensitive")
	}
	if s["password"].Required {
		t.Error("password must be optional now that password_wo can stand in for it")
	}
	for _, name := range []string{"password", "password_wo"} {
		got := s[name].ExactlyOneOf
		if len(got) != 2 || got[0] != "password" || got[1] != "password_wo" {
			t.Errorf("%s: ExactlyOneOf = %v, want [password password_wo]", name, got)
		}
	}
	if got := s["password_wo_version"].RequiredWith; len(got) != 1 || got[0] != "password_wo" {
		t.Errorf("password_wo_version: RequiredWith = %v, want [password_wo]", got)
	}
	for _, name := range []string{"password", "password_wo", "password_wo_version"} {
		if !isElementIn(name, &passwordFields) {
			t.Errorf("%s must be kept out of the generic create/update/read helpers", name)
		}
	}
}

func TestMailboxPassword(t *testing.T) {
	testCases := []struct {
		name      string
		password  string
		writeOnly cty.Value
		want      string
		wantErr   bool
	}{
		{"password wins", "plain", cty.StringVal("wo"), "plain", false},
		{"write-only when password is unset", "", cty.StringVal("wo"), "wo", false},
		{"neither", "", cty.NullVal(cty.String), "", true},
		{"write-only not known yet", "", cty.UnknownVal(cty.String), "", true},
		{"write-only empty", "", cty.StringVal(""), "", true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := mailboxPassword(tc.password, tc.writeOnly)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
