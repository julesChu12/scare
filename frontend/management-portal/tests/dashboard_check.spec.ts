import { test, expect } from '@playwright/test';

test('verify dashboard layout', async ({ page }) => {
  // 1. 登录
  await page.goto('http://localhost:3001/login');
  await page.fill('input[placeholder="请输入手机号"]', '13800000001');
  await page.fill('input[placeholder="请输入密码"]', 'Test@123');
  await page.click('button:has-text("登 录")');
  await page.waitForURL('**/dashboard');

  // 2. 等待加载
  await page.waitForTimeout(2000);

  // 3. 检查核心区域是否存在
  const taskCard = page.locator('.task-list-card');
  const actionCard = page.locator('.action-card');
  
  await expect(taskCard).toBeVisible();
  await expect(actionCard).toBeVisible();

  // 4. 检查具体内容
  // 左侧标题
  await expect(taskCard.locator('text=最新待处理任务')).toBeVisible();
  // 右侧快捷操作
  await expect(actionCard.locator('text=快捷操作')).toBeVisible();
  await expect(actionCard.locator('text=任务池')).toBeVisible();
  
  // 5. 截图留证
  await page.screenshot({ path: 'test-results/dashboard-layout.png', fullPage: true });
});
