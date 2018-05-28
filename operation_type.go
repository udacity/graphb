package graphb

type operationType string

// 3 types of operation.
const (
	TypeQuery        operationType = "query"
	TypeMutation     operationType = "mutation"
	TypeSubscription operationType = "subscription"
)
