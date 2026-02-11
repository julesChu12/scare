import { ref, watch, onMounted } from 'vue'

export type FontSize = 'small' | 'medium' | 'large'

const FONT_SIZE_KEY = 'c_font_size'

// 字体大小配置
export const fontSizeConfig = {
  small: { label: '小', value: '14px' },
  medium: { label: '中', value: '16px' },
  large: { label: '大', value: '18px' }
}

export function useFontSize() {
  // 从 localStorage 读取初始值
  const savedSize = localStorage.getItem(FONT_SIZE_KEY) as FontSize | null
  const fontSize = ref<FontSize>(savedSize || 'medium')

  // 应用字体大小到 document
  const applyFontSize = (size: FontSize) => {
    document.documentElement.setAttribute('data-font-size', size)
  }

  // 设置字体大小
  const setFontSize = (size: FontSize) => {
    fontSize.value = size
    localStorage.setItem(FONT_SIZE_KEY, size)
    applyFontSize(size)
  }

  // 监听变化
  watch(fontSize, (newSize) => {
    applyFontSize(newSize)
  })

  // 初始化时应用
  onMounted(() => {
    applyFontSize(fontSize.value)
  })

  return {
    fontSize,
    setFontSize,
    fontSizeConfig
  }
}
