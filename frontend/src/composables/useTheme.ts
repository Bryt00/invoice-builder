import { ref, computed } from 'vue'

export type ThemeMode = 'light' | 'dark' | 'system'

const THEME_KEY = 'teks_invoice_theme'

// Read initial stored theme or fallback to system
const storedTheme = (localStorage.getItem(THEME_KEY) as ThemeMode) || 'light'
const currentTheme = ref<ThemeMode>(storedTheme)

const mediaQuery = typeof window !== 'undefined' ? window.matchMedia('(prefers-color-scheme: dark)') : null
const systemPrefersDark = ref<boolean>(mediaQuery ? mediaQuery.matches : false)

if (mediaQuery) {
  mediaQuery.addEventListener('change', (e) => {
    systemPrefersDark.value = e.matches
    if (currentTheme.value === 'system') {
      applyTheme()
    }
  })
}

function applyTheme() {
  if (typeof document === 'undefined') return
  const isDark = currentTheme.value === 'dark' || (currentTheme.value === 'system' && systemPrefersDark.value)
  if (isDark) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}

// Initial application
applyTheme()

export function useTheme() {
  const isDark = computed(() => {
    return currentTheme.value === 'dark' || (currentTheme.value === 'system' && systemPrefersDark.value)
  })

  function setTheme(theme: ThemeMode) {
    currentTheme.value = theme
    localStorage.setItem(THEME_KEY, theme)
    applyTheme()
  }

  function toggleTheme() {
    if (isDark.value) {
      setTheme('light')
    } else {
      setTheme('dark')
    }
  }

  return {
    theme: currentTheme,
    isDark,
    setTheme,
    toggleTheme,
  }
}
