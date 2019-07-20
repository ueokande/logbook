package types

type PodStatus string

const (
	PodRunning      PodStatus = "Running"
	PodSucceeded              = "Succeeded"
	PodPending                = "Pending"
	PodTerminating            = "Terminating"
	PodInitializing           = "Initializing"
	PodFailed                 = "Failed"
	PodUnknown                = "Unknown"
)
