import { test, expect } from '@playwright/test';

test.describe('Permission System Test (Real Backend)', () => {
  // Test Case 1: Admin should see everything
  test('Admin role should see system management menu', async ({ page }) => {
    // 1. Login (Admin)
    await page.goto('http://localhost:3001/login');
    await page.fill('input[placeholder="请输入手机号"]', '13800000001');
    await page.fill('input[placeholder="请输入密码"]', 'Test@123');
    await page.click('button:has-text("登 录")');
    await page.waitForURL('**/dashboard');

    // 2. Verify System Menu is visible
    // Wait for menu to load (since it comes from backend now)
    await page.waitForTimeout(1000); 
    const systemMenu = page.locator('.el-sub-menu__title', { hasText: '系统管理' });
    await expect(systemMenu).toBeVisible();
    
    // Verify user management (submenu item)
    // Need to expand menu first if it's collapsed, but usually Element Plus submenus are clickable
    await systemMenu.click();
    await expect(page.locator('.el-menu-item', { hasText: '用户管理' })).toBeVisible();
  });

  // Test Case 2: Staff should NOT see restricted menus
  test('Staff role should NOT see system management menu', async ({ page }) => {
    // 1. Login (Staff)
    await page.goto('http://localhost:3001/login');
    await page.fill('input[placeholder="请输入手机号"]', '13800000004'); // Staff phone
    await page.fill('input[placeholder="请输入密码"]', 'Test@123');
    await page.click('button:has-text("登 录")');
    await page.waitForURL('**/dashboard');

    // 2. Verify System Menu is NOT visible
    // Wait a bit to ensure menu is fully loaded
    await page.waitForTimeout(1000);
    const systemMenu = page.locator('.el-sub-menu__title', { hasText: '系统管理' });
    await expect(systemMenu).not.toBeVisible();

    // 3. Verify Direct Access Block (Router Guard)
    // Try to go to restricted page directly
    await page.goto('http://localhost:3001/system/users');
    
    // Should be redirected back or show error
    // In our implementation, permission guard redirects to '/' or previous page
    // We just check that we are NOT on the restricted page
    await expect(page).not.toHaveURL(/.*\/system\/users/);
  });
});
