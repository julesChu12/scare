package service

import (
	"context"
	"encoding/json"
	"errors"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/consts"
	"community-elderly-care-platform/internal/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrTaskInvalid  = errors.New("invalid task")
	ErrTaskConflict = errors.New("task conflict")
)

type TaskService struct {
	db          *gorm.DB
	taskRepo    *repository.TaskRepository
	requestRepo *repository.RequestRepository
	notifySvc   *NotificationService
}

func NewTaskService(db *gorm.DB, taskRepo *repository.TaskRepository, requestRepo *repository.RequestRepository, notifySvc *NotificationService) *TaskService {
	return &TaskService{
		db:          db,
		taskRepo:    taskRepo,
		requestRepo: requestRepo,
		notifySvc:   notifySvc,
	}
}

// TaskPoolFilter 任务池筛选条件
type TaskPoolFilter struct {
	StationID int64  // 站点ID，0 表示所有站点（仅 admin）
	Status    string // 任务状态，空表示默认 dispatched
}

// ListPoolWithFilter 根据筛选条件查询任务池
func (s *TaskService) ListPoolWithFilter(filter TaskPoolFilter, page, pageSize int) ([]*model.TaskAssignment, int64, error) {
	offset := (page - 1) * pageSize
	repoFilter := repository.TaskPoolFilter{
		StationID: filter.StationID,
		Status:    filter.Status,
	}
	return s.taskRepo.ListPool(repoFilter, offset, pageSize)
}

func (s *TaskService) ListByStaff(staffID int64, page, pageSize int) ([]*model.TaskAssignment, int64, error) {
	if staffID == 0 {
		return nil, 0, ErrTaskInvalid
	}
	offset := (page - 1) * pageSize
	return s.taskRepo.ListByStaff(staffID, offset, pageSize)
}

// Claim 认领任务
//
// 业务目的：工作人员从任务池中认领待处理的任务
//
// 主要流程：
// 1. 使用行锁查询任务（FOR UPDATE，防止并发认领）
// 2. 检查任务状态（必须为 dispatched）
// 3. 检查任务是否已被其他人认领
// 4. 更新任务：设置 staff_id、状态改为 claimed
// 5. 更新关联的服务请求状态为 claimed
// 6. 异步发送通知给客户（告知任务已被认领）
//
// 并发控制：使用数据库行锁确保同一任务不会被多人同时认领
//
// 幂等性：如果同一工作人员重复认领同一任务，返回成功但 changed=false
//
// 返回值：
// - task: 任务信息
// - changed: 是否发生了状态变更
// - error: 错误信息
func (s *TaskService) Claim(taskID, staffID int64) (*model.TaskAssignment, bool, error) {
	if taskID == 0 || staffID == 0 {
		return nil, false, ErrTaskInvalid
	}
	var task model.TaskAssignment
	changed := false

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&task, taskID).Error; err != nil {
			return err
		}

		if task.Status != consts.TaskStatusDispatched {
			if task.Status == consts.TaskStatusClaimed && task.StaffID == staffID {
				return nil
			}
			return ErrTaskConflict
		}
		if task.StaffID != 0 && task.StaffID != staffID {
			return ErrTaskConflict
		}

		task.StaffID = staffID
		task.Status = consts.TaskStatusClaimed
		if err := tx.Model(&task).Updates(map[string]interface{}{
			"staff_id": staffID,
			"status":   consts.TaskStatusClaimed,
		}); err.Error != nil {
			return err.Error
		}
		changed = true

		result := tx.Model(&model.ServiceRequest{}).Where("id = ? AND status = ?", task.RequestID, consts.RequestStatusDispatched).Update("status", consts.RequestStatusClaimed)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrTaskConflict
		}

		return nil
	})
	if err != nil {
		return nil, false, err
	}

	if changed {
		s.sendClaimNotification(task.RequestID)
	}
	return &task, changed, nil
}

