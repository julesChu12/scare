-- =====================================================
-- 新闻和轮播图数据种子文件
-- 版本: v3.0.0
-- 说明: 新闻内容 + 轮播图，覆盖各站点，补充时间戳
-- 注意: 图片使用相对路径 /static/，需配合 nginx 反向代理或后端静态文件服务
-- =====================================================

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

USE `scare_db`;

-- 清空现有数据
TRUNCATE TABLE `news`;
TRUNCATE TABLE `banners`;

-- =====================================================
-- 新闻数据（15条，覆盖各站点和类型）
-- =====================================================
INSERT INTO `news` (
  `title`, `summary`, `content`, `cover_url`, `type`, `status`,
  `station_id`, `author_id`, `publish_at`, `view_count`,
  `created_at`, `updated_at`, `deleted_at`
) VALUES
-- ========== 全局新闻（station_id=0，6条）==========
(
  '霍营街道社区养老服务正式上线',
  '霍营街道社区养老服务中心正式投入运营，为辖区老年人提供全方位养老服务。',
  '<h2>服务中心简介</h2><p>霍营街道社区养老服务中心占地面积1000平方米，设有休息区、活动室、健康检查室、餐厅等功能区域。中心配备专业护理人员和志愿者团队。</p><h2>主要服务</h2><ul><li>日间照料：为白天无人照料的老人提供专业照护服务</li><li>健康管理：定期体检、健康讲座、慢病管理</li><li>营养配餐：专业营养师搭配的老年餐食</li><li>文娱活动：书法、太极、手工等丰富的文化活动</li><li>紧急救助：24小时紧急响应机制</li></ul><h2>联系我们</h2><p>服务热线：010-12345678</p>',
  '/static/b_end/20260329/服务中心正式启动.png',
  'news', 'published', 0, 1,
  DATE_SUB(NOW(), INTERVAL 1 DAY), 1284,
  DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY), NULL
),
(
  '2026年助餐服务全面升级',
  '为方便辖区老年人用餐，现推出助餐服务全面升级，提供营养均衡的午餐和晚餐配送。',
  '<h2>助餐服务说明</h2><p>我们提供健康、营养的老年餐食，由专业营养师精心搭配，严格控制油盐摄入量。</p><h2>餐食标准</h2><ul><li>少油少盐，适合老年人消化系统</li><li>软烂易嚼，适合咀嚼困难老人</li><li>荤素搭配，营养均衡</li></ul><h2>订餐方式</h2><p>可通过电话、小程序或到站预约。</p>',
  '/static/b_end/20260329/上门送饭.jpg',
  'news', 'published', 0, 1,
  DATE_SUB(NOW(), INTERVAL 3 DAY), 986,
  DATE_SUB(NOW(), INTERVAL 3 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY), NULL
),
(
  '老年人健康知识讲座圆满结束',
  '近日举办的健康知识讲座受到老年人热烈欢迎，现场座无虚席，共有200余名老人参加。',
  '<h2>讲座内容</h2><p>本次讲座邀请知名老年医学专家讲解老年人健康知识，内容丰富实用。</p><h2>主要内容</h2><ul><li>高血压的预防与控制</li><li>糖尿病的日常管理</li><li>合理膳食与营养搭配</li><li>适量运动与康复锻炼</li><li>心理健康与情绪调节</li></ul><h2>讲座效果</h2><p>讲座结束后，许多老人表示受益匪浅，希望中心多举办此类活动。</p>',
  '/static/b_end/20260329/开会.jpg',
  'news', 'published', 0, 1,
  DATE_SUB(NOW(), INTERVAL 7 DAY), 743,
  DATE_SUB(NOW(), INTERVAL 7 DAY), DATE_SUB(NOW(), INTERVAL 7 DAY), NULL
),
(
  '社区志愿者培训活动成功举办',
  '为提升志愿者服务水平，特举办志愿服务培训活动，50余名志愿者参加培训。',
  '<h2>培训内容</h2><p>培训涵盖服务礼仪、急救知识、沟通技巧、老年心理学等方面。</p><h2>培训师资</h2><p>由专业社工和医护人员担任讲师。</p><h2>参与情况</h2><p>共有50余名志愿者参加培训，通过考核获得志愿服务证书。</p>',
  '/static/b_end/20260329/社区活动.jpg',
  'news', 'published', 0, 1,
  DATE_SUB(NOW(), INTERVAL 10 DAY), 568,
  DATE_SUB(NOW(), INTERVAL 10 DAY), DATE_SUB(NOW(), INTERVAL 10 DAY), NULL
),
(
  '便民服务进社区活动公告',
  '本周六将在小区广场举办便民服务活动，提供多项免费服务，欢迎广大居民参加。',
  '<h2>服务项目</h2><ul><li>理发服务：专业理发师免费理发</li><li>健康咨询：医生现场义诊</li><li>法律咨询：律师提供免费法律咨询</li><li>防诈骗宣传：民警讲解防骗知识</li><li>磨刀服务：传统磨刀手艺</li></ul><h2>活动时间</h2><p>周六 9:00-16:00</p><h2>活动地点</h2><p>华龙苑北里小区广场</p>',
  '/static/b_end/20260329/重要通知.jpg',
  'news', 'published', 0, 1,
  DATE_SUB(NOW(), INTERVAL 14 DAY), 892,
  DATE_SUB(NOW(), INTERVAL 14 DAY), DATE_SUB(NOW(), INTERVAL 14 DAY), NULL
),
(
  '养老服务政策解读专题',
  '最新养老服务政策解读，帮助居民了解相关福利政策，及时申请各项补贴。',
  '<h2>政策要点</h2><p>详细介绍国家及北京市养老服务相关政策。</p><h2>主要内容</h2><ul><li>居家养老补贴申请条件和流程</li><li>高龄津贴最新标准</li><li>养老机构优惠政策</li><li>长期护理保险制度介绍</li><li>老年优待证办理指南</li></ul>',
  '/static/b_end/20260329/服务中心正式启动1.png',
  'news', 'published', 0, 1,
  DATE_SUB(NOW(), INTERVAL 20 DAY), 1247,
  DATE_SUB(NOW(), INTERVAL 20 DAY), DATE_SUB(NOW(), INTERVAL 20 DAY), NULL
),

