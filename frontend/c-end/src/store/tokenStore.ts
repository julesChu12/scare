import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useTokenStore = defineStore('token', () => {
  const TOKEN_KEY = 'c_token'
  const REFRESH_TOKEN_KEY = 'c_refresh_token'

  // State
  const token = ref<string | null>(localStorage.getItem(TOKEN_KEY))
  const refreshToken = ref<string | null>(localStorage.getItem(REFRESH_TOKEN_KEY))

  // Getters
  const isLoggedIn = computed(() => !!token.value)

  // Actions
  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem(TOKEN_KEY, newToken)
  }

  function setRefreshToken(newRefreshToken: string) {
    refreshToken.value = newRefreshToken
    localStorage.setItem(REFRESH_TOKEN_KEY, newRefreshToken)
  }

  function clearToken() {
    token.value = null
    refreshToken.value = null
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(REFRESH_TOKEN_KEY)
  }

  return {
    token,
    refreshToken,
    isLoggedIn,
    setToken,
    setRefreshToken,
    clearToken
  }
})