// Complete 完成任务
//
// 业务目的：工作人员标记任务完成并上传服务凭证
//
// 主要流程：
// 1. 使用行锁查询任务（FOR UPDATE）
// 2. 验证操作权限（只有认领人可以完成）
// 3. 检查任务状态（不能是已取消的任务）
// 4. 保存完成图片（JSON 格式存储）
// 5. 更新任务：状态改为 completed
// 6. 更新关联的服务请求状态为 completed
// 7. 异步发送通知给客户（告知服务已完成，可以评价）
//
// 幂等性：如果任务已完成，返回成功但 changed=false
//
// 返回值：
// - task: 任务信息
// - changed: 是否发生了状态变更
// - error: 错误信息
func (s *TaskService) Complete(taskID, staffID int64, images []string) (*model.TaskAssignment, bool, error) {
	if taskID == 0 || staffID == 0 {
		return nil, false, ErrTaskInvalid
	}
	var task model.TaskAssignment
	changed := false

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&task, taskID).Error; err != nil {
			return err
		}
		if task.StaffID != staffID {
			return ErrTaskConflict
		}
		if task.Status == consts.TaskStatusCompleted {
			return nil
		}
		if task.Status == consts.TaskStatusCancelled {
			return ErrTaskConflict
		}

		payload, err := json.Marshal(images)
		if err != nil {
			return err
		}
		task.Status = consts.TaskStatusCompleted
		task.Images = string(payload)
		if err := tx.Model(&task).Updates(map[string]interface{}{
			"status": consts.TaskStatusCompleted,
			"images": string(payload),
		}); err.Error != nil {
			return err.Error
		}
		changed = true

		result := tx.Model(&model.ServiceRequest{}).Where("id = ?", task.RequestID).Update("status", consts.RequestStatusCompleted)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return nil, false, err
	}

	if changed {
		s.sendCompleteNotification(task.RequestID)
	}
	return &task, changed, nil
}

func (s *TaskService) sendClaimNotification(requestID int64) {
	if s.notifySvc == nil {
		return
	}
	req, err := s.requestRepo.GetByID(requestID)
	if err != nil || req.UserID == 0 {
		return
	}
	go func(userID int64) {
		_ = s.notifySvc.SendInApp(context.Background(), userID, "task claimed", "Your request has been claimed by a staff member.")
	}(req.UserID)
}

func (s *TaskService) sendCompleteNotification(requestID int64) {
	if s.notifySvc == nil {
		return
	}
	req, err := s.requestRepo.GetByID(requestID)
	if err != nil || req.UserID == 0 {
		return
	}
	go func(userID int64) {
		_ = s.notifySvc.SendInApp(context.Background(), userID, "task completed", "Your request has been completed.")
		_ = s.notifySvc.SendEmail(context.Background(), userID, "Request completed", "Your request has been completed.")
	}(req.UserID)
}

// GetByID 获取任务详情（含关联的服务请求）
func (s *TaskService) GetByID(taskID int64) (*repository.TaskWithRequest, error) {
	if taskID == 0 {
		return nil, ErrTaskInvalid
	}
	return s.taskRepo.GetByIDWithRequest(taskID)
}

// ListWithFilter 根据筛选条件查询任务列表（含关联数据）
func (s *TaskService) ListWithFilter(filter repository.TaskListFilter, page, pageSize int) ([]*repository.TaskWithRequest, int64, error) {
	offset := (page - 1) * pageSize
	return s.taskRepo.ListWithRequest(filter, offset, pageSize)
}

// Transfer 转派任务
//
// 业务目的：管理员或站长将任务转派给其他工作人员
//
// 主要流程：
// 1. 使用行锁查询任务（FOR UPDATE）
// 2. 检查任务状态（必须为 dispatched 或 claimed）
// 3. 更新任务：设置新的 staff_id、状态改为 claimed
// 4. 异步发送通知给新的工作人员
//
// 返回值：
// - task: 任务信息
// - error: 错误信息
func (s *TaskService) Transfer(taskID, newStaffID int64) (*model.TaskAssignment, error) {
	if taskID == 0 || newStaffID == 0 {
		return nil, ErrTaskInvalid
	}
	var task model.TaskAssignment

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&task, taskID).Error; err != nil {
			return err
		}

		// 只有 dispatched 或 claimed 状态可以转派
		if task.Status != consts.TaskStatusDispatched && task.Status != consts.TaskStatusClaimed {
			return ErrTaskConflict
		}

		// 如果已经是同一个人，直接返回
		if task.StaffID == newStaffID {
			return nil
		}

		task.StaffID = newStaffID
		task.Status = consts.TaskStatusClaimed
		// 使用 Updates 只更新指定字段，避免零值时间字段问题
		if err := tx.Model(&task).Updates(map[string]interface{}{
			"staff_id": newStaffID,
			"status":   consts.TaskStatusClaimed,
		}).Error; err != nil {
			return err
		}

		// 更新关联的服务请求状态
		if err := tx.Model(&model.ServiceRequest{}).Where("id = ?", task.RequestID).Update("status", consts.RequestStatusClaimed).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// 发送通知给新的工作人员
	s.sendTransferNotification(task.RequestID, newStaffID)

	return &task, nil
}

func (s *TaskService) sendTransferNotification(requestID, staffID int64) {
	if s.notifySvc == nil {
		return
	}
	go func(userID int64) {
		_ = s.notifySvc.SendInApp(context.Background(), userID, "任务指派", "您有新的任务被指派，请及时处理。")
	}(staffID)
}
