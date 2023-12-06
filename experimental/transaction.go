package experimental

import (
	"github.com/crowdsecurity/coraza/v3/collection"
	"github.com/crowdsecurity/coraza/v3/experimental/plugins/plugintypes"
	"github.com/crowdsecurity/coraza/v3/types"
	"github.com/crowdsecurity/coraza/v3/types/variables"
)

func ToFullInterface(waf types.Transaction) FullTransaction {
	x := waf.(FullTransaction)
	return x
}

type FullTransaction interface {
	types.Transaction
	RemoveRuleByID(int)
	RemoveRulesByID(...int)
	RemoveRuleByTag(string)
	RemoveRulesByTag(...string)
	Variables() plugintypes.TransactionVariables
	Collection(variables.RuleVariable) collection.Collection
}