-- ========== 霍营站新闻（station_id=1，3条）==========
(
  '霍营站开展春季健康体检活动',
  '霍营街道养老服务中心开展春季健康体检活动，为辖区60岁以上老人提供免费体检。',
  '<h2>体检项目</h2><ul><li>血压、血糖检测</li><li>心电图检查</li><li>中医体质辨识</li><li>视力听力筛查</li><li>骨密度检测</li></ul><h2>体检时间</h2><p>3月15日-3月31日，周一至周五上午8:30-11:30</p><h2>预约方式</h2><p>电话预约或到站预约，预约电话：010-12345678</p>',
  '/static/b_end/20260329/开会.jpg',
  'news', 'published', 1, 2,
  DATE_SUB(NOW(), INTERVAL 5 DAY), 342,
  DATE_SUB(NOW(), INTERVAL 5 DAY), DATE_SUB(NOW(), INTERVAL 5 DAY), NULL
),
(
  '霍营站推出"暖心送餐"服务',
  '霍营街道养老服务中心推出暖心送餐服务，为行动不便的老人提供免费送餐上门服务。',
  '<h2>服务对象</h2><p>霍营街道辖区80岁以上独居老人，或行动不便的残疾老人</p><h2>服务内容</h2><ul><li>每日午餐和晚餐配送</li><li>营养师搭配的老年餐</li><li>免费送餐上门</li></ul><h2>申请条件</h2><p>持老年优待证，经评估后符合条件的老人可享受服务</p>',
  '/static/b_end/20260329/上门送饭.jpg',
  'news', 'published', 1, 2,
  DATE_SUB(NOW(), INTERVAL 12 DAY), 456,
  DATE_SUB(NOW(), INTERVAL 12 DAY), DATE_SUB(NOW(), INTERVAL 12 DAY), NULL
),
(
  '霍营站手工兴趣班开课通知',
  '霍营街道养老服务中心手工兴趣班本周开课，欢迎有兴趣的老人报名参加。',
  '<h2>课程内容</h2><p>剪纸、折纸、编织等传统手工艺</p><h2>上课时间</h2><p>每周三下午14:00-16:00</p><h2>上课地点</h2><p>霍营街道养老服务中心活动室</p><h2>报名方式</h2><p>电话报名或到站报名</p>',
  '/static/b_end/20260329/社区活动.jpg',
  'news', 'published', 1, 2,
  DATE_SUB(NOW(), INTERVAL 18 DAY), 198,
  DATE_SUB(NOW(), INTERVAL 18 DAY), DATE_SUB(NOW(), INTERVAL 18 DAY), NULL
),

