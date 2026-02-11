package consts

// 服务请求状态
const (
	RequestStatusPending    = "pending"
	RequestStatusDispatched = "dispatched"
	RequestStatusClaimed    = "claimed"
	RequestStatusProcessing = "processing"
	RequestStatusCompleted  = "completed"
	RequestStatusCancelled  = "cancelled"
	RequestStatusRejected   = "rejected"
)

// IsValidRequestStatus 验证请求状态是否合法
func IsValidRequestStatus(status string) bool {
	switch status {
	case RequestStatusPending, RequestStatusDispatched, RequestStatusClaimed,
		RequestStatusProcessing, RequestStatusCompleted, RequestStatusCancelled, RequestStatusRejected:
		return true
	}
	return false
}

// 任务状态
const (
	TaskStatusDispatched = "dispatched"
	TaskStatusClaimed    = "claimed"
	TaskStatusProcessing = "processing"
	TaskStatusCompleted  = "completed"
	TaskStatusCancelled  = "cancelled"
)
