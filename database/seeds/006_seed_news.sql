-- =====================================================
-- 新闻数据种子文件
-- 版本: v1.0.0
-- 说明: 新闻内容
-- =====================================================

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

USE `scare_db`;

-- 清空现有数据
DELETE FROM `news` WHERE `deleted_at` IS NULL;

-- =====================================================
-- 新闻数据
-- =====================================================
INSERT INTO `news` (`title`, `summary`, `content`, `cover_url`, `type`, `status`, `station_id`, `author_id`, `publish_at`, `view_count`, `created_at`, `updated_at`) VALUES
(
  '社区养老服务中心正式启动，为老年人提供全方位照护',
  '本市首家综合性社区养老服务中心今日正式启动运营，将为周边社区老年人提供日间照料、健康管理、文娱活动等多项服务。',
  '<h2>服务中心简介</h2><p>社区养老服务中心占地面积1000平方米，设有休息区、活动室、健康检查室、餐厅等功能区域。中心配备专业护理人员20名，可同时为100位老年人提供服务。</p><h2>主要服务项目</h2><ul><li>日间照料：提供安全舒适的休息环境</li><li>健康管理：定期体检、慢病管理</li><li>营养配餐：根据老年人需求提供营养餐食</li><li>文娱活动：书法、绘画、太极拳等兴趣班</li><li>心理关怀：专业心理咨询服务</li></ul><h2>服务时间</h2><p>周一至周日 8:00-18:00，节假日正常开放。</p>',
  'https://via.placeholder.com/800x450/4A90E2/FFFFFF?text=养老服务中心',
  'news', 'published', 0, 1, NOW(), 156, NOW(3), NOW(3)
),
(
  '智慧养老新模式：一键呼叫，24小时守护老人安全',
  '我市推出智慧养老服务平台，老年人通过智能手环或手机APP即可一键呼叫紧急救援、家政服务、健康咨询等，真正实现"科技助老"。',
  '<h2>智慧养老平台功能</h2><p>智慧养老服务平台整合了社区、医疗、家政等多方资源，为老年人提供便捷的一站式服务。</p><h2>核心功能</h2><ul><li><strong>紧急呼叫</strong>：一键SOS，5分钟内响应</li><li><strong>健康监测</strong>：实时监测心率、血压等健康指标</li><li><strong>定位服务</strong>：防走失定位，家属随时掌握老人位置</li><li><strong>服务预约</strong>：在线预约家政、陪诊、送餐等服务</li><li><strong>亲情互动</strong>：视频通话、照片分享</li></ul><h2>使用方式</h2><p>老年人可通过智能手环、手机APP或固定电话使用平台服务。平台提供24小时客服热线：400-XXX-XXXX。</p>',
  'https://via.placeholder.com/800x450/50C878/FFFFFF?text=智慧养老',
  'news', 'published', 0, 1, DATE_SUB(NOW(), INTERVAL 3 DAY), 289, DATE_SUB(NOW(3), INTERVAL 3 DAY), DATE_SUB(NOW(3), INTERVAL 3 DAY)
),
(
  '关爱独居老人：社区志愿者定期上门探访服务启动',
  '为关爱社区独居老人，我市启动"温暖敲门"志愿服务项目，组织志愿者定期上门探访，了解老人需求，提供力所能及的帮助。',
  '<h2>项目背景</h2><p>据统计，我市60岁以上独居老人超过5万人。为让这些老人感受到社会关爱，市民政局联合社区启动"温暖敲门"志愿服务项目。</p><h2>服务内容</h2><ul><li>定期上门探访，了解老人生活状况</li><li>协助老人解决日常生活困难</li><li>陪伴聊天，缓解老人孤独感</li><li>发现问题及时反馈给社区和家属</li><li>普及防诈骗、安全用电等知识</li></ul><h2>志愿者招募</h2><p>项目面向全社会招募志愿者，要求有爱心、耐心，能定期参与服务。报名电话：XXX-XXXX-XXXX。</p>',
  'https://via.placeholder.com/800x450/FF6B6B/FFFFFF?text=志愿服务',
  'news', 'published', 0, 1, DATE_SUB(NOW(), INTERVAL 7 DAY), 423, DATE_SUB(NOW(3), INTERVAL 7 DAY), DATE_SUB(NOW(3), INTERVAL 7 DAY)
),
(
  '老年大学秋季班开始招生，丰富老年人精神文化生活',
  '市老年大学2024年秋季班即日起开始招生，开设书法、绘画、声乐、舞蹈、摄影、智能手机使用等20余门课程，欢迎老年朋友报名。',
  '<h2>课程设置</h2><p>本学期老年大学共开设20余门课程，涵盖艺术、健康、科技等多个领域。</p><h2>热门课程</h2><ul><li><strong>书法班</strong>：楷书、行书入门与提高</li><li><strong>国画班</strong>：花鸟、山水画技法</li><li><strong>声乐班</strong>：民族唱法、通俗唱法</li><li><strong>太极拳班</strong>：24式、42式太极拳</li><li><strong>智能手机班</strong>：微信使用、网上购物、健康码等</li><li><strong>摄影班</strong>：手机摄影技巧、照片后期处理</li></ul><h2>报名方式</h2><p>报名时间：即日起至9月15日<br>报名地点：市老年大学教务处（XX路XX号）<br>咨询电话：XXX-XXXX-XXXX<br>学费：每门课程200元/学期</p>',
  'https://via.placeholder.com/800x450/9B59B6/FFFFFF?text=老年大学',
  'news', 'published', 0, 1, DATE_SUB(NOW(), INTERVAL 10 DAY), 567, DATE_SUB(NOW(3), INTERVAL 10 DAY), DATE_SUB(NOW(3), INTERVAL 10 DAY)
),
(
  '医养结合新进展：社区卫生服务中心开设老年病专科',
  '为满足老年人医疗需求，多家社区卫生服务中心开设老年病专科，提供慢病管理、康复理疗、中医养生等服务，实现"小病不出社区"。',
  '<h2>老年病专科服务</h2><p>社区卫生服务中心老年病专科配备经验丰富的医护团队，为老年人提供专业的医疗服务。</p><h2>主要服务项目</h2><ul><li><strong>慢病管理</strong>：高血压、糖尿病等慢性病的规范化管理</li><li><strong>健康体检</strong>：老年人免费健康体检</li><li><strong>康复理疗</strong>：针灸、推拿、理疗等中医康复服务</li><li><strong>用药指导</strong>：合理用药咨询，避免药物不良反应</li><li><strong>家庭病床</strong>：为行动不便的老人提供上门医疗服务</li></ul><h2>就诊指南</h2><p>老年人可持医保卡到社区卫生服务中心就诊，享受医保报销政策。预约电话：XXX-XXXX-XXXX。</p>',
  'https://via.placeholder.com/800x450/3498DB/FFFFFF?text=医养结合',
  'news', 'published', 0, 1, DATE_SUB(NOW(), INTERVAL 14 DAY), 712, DATE_SUB(NOW(3), INTERVAL 14 DAY), DATE_SUB(NOW(3), INTERVAL 14 DAY)
),
(
  '适老化改造惠民生：为困难老人家庭免费安装安全设施',
  '市政府启动困难老人家庭适老化改造项目，为符合条件的老年人家庭免费安装扶手、防滑垫、紧急呼叫器等安全设施，预防老人跌倒等意外发生。',
  '<h2>项目介绍</h2><p>适老化改造项目旨在改善老年人居家生活环境，降低居家安全风险，提高老年人生活质量。</p><h2>改造内容</h2><ul><li><strong>卫生间改造</strong>：安装扶手、防滑垫、坐便椅</li><li><strong>卧室改造</strong>：安装床边扶手、夜间感应灯</li><li><strong>厨房改造</strong>：调整灶台高度、安装防烫设施</li><li><strong>客厅改造</strong>：清除地面障碍物、安装紧急呼叫器</li><li><strong>通道改造</strong>：加装扶手、改善照明</li></ul><h2>申请条件</h2><p>1. 年满60周岁<br>2. 属于低保、低收入或特困供养人员<br>3. 有适老化改造需求<br><br>申请方式：向所在社区提出申请，填写申请表。</p>',
  'https://via.placeholder.com/800x450/E74C3C/FFFFFF?text=适老化改造',
  'news', 'published', 0, 1, DATE_SUB(NOW(), INTERVAL 20 DAY), 834, DATE_SUB(NOW(3), INTERVAL 20 DAY), DATE_SUB(NOW(3), INTERVAL 20 DAY)
);

-- 验证数据
SELECT '新闻数据' AS 'Table', COUNT(*) AS 'Count' FROM `news` WHERE `deleted_at` IS NULL;
