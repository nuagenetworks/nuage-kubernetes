package entity

// State represents possible states for the Nuage entity
type State int

//EntityEvents contains state description of an event
type EntityEvents struct {
	EntityEventCategory EventCategory
	EntityEventType     Event
	EntityState         State
	EntityReason        SubState
}

// Possible states for an entity
const (
	NoState     State = 0
	Running     State = 1
	Blocked     State = 2
	Paused      State = 3
	Shutdown    State = 4
	Shutoff     State = 5
	Crashed     State = 6
	PMSuspended State = 7
)

// SubState represents possible substates for each of the state
type SubState int

// Possible sub states for Running state
const (
	RunningUnknown           SubState = 0
	RunningBooted            SubState = 1
	RunningMigrated          SubState = 2
	RunningRestored          SubState = 3
	RunningFromSnapshot      SubState = 4
	RunningUnpaused          SubState = 5
	RunningMigrationCanceled SubState = 6
	RunningSaveCanceled      SubState = 7
	RunningWakeup            SubState = 8
)

// Possible sub states for Blocked state
const (
	BlockedUnknown SubState = 0
)

// Possible sub states for Paused state
const (
	PausedUnknown      SubState = 0
	PausedUser         SubState = 1
	PausedMigration    SubState = 2
	PausedSave         SubState = 3
	PausedDump         SubState = 4
	PausedIoerror      SubState = 5
	PausedWatchdog     SubState = 6
	PausedFromSnapshot SubState = 7
	PausedShuttingDown SubState = 8
	PausedSnapshot     SubState = 9
)

// Possible sub states for Shutdown state
const (
	ShutdownUnknown SubState = 0
	ShutdownUser    SubState = 1
)

// Possible sub states for Shutoff state
const (
	ShutoffUnknown      SubState = 0
	ShutoffShutdown     SubState = 1
	ShutoffDestroyed    SubState = 2
	ShutoffCrashed      SubState = 3
	ShutoffMigrated     SubState = 4
	ShutoffSaved        SubState = 5
	ShutoffFailed       SubState = 6
	ShutoffFromSnapshot SubState = 7
)

// Possible sub states for crashed state
const (
	CrashedUnknown SubState = 0
)
