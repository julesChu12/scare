-- =====================================================
-- 昌平区霍营街道社区养老信息分发平台 - 数据库表结构
-- 版本: v1.1.0
-- 更新日期: 2026-02-24
-- 说明: 所有表使用逻辑外键，不使用数据库外键约束
--       本文件以 GORM Gen 模型 (.gen.go) 为唯一真相源
-- =====================================================

SET NAMES utf8mb4;
SET TIME_ZONE = '+08:00';

CREATE DATABASE IF NOT EXISTS `scare_db`
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE `scare_db`;

-- =====================================================
-- 1. 用户表 (users)
-- =====================================================
CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `phone` VARCHAR(20) NOT NULL,
  `password_hash` VARCHAR(255) DEFAULT NULL,
  `name` VARCHAR(50) DEFAULT NULL,
  `email` VARCHAR(100) DEFAULT NULL,
  `avatar` VARCHAR(255) DEFAULT NULL,
  `gender` VARCHAR(10) DEFAULT NULL,
  `birth_date` DATE DEFAULT NULL,
  `id_card` VARCHAR(64) DEFAULT NULL,
  `station_id` BIGINT DEFAULT NULL,
  `status` VARCHAR(20) DEFAULT 'active',
  `created_at` DATETIME DEFAULT NULL,
  `updated_at` DATETIME DEFAULT NULL,
  `deleted_at` DATETIME DEFAULT NULL,
  `id_card_hmac` VARCHAR(64) DEFAULT NULL COMMENT '身份证号HMAC摘要',
  `id_card_masked` VARCHAR(20) DEFAULT NULL COMMENT '身份证号脱敏值',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_phone` (`phone`),
  KEY `idx_users_email` (`email`),
  KEY `idx_id_card` (`id_card`),
  KEY `idx_users_id_card` (`id_card`),
  KEY `idx_station_id` (`station_id`),
  KEY `idx_users_station_id` (`station_id`),
  KEY `idx_status` (`status`),
  KEY `idx_users_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_users_deleted_at` (`deleted_at`),
  KEY `idx_users_id_card_hmac` (`id_card_hmac`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- =====================================================
-- 2. 客户档案表 (customer_profiles)
-- =====================================================
CREATE TABLE IF NOT EXISTS `customer_profiles` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` BIGINT NOT NULL COMMENT '关联用户ID',
  `id_card` VARCHAR(64) DEFAULT NULL COMMENT '身份证号',
  `address` TEXT DEFAULT NULL COMMENT '居住地址',
  `latitude` DECIMAL(10,7) DEFAULT NULL COMMENT '纬度',
  `longitude` DECIMAL(10,7) DEFAULT NULL COMMENT '经度',
  `customer_type` VARCHAR(20) DEFAULT NULL COMMENT '客户类型：elderly/disabled/pregnant/child/other',
  `emergency_contact` JSON DEFAULT NULL COMMENT '紧急联系人',
  `created_at` DATETIME DEFAULT NULL,
  `updated_at` DATETIME DEFAULT NULL,
  `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间(软删除)',
  `gender` VARCHAR(10) DEFAULT NULL COMMENT '性别',
  `birth_date` DATE DEFAULT NULL COMMENT '出生日期',
  `health_status` VARCHAR(50) DEFAULT NULL COMMENT '健康状况',
  `disability_level` VARCHAR(20) DEFAULT NULL COMMENT '失能等级：自理/轻度/中度/重度',
  `medical_history` TEXT DEFAULT NULL COMMENT '病史',
  `special_needs` TEXT DEFAULT NULL COMMENT '特殊需求',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_id` (`user_id`),
  KEY `idx_id_card` (`id_card`),
  KEY `idx_customer_type` (`customer_type`),
  KEY `idx_customer_profiles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='C端客户档案表';

-- =====================================================
-- 3. 用户身份表 (user_identities)
-- =====================================================
CREATE TABLE IF NOT EXISTS `user_identities` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `identity_type` VARCHAR(20) NOT NULL COMMENT '身份类型: B端(admin/station_manager/staff) C端(elderly/family)',
  `is_primary` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否主身份(每用户有且仅有一个)',
  `station_id` BIGINT DEFAULT NULL COMMENT '所属服务站ID(仅B端身份)',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '身份状态(active/inactive)',
  `granted_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '授予时间',
  `granted_by` BIGINT DEFAULT NULL COMMENT '授予人ID',
  `revoked_at` DATETIME DEFAULT NULL COMMENT '撤销时间',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间(软删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_identity` (`user_id`, `identity_type`, `deleted_at`),
  KEY `idx_ui_user_id` (`user_id`),
  KEY `idx_ui_identity_type` (`identity_type`),
  KEY `idx_ui_station_id` (`station_id`),
  KEY `idx_ui_status` (`status`),
  KEY `idx_ui_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户身份表';

-- =====================================================
-- 4. 服务站点表 (service_stations)
-- =====================================================
CREATE TABLE IF NOT EXISTS `service_stations` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` VARCHAR(100) NOT NULL,
  `code` VARCHAR(50) DEFAULT NULL,
  `address` VARCHAR(200) DEFAULT NULL,
  `phone` VARCHAR(20) DEFAULT NULL,
  `latitude` DECIMAL(10,7) DEFAULT NULL,
  `longitude` DECIMAL(10,7) DEFAULT NULL,
  `service_area` VARCHAR(200) DEFAULT NULL,
  `capacity` BIGINT DEFAULT 10,
  `work_hours` VARCHAR(100) DEFAULT NULL,
  `status` VARCHAR(20) DEFAULT 'active',
  `created_at` DATETIME DEFAULT NULL,
  `updated_at` DATETIME DEFAULT NULL,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_service_stations_code` (`code`),
  KEY `idx_service_stations_status` (`status`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_service_stations_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='服务站点表';

-- =====================================================
-- 5. 服务围栏表 (service_zones)
-- =====================================================
CREATE TABLE IF NOT EXISTS `service_zones` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `station_id` BIGINT NOT NULL,
  `name` VARCHAR(100) NOT NULL,
  `points` JSON NOT NULL,
  `priority` BIGINT DEFAULT NULL,
  `status` VARCHAR(20) DEFAULT 'active',
  `created_at` DATETIME DEFAULT NULL,
  `updated_at` DATETIME DEFAULT NULL,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_service_zones_station_id` (`station_id`),
  KEY `idx_station_id` (`station_id`),
  KEY `idx_service_zones_status` (`status`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_service_zones_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='服务围栏表';

-- =====================================================
-- 6. 服务需求表 (service_requests)
-- =====================================================
CREATE TABLE IF NOT EXISTS `service_requests` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `request_no` VARCHAR(50) NOT NULL,
  `user_id` BIGINT NOT NULL,
  `service_type` VARCHAR(20) NOT NULL,
  `status` VARCHAR(20) DEFAULT 'pending',
  `description` TEXT DEFAULT NULL,
  `submit_location_lat` DECIMAL(10,7) DEFAULT NULL,
  `submit_location_lng` DECIMAL(10,7) DEFAULT NULL,
  `service_location_lat` DECIMAL(10,7) DEFAULT NULL,
  `service_location_lng` DECIMAL(10,7) DEFAULT NULL,
  `contact_name` VARCHAR(50) DEFAULT NULL,
  `contact_phone` VARCHAR(20) DEFAULT NULL,
  `address` VARCHAR(200) DEFAULT NULL,
  `appointment_time` DATETIME DEFAULT NULL,
  `urgency` VARCHAR(20) DEFAULT 'normal',
  `source_station_id` BIGINT DEFAULT NULL,
  `station_id` BIGINT DEFAULT NULL,
  `dispatch_basis` VARCHAR(50) DEFAULT NULL,
  `needs_manual_verify` TINYINT(1) DEFAULT 0,
  `reject_reason` TEXT DEFAULT NULL,
  `images` JSON DEFAULT NULL,
  `created_at` DATETIME DEFAULT NULL,
  `updated_at` DATETIME DEFAULT NULL,
  `deleted_at` DATETIME DEFAULT NULL,
  `rating` BIGINT DEFAULT NULL,
  `feedback` TEXT DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_service_requests_request_no` (`request_no`),
  KEY `idx_service_requests_user_id` (`user_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_service_requests_status` (`status`),
  KEY `idx_status` (`status`),
  KEY `idx_service_requests_source_station_id` (`source_station_id`),
  KEY `idx_service_requests_station_id` (`station_id`),
  KEY `idx_station_id` (`station_id`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_service_requests_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='服务需求表';

-- =====================================================
-- 7. 任务分配表 (task_assignments)
-- =====================================================
CREATE TABLE IF NOT EXISTS `task_assignments` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `request_id` BIGINT NOT NULL,
  `station_id` BIGINT NOT NULL,
  `staff_id` BIGINT DEFAULT NULL,
  `status` VARCHAR(20) DEFAULT 'dispatched',
  `claimed_at` DATETIME DEFAULT NULL,
  `completed_at` DATETIME DEFAULT NULL,
  `rating` BIGINT DEFAULT NULL,
  `feedback` TEXT DEFAULT NULL,
  `staff_notes` TEXT DEFAULT NULL,
  `images` JSON DEFAULT NULL,
  `created_at` DATETIME DEFAULT NULL,
  `updated_at` DATETIME DEFAULT NULL,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_request_id` (`request_id`),
  KEY `idx_task_assignments_request_id` (`request_id`),
  KEY `idx_station_id` (`station_id`),
  KEY `idx_task_assignments_station_id` (`station_id`),
  KEY `idx_staff_id` (`staff_id`),
  KEY `idx_task_assignments_staff_id` (`staff_id`),
  KEY `idx_status` (`status`),
  KEY `idx_task_assignments_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_task_assignments_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务分配表';

-- =====================================================
-- 8. 任务历史表 (task_histories)
-- =====================================================
CREATE TABLE IF NOT EXISTS `task_histories` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `task_id` BIGINT NOT NULL,
  `request_id` BIGINT NOT NULL,
  `action` VARCHAR(20) NOT NULL,
  `operator_id` BIGINT NOT NULL,
  `from_staff_id` BIGINT DEFAULT NULL,
  `to_staff_id` BIGINT DEFAULT NULL,
  `from_station_id` BIGINT DEFAULT NULL,
  `to_station_id` BIGINT DEFAULT NULL,
  `status_before` VARCHAR(20) DEFAULT NULL,
  `status_after` VARCHAR(20) DEFAULT NULL,
  `remark` TEXT DEFAULT NULL,
  `created_at` DATETIME DEFAULT NULL,
  `updated_at` DATETIME DEFAULT NULL,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_task_histories_task_id` (`task_id`),
  KEY `idx_task_id` (`task_id`),
  KEY `idx_request_id` (`request_id`),
  KEY `idx_task_histories_request_id` (`request_id`),
  KEY `idx_operator_id` (`operator_id`),
  KEY `idx_task_histories_operator_id` (`operator_id`),
  KEY `idx_from_staff_id` (`from_staff_id`),
  KEY `idx_task_histories_from_staff_id` (`from_staff_id`),
  KEY `idx_task_histories_to_staff_id` (`to_staff_id`),
  KEY `idx_to_staff_id` (`to_staff_id`),
  KEY `idx_from_station_id` (`from_station_id`),
  KEY `idx_task_histories_from_station_id` (`from_station_id`),
  KEY `idx_task_histories_to_station_id` (`to_station_id`),
  KEY `idx_to_station_id` (`to_station_id`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_task_histories_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务历史表';

-- =====================================================
-- 9. 通知表 (notifications)
-- =====================================================
CREATE TABLE IF NOT EXISTS `notifications` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` BIGINT NOT NULL,
  `title` VARCHAR(100) NOT NULL,
  `body` TEXT DEFAULT NULL,
  `type` VARCHAR(20) DEFAULT NULL,
  `related_id` BIGINT DEFAULT NULL,
  `related_type` VARCHAR(20) DEFAULT NULL,
  `channel` VARCHAR(20) NOT NULL,
  `send_status` VARCHAR(20) DEFAULT 'pending',
  `sent_at` DATETIME DEFAULT NULL,
  `is_read` TINYINT(1) DEFAULT NULL,
  `read_at` DATETIME DEFAULT NULL,
  `retry_count` BIGINT DEFAULT NULL,
  `created_at` DATETIME DEFAULT NULL,
  `updated_at` DATETIME DEFAULT NULL,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_notifications_user_id` (`user_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_notifications_type` (`type`),
  KEY `idx_type` (`type`),
  KEY `idx_notifications_related_id` (`related_id`),
  KEY `idx_related_id` (`related_id`),
  KEY `idx_channel` (`channel`),
  KEY `idx_notifications_channel` (`channel`),
  KEY `idx_notifications_send_status` (`send_status`),
  KEY `idx_send_status` (`send_status`),
  KEY `idx_is_read` (`is_read`),
  KEY `idx_notifications_is_read` (`is_read`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_notifications_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通知表';

-- =====================================================
-- 10. 新闻资讯表 (news)
-- =====================================================
CREATE TABLE IF NOT EXISTS `news` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `station_id` BIGINT DEFAULT NULL,
  `title` VARCHAR(200) NOT NULL,
  `summary` VARCHAR(500) DEFAULT NULL,
  `content` TEXT DEFAULT NULL,
  `cover_url` VARCHAR(255) DEFAULT NULL,
  `type` VARCHAR(20) DEFAULT 'news',
  `status` VARCHAR(20) DEFAULT 'draft',
  `author_id` BIGINT DEFAULT NULL,
  `publish_at` DATETIME DEFAULT NULL,
  `view_count` BIGINT DEFAULT NULL,
  `created_at` DATETIME DEFAULT NULL,
  `updated_at` DATETIME DEFAULT NULL,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_news_station_id` (`station_id`),
  KEY `idx_news_type` (`type`),
  KEY `idx_news_status` (`status`),
  KEY `idx_news_author_id` (`author_id`),
  KEY `idx_news_publish_at` (`publish_at`),
  KEY `idx_news_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='新闻资讯表';

-- =====================================================
-- 11. 轮播图表 (banners)
-- =====================================================
CREATE TABLE IF NOT EXISTS `banners` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `station_id` BIGINT NOT NULL COMMENT '站点ID(0=平台级/全局)',
  `title` VARCHAR(100) DEFAULT NULL COMMENT '标题',
  `image_url` VARCHAR(255) NOT NULL COMMENT '图片URL',
  `link_type` VARCHAR(20) NOT NULL DEFAULT 'none' COMMENT '跳转类型(none/news/url)',
  `link_value` VARCHAR(255) DEFAULT NULL COMMENT '跳转值(新闻ID或URL)',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序(数字越大越靠前)',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态(active/inactive)',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间(软删除)',
  PRIMARY KEY (`id`),
  KEY `idx_banners_station_status` (`station_id`, `status`),
  KEY `idx_banners_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='轮播图表';

-- =====================================================
-- 12. 统计报表表 (reports)
-- =====================================================
CREATE TABLE IF NOT EXISTS `reports` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  `name` VARCHAR(255) NOT NULL COMMENT '报表名称',
  `type` VARCHAR(32) NOT NULL COMMENT '报表类型：service/performance/request/station',
  `format` VARCHAR(16) NOT NULL COMMENT '文件格式：xlsx/csv',
  `file_path` VARCHAR(512) NOT NULL COMMENT '文件存储路径',
  `file_size` BIGINT NOT NULL DEFAULT 0 COMMENT '文件大小（字节）',
  `station_id` BIGINT DEFAULT NULL COMMENT '站点ID，NULL表示全局',
  `start_date` DATE NOT NULL COMMENT '统计开始日期',
  `end_date` DATE NOT NULL COMMENT '统计结束日期',
  `created_by` BIGINT NOT NULL COMMENT '创建人ID',
  PRIMARY KEY (`id`),
  KEY `idx_reports_type` (`type`),
  KEY `idx_reports_station_id` (`station_id`),
  KEY `idx_reports_created_at` (`created_at`),
  KEY `idx_reports_created_by` (`created_by`),
  KEY `idx_reports_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='统计报表表';

-- =====================================================
-- 13. 权限定义表 (permissions)
-- =====================================================
CREATE TABLE IF NOT EXISTS `permissions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `code` VARCHAR(100) NOT NULL COMMENT '权限码，格式: module:resource:action',
  `name` VARCHAR(100) NOT NULL COMMENT '权限名称',
  `description` VARCHAR(500) DEFAULT NULL COMMENT '权限描述',
  `module` VARCHAR(50) NOT NULL COMMENT '所属模块',
  `type` ENUM('menu','button','resource') NOT NULL DEFAULT 'resource' COMMENT '权限类型: menu-菜单 button-按钮 resource-资源',
  `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父权限ID，0表示顶级',
  `api_path` VARCHAR(200) DEFAULT NULL COMMENT 'API路径，支持通配符*',
  `api_method` VARCHAR(20) DEFAULT NULL COMMENT 'HTTP方法: GET/POST/PUT/DELETE',
  `is_public` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否公共权限，1-是 0-否',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态 active-启用 inactive-禁用',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序值',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` TIMESTAMP DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`, `deleted_at`),
  KEY `idx_module` (`module`),
  KEY `idx_type` (`type`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_api` (`api_path`, `api_method`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限定义表';

-- =====================================================
-- 14. 角色定义表 (roles)
-- =====================================================
CREATE TABLE IF NOT EXISTS `roles` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `code` VARCHAR(50) NOT NULL COMMENT '角色编码，如 admin, staff',
  `name` VARCHAR(100) NOT NULL COMMENT '角色名称',
  `description` VARCHAR(500) DEFAULT NULL COMMENT '角色描述',
  `is_system` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否系统内置角色，1-是 0-否',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态 active-启用 inactive-禁用',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序值',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` TIMESTAMP DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`, `deleted_at`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色定义表';

-- =====================================================
-- 15. 角色权限关联表 (role_permissions)
-- =====================================================
CREATE TABLE IF NOT EXISTS `role_permissions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `role_id` BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
  `permission_id` BIGINT UNSIGNED NOT NULL COMMENT '权限ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_perm` (`role_id`, `permission_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_permission_id` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- =====================================================
-- 16. 菜单配置表 (menus)
-- =====================================================
CREATE TABLE IF NOT EXISTS `menus` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父菜单ID，0表示顶级菜单',
  `name` VARCHAR(50) NOT NULL COMMENT '菜单名称',
  `path` VARCHAR(200) DEFAULT NULL COMMENT '路由路径',
  `component` VARCHAR(200) DEFAULT NULL COMMENT '前端组件路径',
  `icon` VARCHAR(50) DEFAULT NULL COMMENT '图标名称',
  `permission_code` VARCHAR(100) DEFAULT NULL COMMENT '权限码，格式: module:resource:action',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序值，数字越小越靠前',
  `hidden` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否隐藏 0-显示 1-隐藏',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态 active-启用 inactive-禁用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` TIMESTAMP DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_permission_code` (`permission_code`),
  KEY `idx_sort` (`sort`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单配置表';

-- =====================================================
-- 17. 数据库迁移记录表 (schema_migrations)
-- =====================================================
CREATE TABLE IF NOT EXISTS `schema_migrations` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `filename` VARCHAR(255) DEFAULT NULL,
  `applied_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_schema_migrations_filename` (`filename`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据库迁移记录表';
