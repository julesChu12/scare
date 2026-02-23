-- =====================================================
-- 昌平区霍营街道社区养老信息分发平台 - 数据库表结构
-- 版本: v1.0.0
-- 创建日期: 2026-01-18
-- 说明: 所有表使用逻辑外键，不使用数据库外键约束
-- =====================================================

-- 设置字符集和时区
SET NAMES utf8mb4;
SET TIME_ZONE = '+08:00';

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS `scare_db`
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE `scare_db`;

-- =====================================================
-- 1. 用户表 (users)
-- 说明: 存储所有用户信息，支持5种角色
-- =====================================================
CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `phone` VARCHAR(20) NOT NULL COMMENT '手机号(登录账号)',
  `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希',
  `name` VARCHAR(50) DEFAULT NULL COMMENT '用户姓名',
  `email` VARCHAR(100) DEFAULT NULL COMMENT '邮箱地址',
  `avatar` VARCHAR(255) DEFAULT NULL COMMENT '头像URL',
  `gender` VARCHAR(10) DEFAULT NULL COMMENT '性别(male/female/other)',
  `birth_date` DATE DEFAULT NULL COMMENT '出生日期',
  `id_card` VARCHAR(64) DEFAULT NULL COMMENT '身份证号',
  `role` VARCHAR(20) NOT NULL COMMENT '角色(已废弃,仅保留兼容性)',
  `primary_role` VARCHAR(20) DEFAULT NULL COMMENT '主角色(elderly/family/staff/station_manager/admin)',
  `station_id` BIGINT DEFAULT NULL COMMENT '关联站点ID(逻辑外键,staff/station_manager有值)',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态(active/inactive/suspended)',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间(软删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_phone` (`phone`),
  KEY `idx_role` (`role`),
  KEY `idx_station_id` (`station_id`),
  KEY `idx_status` (`status`),
  KEY `idx_id_card` (`id_card`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- =====================================================
-- 2. 服务站点表 (service_stations)
-- 说明: 社区服务站点信息
-- =====================================================
CREATE TABLE IF NOT EXISTS `service_stations` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` VARCHAR(100) NOT NULL COMMENT '站点名称',
  `code` VARCHAR(50) DEFAULT NULL COMMENT '站点编码',
  `address` VARCHAR(200) DEFAULT NULL COMMENT '详细地址',
  `phone` VARCHAR(20) DEFAULT NULL COMMENT '联系电话',
  `latitude` DECIMAL(10,7) DEFAULT NULL COMMENT '纬度',
  `longitude` DECIMAL(10,7) DEFAULT NULL COMMENT '经度',
  `service_area` VARCHAR(200) DEFAULT NULL COMMENT '服务范围描述',
  `capacity` INT NOT NULL DEFAULT 10 COMMENT '服务容量(同时处理任务数)',
  `work_hours` VARCHAR(100) DEFAULT NULL COMMENT '工作时间',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态(active/inactive)',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间(软删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='服务站点表';

-- =====================================================
-- 4. 服务围栏表 (service_zones)
-- 说明: 地理围栏，存储多边形顶点坐标(JSON格式)
-- =====================================================
CREATE TABLE IF NOT EXISTS `service_zones` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `station_id` BIGINT NOT NULL COMMENT '关联站点ID(逻辑外键)',
  `name` VARCHAR(100) NOT NULL COMMENT '围栏名称',
  `points` JSON NOT NULL COMMENT '多边形顶点数组 [[lng,lat],[lng,lat],...]',
  `priority` INT NOT NULL DEFAULT 0 COMMENT '优先级(数字越大优先级越高)',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态(active/inactive)',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间(软删除)',
  PRIMARY KEY (`id`),
  KEY `idx_station_id` (`station_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='服务围栏表';

-- =====================================================
-- 5. 服务需求表 (service_requests)
-- 说明: 用户提交的养老服务需求
-- =====================================================
CREATE TABLE IF NOT EXISTS `service_requests` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `request_no` VARCHAR(50) NOT NULL COMMENT '需求编号(唯一)',
  `user_id` BIGINT NOT NULL COMMENT '提交用户ID(逻辑外键)',
  `service_type` VARCHAR(20) NOT NULL COMMENT '服务类型(常量)',
  `status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '需求状态(pending/dispatched/claimed/processing/completed/cancelled)',
  `description` TEXT DEFAULT NULL COMMENT '需求详细描述',
  `submit_location_lat` DECIMAL(10,7) DEFAULT NULL COMMENT '提交位置纬度',
  `submit_location_lng` DECIMAL(10,7) DEFAULT NULL COMMENT '提交位置经度',
  `contact_name` VARCHAR(50) DEFAULT NULL COMMENT '联系人姓名',
  `contact_phone` VARCHAR(20) DEFAULT NULL COMMENT '联系电话',
  `address` VARCHAR(200) DEFAULT NULL COMMENT '详细地址',
  `appointment_time` DATETIME DEFAULT NULL COMMENT '预约时间',
  `urgency` VARCHAR(20) NOT NULL DEFAULT 'normal' COMMENT '紧急程度(normal/urgent)',
  `station_id` BIGINT DEFAULT NULL COMMENT '分配的站点ID(逻辑外键)',
  `reject_reason` TEXT DEFAULT NULL COMMENT '拒绝原因',
  `images` JSON DEFAULT NULL COMMENT '图片URL列表',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间(软删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_request_no` (`request_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_station_id` (`station_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='服务需求表';

-- =====================================================
-- 6. 任务分配表 (task_assignments)
-- 说明: 任务派单与处理，支持一对多关系(一个需求可多次派单)
-- =====================================================
CREATE TABLE IF NOT EXISTS `task_assignments` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `request_id` BIGINT NOT NULL COMMENT '关联需求ID(逻辑外键)',
  `station_id` BIGINT NOT NULL COMMENT '分配站点ID(逻辑外键)',
  `staff_id` BIGINT DEFAULT NULL COMMENT '认领工作人员ID(逻辑外键)',
  `status` VARCHAR(20) NOT NULL DEFAULT 'dispatched' COMMENT '任务状态(dispatched/claimed/processing/completed/cancelled)',
  `claimed_at` DATETIME DEFAULT NULL COMMENT '认领时间',
  `completed_at` DATETIME DEFAULT NULL COMMENT '完成时间',
  `rating` INT NOT NULL DEFAULT 0 COMMENT '服务评分(0-5)',
  `feedback` TEXT DEFAULT NULL COMMENT '用户评价',
  `staff_notes` TEXT DEFAULT NULL COMMENT '工作人员备注',
  `images` JSON DEFAULT NULL COMMENT '完成照片URL列表',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间(软删除)',
  PRIMARY KEY (`id`),
  KEY `idx_request_id` (`request_id`),
  KEY `idx_station_id` (`station_id`),
  KEY `idx_staff_id` (`staff_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务分配表';

-- =====================================================
-- 7. 任务历史表 (task_histories)
-- 说明: 记录任务的所有操作历史(派单/认领/转派/完成等)
-- =====================================================
CREATE TABLE IF NOT EXISTS `task_histories` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `task_id` BIGINT NOT NULL COMMENT '关联任务ID(逻辑外键)',
  `request_id` BIGINT NOT NULL COMMENT '关联需求ID(逻辑外键)',
  `action` VARCHAR(20) NOT NULL COMMENT '操作类型(dispatched/claimed/transferred/completed/cancelled)',
  `operator_id` BIGINT NOT NULL COMMENT '操作人ID(逻辑外键)',
  `from_staff_id` BIGINT DEFAULT NULL COMMENT '转派来源工作人员ID(逻辑外键)',
  `to_staff_id` BIGINT DEFAULT NULL COMMENT '转派目标工作人员ID(逻辑外键)',
  `from_station_id` BIGINT DEFAULT NULL COMMENT '转派来源站点ID(逻辑外键)',
  `to_station_id` BIGINT DEFAULT NULL COMMENT '转派目标站点ID(逻辑外键)',
  `status_before` VARCHAR(20) DEFAULT NULL COMMENT '操作前状态',
  `status_after` VARCHAR(20) DEFAULT NULL COMMENT '操作后状态',
  `remark` TEXT DEFAULT NULL COMMENT '备注说明',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间(软删除)',
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`),
  KEY `idx_request_id` (`request_id`),
  KEY `idx_operator_id` (`operator_id`),
  KEY `idx_from_staff_id` (`from_staff_id`),
  KEY `idx_to_staff_id` (`to_staff_id`),
  KEY `idx_from_station_id` (`from_station_id`),
  KEY `idx_to_station_id` (`to_station_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务历史表';

-- =====================================================
-- 8. 通知表 (notifications)
-- 说明: 多渠道通知记录(站内信/邮件/短信)
-- =====================================================
CREATE TABLE IF NOT EXISTS `notifications` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` BIGINT NOT NULL COMMENT '接收用户ID(逻辑外键)',
  `title` VARCHAR(100) NOT NULL COMMENT '通知标题',
  `body` TEXT DEFAULT NULL COMMENT '通知内容',
  `type` VARCHAR(20) DEFAULT NULL COMMENT '通知类型(request_created/task_claimed/task_completed等)',
  `related_id` BIGINT DEFAULT NULL COMMENT '关联业务ID',
  `related_type` VARCHAR(20) DEFAULT NULL COMMENT '关联业务类型(request/task)',
  `channel` VARCHAR(20) NOT NULL COMMENT '渠道(in_app/email/sms)',
  `send_status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '发送状态(pending/sent/failed)',
  `sent_at` DATETIME DEFAULT NULL COMMENT '发送时间',
  `is_read` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否已读(0未读/1已读)',
  `read_at` DATETIME DEFAULT NULL COMMENT '阅读时间',
  `retry_count` INT NOT NULL DEFAULT 0 COMMENT '重试次数',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间(软删除)',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_type` (`type`),
  KEY `idx_related_id` (`related_id`),
  KEY `idx_channel` (`channel`),
  KEY `idx_send_status` (`send_status`),
  KEY `idx_is_read` (`is_read`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通知表';

-- =====================================================
-- 索引说明
-- =====================================================
-- 1. 主键索引: 所有表的 `id` 字段
-- 2. 唯一索引: phone, request_no, code 等业务唯一字段
-- 3. 普通索引: 外键字段、状态字段、时间字段
-- 4. 软删除索引: deleted_at 字段（支持GORM软删除查询）
-- =====================================================

-- =====================================================
-- 9. 新闻资讯表 (news)
-- 说明: 站点新闻、公告，station_id=0 为平台级内容
-- =====================================================
CREATE TABLE IF NOT EXISTS `news` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `station_id` BIGINT NOT NULL DEFAULT 0 COMMENT '站点ID(0=平台级/全局)',
  `title` VARCHAR(200) NOT NULL COMMENT '标题',
  `summary` VARCHAR(500) DEFAULT NULL COMMENT '摘要',
  `content` TEXT DEFAULT NULL COMMENT '正文内容',
  `cover_url` VARCHAR(255) DEFAULT NULL COMMENT '封面图片URL',
  `type` VARCHAR(20) NOT NULL DEFAULT 'news' COMMENT '类型(news/notice/activity)',
  `status` VARCHAR(20) NOT NULL DEFAULT 'draft' COMMENT '状态(draft/published/archived)',
  `author_id` BIGINT DEFAULT NULL COMMENT '作者ID(逻辑外键)',
  `publish_at` DATETIME DEFAULT NULL COMMENT '发布时间',
  `view_count` BIGINT NOT NULL DEFAULT 0 COMMENT '浏览次数',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间(软删除)',
  PRIMARY KEY (`id`),
  KEY `idx_news_station_id` (`station_id`),
  KEY `idx_news_type` (`type`),
  KEY `idx_news_status` (`status`),
  KEY `idx_news_author_id` (`author_id`),
  KEY `idx_news_publish_at` (`publish_at`),
  KEY `idx_news_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='新闻资讯表';

-- =====================================================
-- 12. 轮播图/Banner表 (banners)
-- 说明: 首页轮播图，station_id=0 为平台级内容
-- =====================================================
CREATE TABLE IF NOT EXISTS `banners` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `station_id` BIGINT NOT NULL DEFAULT 0 COMMENT '站点ID(0=平台级/全局)',
  `title` VARCHAR(100) DEFAULT NULL COMMENT '标题',
  `image_url` VARCHAR(255) NOT NULL COMMENT '图片URL',
  `link_type` VARCHAR(20) NOT NULL DEFAULT 'none' COMMENT '跳转类型(none/news/url)',
  `link_value` VARCHAR(255) DEFAULT NULL COMMENT '跳转值(新闻ID或URL)',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序(数字越大越靠前)',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态(active/inactive)',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间(软删除)',
  PRIMARY KEY (`id`),
  KEY `idx_banners_station_status` (`station_id`, `status`),
  KEY `idx_banners_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='轮播图表';

-- =====================================================
-- 13. C端客户档案表 (customer_profiles)
-- 说明: C端用户的扩展信息（地址、紧急联系人等）
-- =====================================================
CREATE TABLE IF NOT EXISTS `customer_profiles` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` BIGINT NOT NULL COMMENT '关联用户ID(逻辑外键)',
  `id_card` VARCHAR(64) DEFAULT NULL COMMENT '身份证号',
  `address` VARCHAR(200) DEFAULT NULL COMMENT '居住地址',
  `latitude` DECIMAL(10,7) DEFAULT NULL COMMENT '纬度',
  `longitude` DECIMAL(10,7) DEFAULT NULL COMMENT '经度',
  `customer_type` VARCHAR(20) DEFAULT NULL COMMENT '客户类型(elderly/family)',
  `emergency_contact` JSON DEFAULT NULL COMMENT '紧急联系人',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间(软删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_customer_user_id` (`user_id`),
  KEY `idx_customer_profiles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='C端客户档案表';
