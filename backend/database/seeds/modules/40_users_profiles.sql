-- 40_users_profiles.sql
-- 用户、身份、客户档案模块初始化

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- 密码哈希: $2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm (Test@123)
-- id_card_masked: 脱敏值直接存储，id_card/id_card_hmac 由 encrypt_seed 工具生成后 UPDATE
INSERT INTO `users` (`id`, `phone`, `password_hash`, `name`, `email`, `gender`, `birth_date`, `id_card_masked`, `station_id`, `status`, `created_at`, `updated_at`) VALUES
-- B端用户（9人）
(1, '13800000001', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '系统管理员', 'admin@scare.com', 'male', '1980-01-01', '1101**********0011', NULL, 'active', NOW(), NOW()),
(2, '13800000002', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '李站长', 'lizhang@scare.com', 'male', '1975-05-15', '1101**********0022', 1, 'active', NOW(), NOW()),
(3, '13800000003', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '王站长', 'wangzhang@scare.com', 'female', '1978-08-20', '1101**********0033', 2, 'active', NOW(), NOW()),
(4, '13800000004', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '王小红', 'xiaohong@scare.com', 'female', '1990-03-10', '1101**********0044', 1, 'active', NOW(), NOW()),
(5, '13800000005', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '刘师傅', 'liushifu@scare.com', 'male', '1985-07-22', '1101**********0055', 1, 'active', NOW(), NOW()),
(6, '13800000006', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '陈护士', 'chenhushi@scare.com', 'female', '1992-11-05', '1101**********0066', 2, 'active', NOW(), NOW()),
(7, '13800000007', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '赵大哥', 'zhaodage@scare.com', 'male', '1988-09-18', '1101**********0077', 2, 'active', NOW(), NOW()),
(24, '13800000024', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '孙小明', 'sunxm@scare.com', 'male', '1995-08-12', '1101**********0244', 1, 'active', NOW(), NOW()),
(25, '13800000025', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '周护工', 'zhouhg@scare.com', 'female', '1991-12-05', '1101**********0255', 2, 'active', NOW(), NOW()),

-- C端用户（16人）
(8, '13800000008', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '张大爷', NULL, 'male', '1950-05-15', '1101**********0088', NULL, 'active', NOW(), NOW()),
(9, '13800000009', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '李奶奶', NULL, 'female', '1955-03-20', '1101**********0099', 1, 'active', NOW(), NOW()),
(10, '13800000010', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '王爷爷', NULL, 'male', '1948-11-10', '1101**********0100', NULL, 'active', NOW(), NOW()),
(11, '13800000011', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '孙女士', NULL, 'female', '1990-06-25', '1101**********0111', NULL, 'active', NOW(), NOW()),
(12, '13800000012', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '赵先生', NULL, 'male', '1965-02-14', '1101**********0122', NULL, 'active', NOW(), NOW()),
(13, '13800000013', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '小明', NULL, 'male', '2018-03-15', '1101**********0133', NULL, 'active', NOW(), NOW()),
(14, '13800000014', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '周阿姨', NULL, 'female', '1949-03-08', '1101**********0144', NULL, 'active', NOW(), NOW()),
(15, '13800000015', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '吴大爷', NULL, 'male', '1951-07-12', '1101**********0155', NULL, 'active', NOW(), NOW()),
(16, '13800000016', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '郑奶奶', NULL, 'female', '1952-09-25', '1101**********0166', NULL, 'active', NOW(), NOW()),
(17, '13800000017', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '冯爷爷', NULL, 'male', '1946-05-18', '1101**********0177', NULL, 'active', NOW(), NOW()),
(18, '13800000018', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '陈阿姨', NULL, 'female', '1953-12-02', '1101**********0188', NULL, 'active', NOW(), NOW()),
(19, '13800000019', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '杨大爷', NULL, 'male', '1950-08-15', '1101**********0199', NULL, 'active', NOW(), NOW()),
(20, '13800000020', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '黄女士', NULL, 'female', '1992-05-10', '1101**********0200', NULL, 'active', NOW(), NOW()),
(21, '13800000021', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '林先生', NULL, 'male', '1970-03-20', '1101**********0211', NULL, 'active', NOW(), NOW()),
(22, '13800000022', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '何大爷', NULL, 'male', '1947-12-08', '1101**********0222', NULL, 'active', NOW(), NOW()),
(23, '13800000023', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '马奶奶', NULL, 'female', '1956-01-15', '1101**********0233', NULL, 'active', NOW(), NOW());

-- =====================================================
-- 4. 初始化用户身份（user_identities）
-- =====================================================
INSERT INTO `user_identities` (`user_id`, `identity_type`, `is_primary`, `station_id`, `status`, `granted_at`, `created_at`, `updated_at`) VALUES
-- B端身份（10条）
(1, 'admin', 1, NULL, 'active', NOW(), NOW(), NOW()),
(2, 'station_manager', 1, 1, 'active', NOW(), NOW(), NOW()),
(3, 'station_manager', 1, 2, 'active', NOW(), NOW(), NOW()),
(4, 'staff', 1, 1, 'active', NOW(), NOW(), NOW()),
(5, 'staff', 1, 1, 'active', NOW(), NOW(), NOW()),
(6, 'staff', 1, 2, 'active', NOW(), NOW(), NOW()),
(7, 'staff', 1, 2, 'active', NOW(), NOW(), NOW()),
(9, 'staff', 0, 1, 'active', NOW(), NOW(), NOW()),
(24, 'staff', 1, 1, 'active', NOW(), NOW(), NOW()),
(25, 'staff', 1, 2, 'active', NOW(), NOW(), NOW()),

