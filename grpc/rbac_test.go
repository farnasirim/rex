package grpc

import (
	"context"
	"testing"

	"github.com/farnasirim/rex"
)

func TestSimpleAccessRule_MatchPrincipal_Wildcard_Reverse(t *testing.T) {
	rule, err := SimpleAccessRuleFromJSON([]byte(`
	{"principal": "username", "effect": "allow", "action": "/Rex/Exec"}
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

type policyMock struct {
	EnforceFunc func(context.Context) (bool, bool)
}

func (m *policyMock) Enforce(ctx context.Context) (bool, bool) {
	return m.EnforceFunc(ctx)
}

func TestPolicyEnforcer_Enforce_DefaultToDeny(t *testing.T) {
	enforcer := NewPolicyEnforcer(
		&policyMock{func(context.Context) (bool, bool) { return true, false }},
		&policyMock{func(context.Context) (bool, bool) { return true, false }})

	ctxUser := rex.WithUserID(context.Background(), "*")
	ctxMethod := withMethodName(ctxUser, "/path/to/action")
	verdict, applies := enforcer.Enforce(ctxMethod)
	if verdict || applies {
		t.Errorf("Rule must not apply. Was expecting (false, false), got (%v, %v)", verdict, applies)
	}
}
