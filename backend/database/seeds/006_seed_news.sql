-- =====================================================
-- 新闻数据种子文件
-- 版本: v2.0.0
-- 说明: 新闻内容 + 轮播图
-- 注意: 图片使用相对路径 /static/，需配合 nginx 反向代理或后端静态文件服务
-- =====================================================

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

USE `scare_db`;

-- 清空现有数据
TRUNCATE TABLE `news`;
TRUNCATE TABLE `banners`;

-- =====================================================
-- 新闻数据
-- =====================================================
INSERT INTO `news` (`title`, `summary`, `content`, `cover_url`, `type`, `status`, `station_id`, `author_id`, `publish_at`, `view_count`, `created_at`, `updated_at`) VALUES
(
  '霍营街道社区养老服务正式上线',
  '霍营街道社区养老服务中心正式投入运营，为辖区老年人提供全方位养老服务。',
  '<h2>服务中心简介</h2><p>霍营街道社区养老服务中心占地面积1000平方米，设有休息区、活动室、健康检查室、餐厅等功能区域。</p><h2>主要服务</h2><ul><li>日间照料</li><li>健康管理</li><li>营养配餐</li><li>文娱活动</li></ul>',
  '/static/b_end/20260329/服务中心正式启动.png',
  'news', 'published', 0, 1, NOW(), 156, NOW(3), NOW(3)
),
(
  '2024年助餐服务安排',
  '为方便辖区老年人用餐，现推出助餐服务，提供营养均衡的午餐配送。',
  '<h2>助餐服务说明</h2><p>我们提供健康、营养的老年餐食，由专业营养师搭配。</p><h2>订餐方式</h2><p>可通过电话或小程序预约。</p>',
  '/static/b_end/20260329/上门送饭.jpg',
  'news', 'published', 0, 1, DATE_SUB(NOW(), INTERVAL 3 DAY), 289, DATE_SUB(NOW(3), INTERVAL 3 DAY), DATE_SUB(NOW(3), INTERVAL 3 DAY)
),
(
  '老年人健康讲座圆满结束',
  '近日举办的健康知识讲座受到老年人热烈欢迎，现场座无虚席。',
  '<h2>讲座内容</h2><p>本次讲座邀请知名专家讲解老年人健康知识。</p><h2>主要内容</h2><ul><li>高血压防治</li><li>糖尿病管理</li><li>合理膳食</li></ul>',
  '/static/b_end/20260329/开会.jpg',
  'news', 'published', 0, 1, DATE_SUB(NOW(), INTERVAL 7 DAY), 423, DATE_SUB(NOW(3), INTERVAL 7 DAY), DATE_SUB(NOW(3), INTERVAL 7 DAY)
),
(
  '社区志愿者培训活动回顾',
  '为提升志愿者服务水平，特举办志愿服务培训活动。',
  '<h2>培训内容</h2><p>培训涵盖服务礼仪、急救知识、沟通技巧等方面。</p><h2>参与情况</h2><p>共有50余名志愿者参加培训。</p>',
  '/static/b_end/20260329/社区活动.jpg',
  'news', 'published', 0, 1, DATE_SUB(NOW(), INTERVAL 10 DAY), 567, DATE_SUB(NOW(3), INTERVAL 10 DAY), DATE_SUB(NOW(3), INTERVAL 10 DAY)
),
(
  '便民服务进社区活动',
  '本周六将在小区广场举办便民服务活动，提供多项免费服务。',
  '<h2>服务项目</h2><ul><li>理发服务</li><li>健康咨询</li><li>法律咨询</li><li>防诈骗宣传</li></ul><h2>活动时间</h2><p>周六 9:00-16:00</p>',
  '/static/b_end/20260329/重要通知.jpg',
  'news', 'published', 0, 1, DATE_SUB(NOW(), INTERVAL 14 DAY), 712, DATE_SUB(NOW(3), INTERVAL 14 DAY), DATE_SUB(NOW(3), INTERVAL 14 DAY)
),
(
  '养老服务政策解读',
  '最新养老服务政策解读，帮助居民了解相关福利政策。',
  '<h2>政策要点</h2><p>详细介绍国家及北京市养老服务相关政策。</p><h2>主要内容</h2><ul><li>居家养老补贴</li><li>高龄津贴申请</li><li>养老机构优惠政策</li></ul>',
  '/static/b_end/20260329/服务中心正式启动1.png',
  'news', 'published', 0, 1, DATE_SUB(NOW(), INTERVAL 20 DAY), 834, DATE_SUB(NOW(3), INTERVAL 20 DAY), DATE_SUB(NOW(3), INTERVAL 20 DAY)
);

-- =====================================================
-- 轮播图数据
-- =====================================================
INSERT INTO `banners` (`station_id`, `title`, `image_url`, `link_type`, `link_value`, `sort`, `status`) VALUES
(0, '社区养老服务上线', '/static/b_end/20260329/服务中心正式启动.png', 'none', NULL, 10, 'active'),
(0, '助餐服务预约', '/static/b_end/20260329/上门送饭.jpg', 'none', NULL, 9, 'active'),
(0, '健康义诊活动', '/static/b_end/20260329/开会.jpg', 'none', NULL, 8, 'active'),
(0, '志愿者招募', '/static/b_end/20260329/社区活动.jpg', 'none', NULL, 7, 'active'),
(0, '便民服务日', '/static/b_end/20260329/上门理发.webp', 'none', NULL, 6, 'active');

-- 验证数据
SELECT '新闻数据' AS 'Table', COUNT(*) AS 'Count' FROM `news` WHERE `deleted_at` IS NULL
UNION ALL SELECT '轮播图数据', COUNT(*) FROM `banners` WHERE `deleted_at` IS NULL;