-- C端身份（16条）
(8, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),
(9, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),
(10, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),
(11, 'pregnant', 1, NULL, 'active', NOW(), NOW(), NOW()),
(12, 'disabled', 1, NULL, 'active', NOW(), NOW(), NOW()),
(13, 'child', 1, NULL, 'active', NOW(), NOW(), NOW()),
(14, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),
(15, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),
(16, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),
(17, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),
(18, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),
(19, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),
(20, 'pregnant', 1, NULL, 'active', NOW(), NOW(), NOW()),
(21, 'family', 1, NULL, 'active', NOW(), NOW(), NOW()),
(22, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),
(23, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW());

-- =====================================================
-- 5. 初始化客户档案（C端服务对象）
-- =====================================================
INSERT INTO `customer_profiles` (
    `user_id`, `customer_type`, `gender`, `birth_date`, `address`,
    `health_status`, `disability_level`, `medical_history`, `special_needs`,
    `emergency_contact`, `created_at`, `updated_at`
) VALUES
(8, 'elderly', 'male', '1950-05-15', '北京市朝阳区幸福小区1号楼101', '良好', '自理',
 '高血压，需要定期测量血压', '每周需要社区医生上门测血压',
 '{"name":"张小明","phone":"13900000001","relation":"子女"}', NOW(), NOW()),

(9, 'elderly', 'female', '1955-03-20', '北京市朝阳区幸福小区2号楼202', '一般', '轻度失能',
 '糖尿病，行动不便', '需要助行器，定期血糖监测',
 '{"name":"李华","phone":"13900000002","relation":"子女"}', NOW(), NOW()),

(10, 'elderly', 'male', '1948-11-10', '北京市朝阳区幸福小区3号楼303', '较差', '中度失能',
 '心脏病，中风后遗症', '需要轮椅，24小时护理',
 '{"name":"王芳","phone":"13900000003","relation":"子女"}', NOW(), NOW()),

(11, 'pregnant', 'female', '1990-06-25', '北京市朝阳区康乐小区5号楼501', '良好', NULL,
 '孕27周，定期产检', '需要产前护理指导',
 '{"name":"孙先生","phone":"13900000004","relation":"配偶"}', NOW(), NOW()),

(12, 'disabled', 'male', '1965-02-14', '北京市朝阳区康乐小区6号楼602', '较差', '重度失能',
 '脊髓损伤，下肢瘫痪', '需要专业康复护理，定期更换导尿管',
 '{"name":"赵女士","phone":"13900000005","relation":"配偶"}', NOW(), NOW()),

(13, 'child', 'male', '2018-03-15', '北京市朝阳区幸福小区4号楼404', '良好', '自理',
 '无重大病史', '课后托管服务',
 '{"name":"小明妈妈","phone":"13900000006","relation":"母亲"}', NOW(), NOW()),

(14, 'elderly', 'female', '1949-03-08', '北京市朝阳区幸福小区5号楼505', '一般', '轻度失能',
 '关节炎，腰椎间盘突出', '需要定期理疗和康复训练',
 '{"name":"周明","phone":"13900000007","relation":"子女"}', NOW(), NOW()),

(15, 'elderly', 'male', '1951-07-12', '北京市朝阳区幸福小区6号楼606', '良好', '自理',
 '轻度白内障', '需要定期眼科检查',
 '{"name":"吴丽","phone":"13900000008","relation":"子女"}', NOW(), NOW()),

(16, 'elderly', 'female', '1952-09-25', '北京市朝阳区康乐小区1号楼101', '一般', '轻度失能',
 '骨质疏松，曾骨折', '需要防跌倒辅助和钙质补充',
 '{"name":"郑强","phone":"13900000009","relation":"子女"}', NOW(), NOW()),

(17, 'elderly', 'male', '1946-05-18', '北京市朝阳区康乐小区2号楼202', '较差', '重度失能',
 '阿尔茨海默症早期，高血压', '需要24小时看护，防走失',
 '{"name":"冯丽华","phone":"13900000010","relation":"子女"}', NOW(), NOW()),

(18, 'elderly', 'female', '1953-12-02', '北京市朝阳区幸福小区7号楼707', '一般', '自理',
 '慢性支气管炎', '需要定期呼吸功能检查',
 '{"name":"陈刚","phone":"13900000011","relation":"子女"}', NOW(), NOW()),

(19, 'elderly', 'male', '1950-08-15', '北京市朝阳区康乐小区3号楼303', '良好', '自理',
 '前列腺增生', '需要定期泌尿科复查',
 '{"name":"杨芳","phone":"13900000012","relation":"子女"}', NOW(), NOW()),

(20, 'pregnant', 'female', '1992-05-10', '北京市朝阳区幸福小区8号楼808', '良好', NULL,
 '孕32周，妊娠期糖尿病', '需要饮食指导和血糖监测',
 '{"name":"黄先生","phone":"13900000013","relation":"配偶"}', NOW(), NOW()),

(22, 'elderly', 'male', '1947-12-08', '北京市朝阳区康乐小区4号楼404', '较差', '中度失能',
 '帕金森病，行动迟缓', '需要康复训练和日常照护',
 '{"name":"何美","phone":"13900000014","relation":"子女"}', NOW(), NOW()),

(23, 'elderly', 'female', '1956-01-15', '北京市朝阳区幸福小区9号楼909', '一般', '自理',
 '高血脂，轻度脂肪肝', '需要饮食管理和定期体检',
 '{"name":"马军","phone":"13900000015","relation":"子女"}', NOW(), NOW());
