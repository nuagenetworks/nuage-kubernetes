package entity

// Type defines the different types of entity
type Type int

// Different types of entities
const (
	VM        = 0
	Host      = 1
	Bridge    = 2
	VMBridge  = 3
	Container = 4
)
