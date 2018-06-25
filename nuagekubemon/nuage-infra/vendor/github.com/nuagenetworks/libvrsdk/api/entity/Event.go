package entity

// EventCategory represents different types of events that can be associated with an Entity
type EventCategory int

// Event identifies a specific Entity event for an EventCategory
type Event int

// Categories of events supported by Nuage Entity
const (
	EventCategoryDefined     EventCategory = 0
	EventCategoryUndefined   EventCategory = 1
	EventCategoryStarted     EventCategory = 2
	EventCategorySuspended   EventCategory = 3
	EventCategoryResumed     EventCategory = 4
	EventCategoryStopped     EventCategory = 5
	EventCategoryShutdown    EventCategory = 6
	EventCategoryPmsuspended EventCategory = 7
	EventCategoryStarting    EventCategory = 10
)

// Defined events
const (
	EventDefinedAdded   Event = 0
	EventDefinedUpdated Event = 1
	EventDefinedLast    Event = 2
)

// Undefined events
const (
	EventUndefinedRemoved Event = 0
	EventUndefinedLast    Event = 1
)

// Start events
const (
	EventStartedBooted       Event = 0
	EventStartedMigrated     Event = 1
	EventStartedRestored     Event = 2
	EventStartedFromSnapshot Event = 3
	EventStartedWakeup       Event = 4
	EventStartedLast         Event = 5
)

// Suspend events
const (
	EventSuspendedPaused       Event = 0
	EventSuspendedMigrated     Event = 1
	EventSuspendedIOError      Event = 2
	EventSuspendedWatchdog     Event = 3
	EventSuspendedRestored     Event = 4
	EventSuspendedFromSnapshot Event = 5
	EventSuspendedAPIError     Event = 6
	EventSuspendedLast         Event = 7
)

// Resume events
const (
	EventResumedUnpaused     Event = 0
	EventResumedMigrated     Event = 1
	EventResumedFromSnapshot Event = 2
	EventResumedLast         Event = 3
)

// Stop events
const (
	EventStoppedShutdown     Event = 0
	EventStoppedDestroyed    Event = 1
	EventStoppedCrashed      Event = 2
	EventStoppedMigrated     Event = 3
	EventStoppedSaved        Event = 4
	EventStoppedFailed       Event = 5
	EventStoppedFromSnapshot Event = 6
	EventStoppedLast         Event = 7
)

// Shutdown events
const (
	EventShutdownFinished Event = 0
	EventShutdownLast     Event = 1
)

// PM suspend events
const (
	EventPMSuspendedMemory Event = 0
	EventPMSuspendedDisk   Event = 1
	EventPMSuspendedLast   Event = 2
)

// ValidateEvent validates the event type given it's category
func ValidateEvent(eventCategory EventCategory, event Event) bool {
	valid := false
	switch eventCategory {
	case EventCategoryDefined:
		valid = int(event) >= 0 && event < EventDefinedLast
	case EventCategoryUndefined:
		valid = int(event) >= 0 && event < EventUndefinedLast
	case EventCategoryStarted:
		valid = int(event) >= 0 && event < EventStartedLast
	case EventCategorySuspended:
		valid = int(event) >= 0 && event < EventSuspendedLast
	case EventCategoryResumed:
		valid = int(event) >= 0 && event < EventResumedLast
	case EventCategoryStopped:
		valid = int(event) >= 0 && event < EventStoppedLast
	case EventCategoryShutdown:
		valid = int(event) >= 0 && event < EventShutdownLast
	case EventCategoryPmsuspended:
		valid = int(event) >= 0 && event < EventPMSuspendedLast
	}

	return valid
}
