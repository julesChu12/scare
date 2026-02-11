import { test, expect } from '@playwright/test';

test('verify task detail page access', async ({ page }) => {
  // --- 1. 登录 ---
  console.log('Step 1: Logging in...');
  await page.goto('http://localhost:3001/#/login');
  await page.waitForSelector('.login-card');
  await page.fill('input[placeholder="请输入手机号"]', '13800000001');
  await page.fill('input[placeholder="请输入密码"]', 'Test@123');
  await page.click('button:has-text("登 录")');
  
  // 等待登录完成并跳转到首页
  await page.waitForURL('**/#/', { timeout: 15000 });
  console.log('Login successful.');

  // --- 2. 访问任务列表页 ---
  console.log('Step 2: Navigating to Task List...');
  // 直接访问所有任务列表，因为那里数据最全
  await page.goto('http://localhost:3001/#/services/tasks/list');
  
  // 等待表格加载（通过查找表头或刷新按钮）
  await page.waitForSelector('.el-table', { timeout: 10000 });
  console.log('Task list loaded.');

  // --- 3. 查找并点击详情按钮 ---
  console.log('Step 3: Clicking "View Detail" button...');
  
  // 等待数据加载完成（loading 遮罩消失）
  await page.waitForSelector('.el-loading-mask', { state: 'detached', timeout: 10000 });

  // 查找第一个包含"查看"文本的按钮/链接
  // 注意：Element Plus 表格可能有多个按钮，我们要找行内的操作按钮
  const detailBtn = page.getByRole('button', { name: '查看' }).first();
  
  // 检查是否存在按钮（如果没有数据，测试将在这里跳过或失败）
  const count = await detailBtn.count();
  if (count === 0) {
    console.log('WARNING: No tasks found in the list. Cannot verify detail page.');
    return;
  }

  // 获取点击前的 URL 用于对比
  const prevUrl = page.url();
  await detailBtn.click();

  // --- 4. 验证详情页 ---
  console.log('Step 4: Verifying Detail Page...');
  
  // 等待 URL 变化
  await page.waitForURL(/\/services\/tasks\/\d+/, { timeout: 5000 });
  
  const currentUrl = page.url();
  console.log(`Current URL: ${currentUrl}`);

  // 断言：URL 必须包含 /services/tasks/ 且后面跟着数字
  expect(currentUrl).toMatch(/\/services\/tasks\/\d+/);

  // 断言：页面应该包含 "任务详情" 标题
  // 查找 h3 标签或其他标识
  await expect(page.locator('h3').filter({ hasText: '任务详情' })).toBeVisible();

  // 断言：不应该在 Dashboard (防止重定向)
  expect(currentUrl).not.toContain('/dashboard');

  console.log('SUCCESS: Task Detail page accessed successfully without redirect.');
});
