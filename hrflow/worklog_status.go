package hrflow

// WorkLogStatus defines the possible statuses for work logs.
type WorkLogStatus string

const (
	WorkLogStatusNew             = "NEW"
	WorkLogStatusSent            = "SENT"
	WorkLogStatusRejected        = "REJECTED"
	WorkLogStatusTransferred     = "TRANSFERRED"
	WorkLogStatusWaitingApproval = "WAITINGAPPROVAL"
)
