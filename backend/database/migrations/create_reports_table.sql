-- 统计报表表
CREATE TABLE IF NOT EXISTS reports (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    
    name VARCHAR(255) NOT NULL COMMENT '报表名称',
    type VARCHAR(32) NOT NULL COMMENT '报表类型：service/performance/request/station',
    format VARCHAR(16) NOT NULL COMMENT '文件格式：xlsx/csv',
    file_path VARCHAR(512) NOT NULL COMMENT '文件存储路径',
    file_size BIGINT NOT NULL DEFAULT 0 COMMENT '文件大小（字节）',
    
    station_id BIGINT DEFAULT NULL COMMENT '站点ID，NULL表示全局',
    start_date DATE NOT NULL COMMENT '统计开始日期',
    end_date DATE NOT NULL COMMENT '统计结束日期',
    
    created_by BIGINT NOT NULL COMMENT '创建人ID',
    
    INDEX idx_reports_type (type),
    INDEX idx_reports_station_id (station_id),
    INDEX idx_reports_created_at (created_at),
    INDEX idx_reports_created_by (created_by),
    INDEX idx_reports_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='统计报表';
