# C端前端设计 AI 提示词 (Prompt)

该文档包含用于指导 AI 设计和实现 sCare 项目 C 端（老年人端）页面的提示词，遵循适老化设计（Age-Friendly Design）原则。

---

## AI 提示词内容

```markdown
# Role
You are a Senior Frontend Engineer and UX Expert specializing in Accessibility (a11y) and Age-Friendly Design.

# Task
Design and implement the Core C-end (Consumer) "Home Page" and "Service Request Flow" for a Community Elderly Care application.

# Tech Stack
- Framework: Vue 3 (Script Setup)
- UI Library: Naive UI (Use NConfigProvider for theming)
- Styling: Tailwind CSS (preferred) or Scoped CSS
- Language: TypeScript
- Target Device: Mobile Web / PWA

# Design Requirements (Crucial for Elderly Users)
1. **Typography**: Base font size must be 18px+. Headings 24px+.
2. **Contrast**: Use high contrast colors. Primary action buttons should be very distinct (e.g., Deep Blue or Emerald Green).
3. **Touch Targets**: All clickable elements must be at least 48px height/width.
4. **Simplicity**: No complex menus. Use large cards with clear icons.

# Specific Components to Build

1. **ServiceGrid.vue (Home Page Component)**
   - Display a grid of 4-6 large cards for services:
     - Meal Delivery (Icon: Rice Bowl)
     - House Cleaning (Icon: Broom)
     - Medical Care (Icon: Medicine/Cross)
     - Companionship (Icon: Chat Bubbles)
   - Feature a prominent "Emergency Call" button at the top (Red color).

2. **RequestWizard.vue (Submission Form)**
   - A simplified form when a service is selected.
   - Field 1: **Location**. Default to "Current Location" (mock functionality) with a large button.
   - Field 2: **Time**. Two big toggle buttons: "Now" vs "Schedule".
   - Field 3: **Voice Note**. A large microphone button icon to simulate "Hold to Speak" for special instructions (mock logic).
   - **Submit Button**: A massive, full-width button fixed at the bottom.

# Output Format
Please provide the Vue component code for `HomeView.vue` (integrating the ServiceGrid) and `RequestView.vue` (the Wizard). Include the necessary script setup and dummy data.
Ensure the code is responsive and mobile-first.
```

---

## 使用说明

1. **复制提示词**：将上述 Markdown 块中的内容复制。
2. **粘贴至 AI 助手**：发送给 Cursor, Claude, ChatGPT 等 AI 助手。
3. **代码集成**：
   - 将生成的 `HomeView.vue` 和 `RequestView.vue` 放置在 `frontend/c-end/src/views/` 目录下。
   - 配置路由 `frontend/c-end/src/router/index.ts` 以访问这些页面。
   - 在 `App.vue` 中配置 `n-config-provider` 以统一字号和配色。
