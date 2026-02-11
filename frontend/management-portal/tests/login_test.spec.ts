import { test, expect } from '@playwright/test';

test('admin login test', async ({ page }) => {
  // 1. 访问登录页
  console.log('Navigating to login page...');
  await page.goto('http://localhost:3001/#/login');

  // 2. 等待页面加载
  await page.waitForSelector('.login-card');

  // 3. 填充账号密码
  console.log('Filling credentials...');
  await page.fill('input[placeholder="请输入手机号"]', '13800000001');
  await page.fill('input[placeholder="请输入密码"]', 'Test@123');

  // 4. 点击登录按钮
  console.log('Clicking login button...');
  await page.click('button:has-text("登 录")');

  // 5. 验证是否登录成功
  // 登录成功后通常会跳转到首页，或者显示"登录成功"的消息
  // 我们可以检查 URL 是否改变，或者页面上是否出现了 Dashboard 相关的元素
  
  // 等待 URL 变化 (假设跳转到 /dashboard 或 /)
  await page.waitForURL('**/#/', { timeout: 10000 });
  console.log('Redirected to home page.');

  // 也可以检查 localStorage 中是否有 token
  const token = await page.evaluate(() => localStorage.getItem('b_token'));
  expect(token).toBeTruthy();
  console.log('Token found in localStorage.');
});
