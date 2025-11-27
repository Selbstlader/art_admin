package operation_log

import (
	"errors"
	"time"

	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"
)

type OperationLogService struct {
	logRepo *repository.OperationLogRepository
}

func NewOperationLogService() *OperationLogService {
	return &OperationLogService{
		logRepo: repository.NewOperationLogRepository(),
	}
}

// GetOperationLogList 获取操作日志列表
func (s *OperationLogService) GetOperationLogList(req *request.OperationLogListRequest) ([]response.OperationLogListItem, int64, error) {
	// 构建查询条件
	query := make(map[string]interface{})
	if req.Module != "" {
		query["module"] = req.Module
	}
	if req.BusinessType != "" {
		query["businessType"] = req.BusinessType
	}
	if req.OperatorName != "" {
		query["operatorName"] = req.OperatorName
	}
	if req.Status != nil {
		query["status"] = *req.Status
	}
	if req.StartTime != "" {
		query["startTime"] = req.StartTime
	}
	if req.EndTime != "" {
		query["endTime"] = req.EndTime
	}

	// 查询操作日志列表
	logs, total, err := s.logRepo.FindWithPagination(query, req.Current, req.Size)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应DTO
	result := make([]response.OperationLogListItem, 0, len(logs))
	for _, log := range logs {
		result = append(result, response.OperationLogListItem{
			ID:            log.ID,
			Module:        log.Module,
			BusinessType:  log.BusinessType,
			RequestMethod: log.RequestMethod,
			RequestURL:    log.RequestURL,
			OperatorName:  log.OperatorName,
			OperatorIP:    log.OperatorIP,
			OperatorAddr:  log.OperatorAddr,
			RequestParam:  log.RequestParam,
			ResponseData:  log.ResponseData,
			Status:        log.Status,
			ErrorMsg:      log.ErrorMsg,
			CostTime:      log.CostTime,
			UserAgent:     log.UserAgent,
			OperationTime: log.OperationTime,
		})
	}

	return result, total, nil
}

// GetOperationLogDetail 获取操作日志详情
func (s *OperationLogService) GetOperationLogDetail(id int64) (*response.OperationLogDetail, error) {
	log, err := s.logRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("操作日志不存在")
	}

	return &response.OperationLogDetail{
		OperationLogListItem: response.OperationLogListItem{
			ID:            log.ID,
			Module:        log.Module,
			BusinessType:  log.BusinessType,
			RequestMethod: log.RequestMethod,
			RequestURL:    log.RequestURL,
			OperatorName:  log.OperatorName,
			OperatorIP:    log.OperatorIP,
			OperatorAddr:  log.OperatorAddr,
			RequestParam:  log.RequestParam,
			ResponseData:  log.ResponseData,
			Status:        log.Status,
			ErrorMsg:      log.ErrorMsg,
			CostTime:      log.CostTime,
			UserAgent:     log.UserAgent,
			OperationTime: log.OperationTime,
		},
	}, nil
}

// CreateOperationLog 创建操作日志
func (s *OperationLogService) CreateOperationLog(req *request.CreateOperationLogRequest) error {
	log := &model.OperationLog{
		Module:        req.Module,
		BusinessType:  req.BusinessType,
		RequestMethod: req.RequestMethod,
		RequestURL:    req.RequestURL,
		OperatorName:  req.OperatorName,
		OperatorIP:    req.OperatorIP,
		OperatorAddr:  req.OperatorAddr,
		RequestParam:  req.RequestParam,
		ResponseData:  req.ResponseData,
		Status:        req.Status,
		ErrorMsg:      req.ErrorMsg,
		CostTime:      req.CostTime,
		UserAgent:     req.UserAgent,
		OperationTime: time.Now(),
	}

	return s.logRepo.Create(log)
}

// DeleteOperationLog 删除操作日志
func (s *OperationLogService) DeleteOperationLog(id int64) error {
	// 检查日志是否存在
	_, err := s.logRepo.FindByID(id)
	if err != nil {
		return errors.New("操作日志不存在")
	}

	return s.logRepo.Delete(id)
}

// BatchDeleteOperationLog 批量删除操作日志
func (s *OperationLogService) BatchDeleteOperationLog(ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的日志")
	}

	return s.logRepo.BatchDelete(ids)
}

// CleanOperationLog 清理操作日志
func (s *OperationLogService) CleanOperationLog(days int) error {
	if days <= 0 {
		return errors.New("保留天数必须大于0")
	}

	return s.logRepo.CleanOldLogs(days)
}
