package grpc

import (
	"context"
	"encoding/json"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-playground/validator/v10"

	"github.com/farnasirim/rex"
)

// Policy is an abstraction for enforcing arbitrary rules with boolean outcome.
type Policy interface {
	// Enforce enforces a yes/no effect verdict on the given context and
	// returns "applies" (second return value) as true if the context falls
	// under the policy conditions. Otherwise it returns "applies" as false.
	Enforce(context.Context) (effect bool, applies bool)
}

// SimpleAccessRule defines access rules of the form
// "User is/is not allowed to execute Action"
type SimpleAccessRule struct {
	Principal string `validate:"required"`
	// TODO: extract method names from the grpc service. Currently I don't see
	// a clean way to do this. We can register a dummy service which will
	// introduce proto._Rex_serviceDesc to the grpc server upon registration,
	// They GetServiceInfo will reveal the method names. However to avoid
	// reconstructing the full method name ourselves, we need to extract those
	// directly from interceptors. One idea would be to loop over the method
	// names of the above dummy service and through server reflection send a
	// dummy request to each of its endpoints, allowing for the interceptor
	// to be invoked. There we steal the full name using UnaryServerInfo.
	// All of this happens before server startup time.
	Action string `validate:"oneof=* /Rex/Exec /Rex/Kill /Rex/GetProcessInfo /Rex/ListProcessInfo /Rex/Read"`
	Effect string `validate:"oneof=allow deny"`
}

// Enforce returns (lowercase(Effect) == "allow", true) if principal and action
// match those in the context, (false, false) otherwise. Its Principal and
// Action fields support "*" to always match.
func (r *SimpleAccessRule) Enforce(ctx context.Context) (bool, bool) {
	userID, ok := rex.UserIDFromContext(ctx)
	if !ok {
		return false, false
	}
	methodName, ok := methodNameFromContext(ctx)
	if !ok {
		return false, false
	}

	return r.effect(), r.matchPrincipal(userID) && r.matchAction(methodName)
}

func (r *SimpleAccessRule) matchPrincipal(principal string) bool {
	return wildcardMatch(r.Principal, principal)
}

func (r *SimpleAccessRule) matchAction(action string) bool {
	return wildcardMatch(r.Action, action)
}

func (r *SimpleAccessRule) effect() bool {
	return strings.ToLower(r.Effect) == "allow"
}

func wildcardMatch(maybeWildcard, str string) bool {
	return maybeWildcard == "*" || maybeWildcard == str
}

// SimpleAccessRuleFromJSON creates an access rule from its json representation
func SimpleAccessRuleFromJSON(marshalledAccessRule []byte) (*SimpleAccessRule, error) {
	validate := validator.New()

	var rule SimpleAccessRule
	if err := json.Unmarshal(marshalledAccessRule, &rule); err != nil {
		return nil, err
	}
	rule.Effect = strings.ToLower(rule.Effect)
	if err := validate.Struct(&rule); err != nil {
		return nil, err
	}

	return &rule, nil
}

// PolicyEnforcer implements Policy by chaining together other policies and
// executing them one by one, defaulting
type PolicyEnforcer struct {
	policies []Policy
}

// NewPolicyEnforcer creates a PolicyEnforcer from a list (ordered chain)
// of objects satisfying the Policy interface
func NewPolicyEnforcer(policies ...Policy) *PolicyEnforcer {
	return &PolicyEnforcer{
		policies: policies,
	}
}

// Enforce applies []policies on the context. Returns (false, false) if
// none of the policies apply, returns (false, true) if at least one of
// the policies returns (false, true), and returns (true, true) otherwise.
func (e *PolicyEnforcer) Enforce(ctx context.Context) (bool, bool) {
	verdict := false
	applies := false

	for _, policy := range e.policies {
		thisVerdict, thisApplies := policy.Enforce(ctx)
		if thisApplies {
			if !thisVerdict {
				return false, true
			}
			applies = true
			verdict = true
		}
	}

	return verdict, applies
}

// PolicyEnforcementInterceptor authorizes the execution of handler
func PolicyEnforcementInterceptor(p Policy) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx = withMethodName(ctx, info.FullMethod)

		if authorized, applies := p.Enforce(ctx); !applies || !authorized {
			return nil, status.Errorf(codes.PermissionDenied,
				rex.ErrAccessDenied.Error())
		}

		return handler(ctx, req)
	}
}
