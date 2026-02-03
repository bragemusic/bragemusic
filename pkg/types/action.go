package types

type Action string

func (a Action) P() *Action {
	return &a
}

const (
	ActionCreate Action = "create"
	ActionRead   Action = "read"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)
