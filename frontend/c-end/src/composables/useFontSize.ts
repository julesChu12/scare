import { ref } from 'vue'

export type FontSize = 'small' | 'medium' | 'large'

const FONT_SIZE_KEY = 'c_font_size'
const DEFAULT_FONT_SIZE: FontSize = 'medium'

// 字体大小配置
export const fontSizeConfig = {
  small: { label: '小', value: '14px' },
  medium: { label: '中', value: '18px' },
  large: { label: '大', value: '22px' }
}

const isValidFontSize = (size: string | null): size is FontSize => {
  return size === 'small' || size === 'medium' || size === 'large'
}

const readSavedFontSize = (): FontSize => {
  if (typeof window === 'undefined') {
    return DEFAULT_FONT_SIZE
  }

  const savedSize = window.localStorage.getItem(FONT_SIZE_KEY)
  return isValidFontSize(savedSize) ? savedSize : DEFAULT_FONT_SIZE
}

const fontSize = ref<FontSize>(readSavedFontSize())

const applyFontSize = (size: FontSize) => {
  if (typeof document === 'undefined') {
    return
  }

  document.documentElement.setAttribute('data-font-size', size)
}

export const initFontSize = () => {
  fontSize.value = readSavedFontSize()
  applyFontSize(fontSize.value)
}

const setFontSize = (size: FontSize) => {
  fontSize.value = size

  if (typeof window !== 'undefined') {
    localStorage.setItem(FONT_SIZE_KEY, size)
  }

  applyFontSize(size)
}

export function useFontSize() {
  return {
    fontSize,
    setFontSize,
    fontSizeConfig
  }
}
