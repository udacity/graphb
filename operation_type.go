package graphb

type operationType string

const (
	TypeQuery        operationType = "query"
	TypeMutation     operationType = "mutation"
	TypeSubscription operationType = "subscription"
)