-- ========== 龙泽站新闻（station_id=2，3条）==========
(
  '龙泽站义诊活动圆满成功',
  '龙泽园养老服务站联合昌平区医院举办义诊活动，为300余名老人提供免费医疗服务。',
  '<h2>义诊项目</h2><ul><li>内科常规检查</li><li>血糖血压检测</li><li>中医问诊</li><li>用药咨询</li></ul><h2>参与情况</h2><p>活动共服务300余名老人，发放健康宣传资料500余份。</p>',
  '/static/b_end/20260329/开会.jpg',
  'news', 'published', 2, 3,
  DATE_SUB(NOW(), INTERVAL 8 DAY), 521,
  DATE_SUB(NOW(), INTERVAL 8 DAY), DATE_SUB(NOW(), INTERVAL 8 DAY), NULL
),
(
  '龙泽站太极班招生中',
  '龙泽园养老服务站太极拳班开始招生，专业教练指导，欢迎参加。',
  '<h2>课程介绍</h2><p>由陈式太极拳传人亲自指导，动作舒缓，适合老年人锻炼。</p><h2>上课时间</h2><p>每周二、周五上午7:00-8:30</p><h2>上课地点</h2><p>龙泽苑社区广场</p><h2>收费标准</h2><p>免费公益课程</p>',
  '/static/b_end/20260329/社区活动.jpg',
  'news', 'published', 2, 3,
  DATE_SUB(NOW(), INTERVAL 15 DAY), 267,
  DATE_SUB(NOW(), INTERVAL 15 DAY), DATE_SUB(NOW(), INTERVAL 15 DAY), NULL
),
(
  '龙泽站推出认知症关爱计划',
  '龙泽园养老服务站推出认知症（老年痴呆）关爱计划，为辖区认知症老人家庭提供支持。',
  '<h2>服务内容</h2><ul><li>认知症筛查评估</li><li>照护者培训支持</li><li>记忆训练活动</li><li>家庭支援服务</li></ul><h2>服务对象</h2><p>辖区认知症老人及其照护家庭</p>',
  '/static/b_end/20260329/服务中心正式启动.png',
  'news', 'published', 2, 3,
  DATE_SUB(NOW(), INTERVAL 22 DAY), 389,
  DATE_SUB(NOW(), INTERVAL 22 DAY), DATE_SUB(NOW(), INTERVAL 22 DAY), NULL
),

-- ========== 公告类型（2条）==========
(
  '清明节放假通知',
  '根据国家规定，2026年清明节放假安排如下，请各位居民知悉。',
  '<h2>放假安排</h2><p>4月4日（周六）至4月6日（周一）放假调休，共3天。</p><h2>服务安排</h2><ul><li>4月4日（周六）：站点休息，不提供上门服务</li><li>4月5日（周日）：部分服务正常运营</li><li>4月6日（周一）：恢复正常服务</li></ul><h2>紧急联系</h2><p>紧急情况请拨打：010-12345678</p>',
  '/static/b_end/20260329/重要通知.jpg',
  'notice', 'published', 0, 1,
  DATE_SUB(NOW(), INTERVAL 4 DAY), 2134,
  DATE_SUB(NOW(), INTERVAL 4 DAY), DATE_SUB(NOW(), INTERVAL 4 DAY), NULL
),
(
  '系统升级维护通知',
  '养老服务系统将于本周日凌晨进行升级维护，届时部分功能暂停使用。',
  '<h2>维护时间</h2><p>本周日凌晨2:00-6:00，预计4小时</p><h2>影响范围</h2><ul><li>小程序预约功能暂停</li><li>在线支付功能暂停</li><li>其他功能不受影响</li></ul><h2>应急服务</h2><p>维护期间紧急求助请拨打服务热线：010-12345678</p>',
  '/static/b_end/20260329/重要通知.jpg',
  'notice', 'published', 0, 1,
  DATE_SUB(NOW(), INTERVAL 2 DAY), 876,
  DATE_SUB(NOW(), INTERVAL 2 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY), NULL
),

