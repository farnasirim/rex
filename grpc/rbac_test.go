package grpc

import (
	"context"
	"testing"

	"github.com/farnasirim/rex"
)

func TestSimpleAccessRule_MatchPrincipal_Wildcard_Reverse(t *testing.T) {
	rule, err := SimpleAccessRuleFromJSON([]byte(`
	{"principal": "username", "effect": "allow", "action": "/path/to/action"}
	`))
	if err != nil {
		t.Errorf("Caught error while creating simple access rule from JSON: %v", err)
	}

	ctxUser := rex.WithUserID(context.Background(), "*")
	ctxMethod := withMethodName(ctxUser, "/path/to/action")
	ctx := ctxMethod
	verdict, applies := rule.Enforce(ctx)
	if !verdict || applies {
		t.Errorf("Rule must not apply. Was expecting (true, false), got (%v, %v)", verdict, applies)
	}
}

func TestPolicyEnforcer_Enforce(t *testing.T) {

}
