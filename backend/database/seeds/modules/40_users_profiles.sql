-- 40_users_profiles.sql
-- 用户、身份、客户档案模块初始化

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- 密码哈希: $2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm
INSERT INTO `users` (`id`, `phone`, `password_hash`, `name`, `email`, `gender`, `birth_date`, `station_id`, `status`, `created_at`, `updated_at`) VALUES
-- B端用户
(1, '13800000001', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '系统管理员', 'admin@scare.com', 'male', '1980-01-01', NULL, 'active', NOW(), NOW()),
(2, '13800000002', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '李站长', 'lizhang@scare.com', 'male', '1975-05-15', 1, 'active', NOW(), NOW()),
(3, '13800000003', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '王站长', 'wangzhang@scare.com', 'female', '1978-08-20', 2, 'active', NOW(), NOW()),
(4, '13800000004', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '王小红', 'xiaohong@scare.com', 'female', '1990-03-10', 1, 'active', NOW(), NOW()),
(5, '13800000005', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '刘师傅', 'liushifu@scare.com', 'male', '1985-07-22', 1, 'active', NOW(), NOW()),
(6, '13800000006', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '陈护士', 'chenhushi@scare.com', 'female', '1992-11-05', 2, 'active', NOW(), NOW()),
(7, '13800000007', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '赵大哥', 'zhaodage@scare.com', 'male', '1988-09-18', 2, 'active', NOW(), NOW()),

-- C端用户（服务对象）
(8, '13800000008', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '张大爷', NULL, 'male', '1950-05-15', NULL, 'active', NOW(), NOW()),
(9, '13800000009', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '李奶奶', NULL, 'female', '1955-03-20', 1, 'active', NOW(), NOW()),  -- 跨端用户（员工+客户）
(10, '13800000010', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '王爷爷', NULL, 'male', '1948-11-10', NULL, 'active', NOW(), NOW()),
(11, '13800000011', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '孙女士', NULL, 'female', '1990-06-25', NULL, 'active', NOW(), NOW()),  -- 孕妇
(12, '13800000012', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '赵先生', NULL, 'male', '1965-02-14', NULL, 'active', NOW(), NOW()),  -- 失能人士
(13, '13800000013', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '小明', NULL, 'male', '2018-03-15', NULL, 'active', NOW(), NOW());  -- 儿童

-- =====================================================
-- 4. 初始化用户身份（user_identities）
-- B端身份：admin, station_manager, staff
-- C端身份：elderly, family, pregnant, disabled, child
-- =====================================================
INSERT INTO `user_identities` (`user_id`, `identity_type`, `is_primary`, `station_id`, `status`, `granted_at`, `created_at`, `updated_at`) VALUES
-- B端身份（8条）
(1, 'admin', 1, NULL, 'active', NOW(), NOW(), NOW()),           -- 系统管理员
(2, 'station_manager', 1, 1, 'active', NOW(), NOW(), NOW()),    -- 李站长 - 站点1
(3, 'station_manager', 1, 2, 'active', NOW(), NOW(), NOW()),    -- 王站长 - 站点2
(4, 'staff', 1, 1, 'active', NOW(), NOW(), NOW()),              -- 王小红 - 站点1
(5, 'staff', 1, 1, 'active', NOW(), NOW(), NOW()),              -- 刘师傅 - 站点1
(6, 'staff', 1, 2, 'active', NOW(), NOW(), NOW()),              -- 陈护士 - 站点2
(7, 'staff', 1, 2, 'active', NOW(), NOW(), NOW()),              -- 赵大哥 - 站点2
(9, 'staff', 0, 1, 'active', NOW(), NOW(), NOW()),              -- 李奶奶的副身份（志愿者）

-- C端身份（6条）
(8, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),         -- 张大爷 - 老年人
(9, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),         -- 李奶奶 - 老年人（主身份）
(10, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),        -- 王爷爷 - 老年人
(11, 'pregnant', 1, NULL, 'active', NOW(), NOW(), NOW()),       -- 孙女士 - 孕妇
(12, 'disabled', 1, NULL, 'active', NOW(), NOW(), NOW()),       -- 赵先生 - 失能人士
(13, 'child', 1, NULL, 'active', NOW(), NOW(), NOW());          -- 小明 - 儿童

-- =====================================================
-- 5. 初始化客户档案（C端服务对象）
-- =====================================================
INSERT INTO `customer_profiles` (
    `user_id`,
    `customer_type`,
    `gender`,
    `birth_date`,
    `address`,
    `health_status`,
    `disability_level`,
    `medical_history`,
    `special_needs`,
    `emergency_contact`,
    `created_at`,
    `updated_at`
) VALUES
-- 张大爷：老年人
(8, 'elderly', 'male', '1950-05-15', '北京市朝阳区幸福小区1号楼101', '良好', '自理',
 '高血压，需要定期测量血压', '每周需要社区医生上门测血压',
 '{"name":"张小明","phone":"13900000001","relation":"子女"}', NOW(), NOW()),

-- 李奶奶：老年人（跨端用户）
(9, 'elderly', 'female', '1955-03-20', '北京市朝阳区幸福小区2号楼202', '一般', '轻度失能',
 '糖尿病，行动不便', '需要助行器，定期血糖监测',
 '{"name":"李华","phone":"13900000002","relation":"子女"}', NOW(), NOW()),

-- 王爷爷：老年人
(10, 'elderly', 'male', '1948-11-10', '北京市朝阳区幸福小区3号楼303', '较差', '中度失能',
 '心脏病，中风后遗症', '需要轮椅，24小时护理',
 '{"name":"王芳","phone":"13900000003","relation":"子女"}', NOW(), NOW()),

-- 孙女士：孕妇
(11, 'pregnant', 'female', '1990-06-25', '北京市朝阳区康乐小区5号楼501', '良好', NULL,
 '孕27周，定期产检', '需要产前护理指导',
 '{"name":"孙先生","phone":"13900000004","relation":"配偶"}', NOW(), NOW()),

-- 赵先生：失能人士
(12, 'disabled', 'male', '1965-02-14', '北京市朝阳区康乐小区6号楼602', '较差', '重度失能',
 '脊髓损伤，下肢瘫痪', '需要专业康复护理，定期更换导尿管',
 '{"name":"赵女士","phone":"13900000005","relation":"配偶"}', NOW(), NOW()),

-- 小明：儿童
(13, 'child', 'male', '2018-03-15', '北京市朝阳区幸福小区4号楼404', '良好', '自理',
 '无重大病史', '课后托管服务',
 '{"name":"小明妈妈","phone":"13900000006","relation":"母亲"}', NOW(), NOW());