-- ========== 草稿箱（1条，未发布）==========
(
  '端午节活动策划（草稿）',
  '端午节活动初步策划案，计划组织包粽子等传统文化活动。',
  '<h2>活动初步计划</h2><p>端午节组织社区老人包粽子活动。</p><h2>待确认事项</h2><ul><li>活动时间地点待定</li><li>材料预算待审批</li></ul>',
  '/static/b_end/20260329/社区活动.jpg',
  'news', 'draft', 0, 1,
  DATE_SUB(NOW(), INTERVAL 25 DAY), 12,
  DATE_SUB(NOW(), INTERVAL 25 DAY), DATE_SUB(NOW(), INTERVAL 25 DAY), NULL
);

-- =====================================================
-- 轮播图数据（8条）
-- =====================================================
INSERT INTO `banners` (
  `station_id`, `title`, `image_url`, `link_type`, `link_value`, `sort`, `status`,
  `created_at`, `updated_at`, `deleted_at`
) VALUES
(0, '社区养老服务正式上线', '/static/b_end/20260329/服务中心正式启动.png', 'none', NULL, 10, 'active',
 DATE_SUB(NOW(), INTERVAL 30 DAY), DATE_SUB(NOW(), INTERVAL 30 DAY), NULL),
(0, '暖心助餐服务预约中', '/static/b_end/20260329/上门送饭.jpg', 'none', NULL, 9, 'active',
 DATE_SUB(NOW(), INTERVAL 30 DAY), DATE_SUB(NOW(), INTERVAL 30 DAY), NULL),
(0, '春季健康义诊活动', '/static/b_end/20260329/开会.jpg', 'none', NULL, 8, 'active',
 DATE_SUB(NOW(), INTERVAL 30 DAY), DATE_SUB(NOW(), INTERVAL 30 DAY), NULL),
(0, '志愿者招募进行中', '/static/b_end/20260329/社区活动.jpg', 'none', NULL, 7, 'active',
 DATE_SUB(NOW(), INTERVAL 30 DAY), DATE_SUB(NOW(), INTERVAL 30 DAY), NULL),
(0, '便民服务进社区', '/static/b_end/20260329/上门理发.webp', 'none', NULL, 6, 'active',
 DATE_SUB(NOW(), INTERVAL 30 DAY), DATE_SUB(NOW(), INTERVAL 30 DAY), NULL),
-- 霍营站专属轮播
(1, '霍营站健康体检预约', '/static/b_end/20260329/开会.jpg', 'none', NULL, 5, 'active',
 DATE_SUB(NOW(), INTERVAL 10 DAY), DATE_SUB(NOW(), INTERVAL 10 DAY), NULL),
(1, '霍营站暖心送餐', '/static/b_end/20260329/上门送饭.jpg', 'none', NULL, 4, 'active',
 DATE_SUB(NOW(), INTERVAL 15 DAY), DATE_SUB(NOW(), INTERVAL 15 DAY), NULL),
-- 龙泽站专属轮播
(2, '龙泽站义诊活动', '/static/b_end/20260329/服务中心正式启动.png', 'none', NULL, 5, 'active',
 DATE_SUB(NOW(), INTERVAL 12 DAY), DATE_SUB(NOW(), INTERVAL 12 DAY), NULL);

-- 验证数据
SELECT '新闻数据' AS `Table`, COUNT(*) AS `Count` FROM `news` WHERE `deleted_at` IS NULL
UNION ALL SELECT '轮播图数据', COUNT(*) FROM `banners` WHERE `deleted_at` IS NULL;
