-- =====================================================
-- 统计报表数据种子文件
-- 版本: v1.0.0
-- 说明: 系统统计报表数据
-- =====================================================

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

USE `scare_db`;

-- 清空现有数据
TRUNCATE TABLE `reports`;

-- =====================================================
-- 报表数据
-- type: service(服务统计) / performance(绩效统计) / request(需求统计) / station(站点统计)
-- format: xlsx / csv
-- station_id: NULL 表示全局报表
-- =====================================================
INSERT INTO `reports` (
  `name`, `type`, `format`, `file_path`, `file_size`,
  `station_id`, `start_date`, `end_date`, `created_by`, `created_at`, `updated_at`
) VALUES
-- ========== 全局服务统计报表 ==========
(
  '2026年3月服务统计月报',
  'service', 'xlsx',
  '/static/reports/2026_03_service_monthly.xlsx',
  245760,
  NULL,
  '2026-03-01', '2026-03-31', 1,
  DATE_SUB(NOW(), INTERVAL 2 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY)
),
(
  '2026年2月服务统计月报',
  'service', 'xlsx',
  '/static/reports/2026_02_service_monthly.xlsx',
  229376,
  NULL,
  '2026-02-01', '2026-02-28', 1,
  DATE_SUB(NOW(), INTERVAL 32 DAY), DATE_SUB(NOW(), INTERVAL 32 DAY)
),
(
  '2026年1月服务统计月报',
  'service', 'xlsx',
  '/static/reports/2026_01_service_monthly.xlsx',
  212992,
  NULL,
  '2026-01-01', '2026-01-31', 1,
  DATE_SUB(NOW(), INTERVAL 62 DAY), DATE_SUB(NOW(), INTERVAL 62 DAY)
),

-- ========== 站点统计报表 ==========
(
  '霍营站2026年3月站点统计',
  'station', 'xlsx',
  '/static/reports/station_1_2026_03.xlsx',
  102400,
  1,
  '2026-03-01', '2026-03-31', 2,
  DATE_SUB(NOW(), INTERVAL 2 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY)
),
(
  '龙泽站2026年3月站点统计',
  'station', 'xlsx',
  '/static/reports/station_2_2026_03.xlsx',
  94208,
  2,
  '2026-03-01', '2026-03-31', 3,
  DATE_SUB(NOW(), INTERVAL 2 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY)
),
(
  '回龙观站2026年3月站点统计',
  'station', 'xlsx',
  '/static/reports/station_3_2026_03.xlsx',
  76800,
  3,
  '2026-03-01', '2026-03-31', 3,
  DATE_SUB(NOW(), INTERVAL 2 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY)
),

-- ========== 绩效统计报表 ==========
(
  '霍营站工作人员2026年3月绩效考核',
  'performance', 'xlsx',
  '/static/reports/station_1_performance_2026_03.xlsx',
  184320,
  1,
  '2026-03-01', '2026-03-31', 2,
  DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY)
),
(
  '龙泽站工作人员2026年3月绩效考核',
  'performance', 'xlsx',
  '/static/reports/station_2_performance_2026_03.xlsx',
  163840,
  2,
  '2026-03-01', '2026-03-31', 3,
  DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY)
),

-- ========== 需求统计报表 ==========
(
  '2026年3月需求统计周报（第1周）',
  'request', 'csv',
  '/static/reports/2026_03_week1_requests.csv',
  15360,
  NULL,
  '2026-03-01', '2026-03-07', 1,
  DATE_SUB(NOW(), INTERVAL 12 DAY), DATE_SUB(NOW(), INTERVAL 12 DAY)
),
(
  '2026年3月需求统计周报（第2周）',
  'request', 'csv',
  '/static/reports/2026_03_week2_requests.csv',
  17408,
  NULL,
  '2026-03-08', '2026-03-14', 1,
  DATE_SUB(NOW(), INTERVAL 5 DAY), DATE_SUB(NOW(), INTERVAL 5 DAY)
),
(
  '2026年3月需求统计周报（第3周）',
  'request', 'csv',
  '/static/reports/2026_03_week3_requests.csv',
  18432,
  NULL,
  '2026-03-15', '2026-03-21', 1,
  DATE_SUB(NOW(), INTERVAL 2 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY)
),

-- ========== 历史季度报表 ==========
(
  '2026年Q1季度服务汇总报告',
  'service', 'xlsx',
  '/static/reports/2026_Q1_quarterly.xlsx',
  524288,
  NULL,
  '2026-01-01', '2026-03-31', 1,
  DATE_SUB(NOW(), INTERVAL 15 DAY), DATE_SUB(NOW(), INTERVAL 15 DAY)
);

-- 验证数据
SELECT '统计报表' AS `Table`, COUNT(*) AS `Count` FROM `reports`;
