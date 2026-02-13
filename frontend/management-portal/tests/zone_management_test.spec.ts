import { test, expect } from '@playwright/test';

test.describe('Zone Management Map Test (Real Backend)', () => {
  test('should load map container and canvas when adding a new zone', async ({ page }) => {
    // 1. Login
    await page.goto('http://localhost:3001/login');
    await page.fill('input[placeholder="请输入手机号"]', '13800000001');
    await page.fill('input[placeholder="请输入密码"]', 'Test@123');
    await page.click('button:has-text("登 录")');
    await page.waitForURL('**/dashboard');

    // 2. Navigate to Zone Management
    await page.goto('http://localhost:3001/stations/zones');
    
    // 3. Open "Add Zone" Dialog
    await page.click('button:has-text("新增围栏")');

    // 4. Verify Map Container
    const mapContainer = page.locator('.map-container');
    await expect(mapContainer).toBeVisible();
    
    // 5. Verify AMap Canvas Rendering
    // Wait for AMap to inject canvas elements (this proves JS API loaded and Key is valid)
    // Note: AMap injects multiple layers, we check if the container has content
    await page.waitForTimeout(3000); // Give it some time to render
    const html = await mapContainer.innerHTML();
    expect(html).toContain('amap-maps'); // AMap root class
    expect(html).toContain('amap-layers'); // AMap layers container
  });
});
