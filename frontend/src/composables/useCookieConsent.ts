import { ref } from 'vue'

export interface CookiePreferences {
  necessary: boolean // Always true
  functional: boolean
  analytics: boolean
  marketing: boolean
  timestamp: string | null
  version: number
}

const STORAGE_KEY = 'teks_cookie_preferences_v1'

const defaultPreferences: CookiePreferences = {
  necessary: true,
  functional: true,
  analytics: true,
  marketing: false,
  timestamp: null,
  version: 1,
}

// Reactive singleton state
const hasConsented = ref<boolean>(false)
const isBannerVisible = ref<boolean>(false)
const isModalVisible = ref<boolean>(false)
const preferences = ref<CookiePreferences>({ ...defaultPreferences })

// Initialize from localStorage if present
let isInitialized = false

function initCookieConsent() {
  if (isInitialized || typeof window === 'undefined') return
  isInitialized = true

  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored) {
      const parsed = JSON.parse(stored) as CookiePreferences
      preferences.value = {
        ...defaultPreferences,
        ...parsed,
        necessary: true, // Strictly necessary must always stay true
      }
      hasConsented.value = true
      isBannerVisible.value = false
    } else {
      // First visit: sprout the cookie banner after a short natural delay
      hasConsented.value = false
      setTimeout(() => {
        isBannerVisible.value = true
      }, 700)
    }
  } catch (e) {
    console.error('Failed to parse cookie preferences from storage:', e)
    hasConsented.value = false
    isBannerVisible.value = true
  }
}

function persistPreferences(prefs: CookiePreferences) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(prefs))
  } catch (e) {
    console.warn('Unable to persist cookie preferences to localStorage:', e)
  }

  // Dispatch custom event for telemetry or third-party script integrations
  if (typeof window !== 'undefined') {
    window.dispatchEvent(
      new CustomEvent('cookie-consent-updated', {
        detail: { ...prefs },
      })
    )
  }
}

export function useCookieConsent() {
  initCookieConsent()

  function acceptAll() {
    const updated: CookiePreferences = {
      necessary: true,
      functional: true,
      analytics: true,
      marketing: true,
      timestamp: new Date().toISOString(),
      version: 1,
    }
    preferences.value = updated
    hasConsented.value = true
    isBannerVisible.value = false
    isModalVisible.value = false
    persistPreferences(updated)
  }

  function acceptEssentialOnly() {
    const updated: CookiePreferences = {
      necessary: true,
      functional: false,
      analytics: false,
      marketing: false,
      timestamp: new Date().toISOString(),
      version: 1,
    }
    preferences.value = updated
    hasConsented.value = true
    isBannerVisible.value = false
    isModalVisible.value = false
    persistPreferences(updated)
  }

  function saveCustomPreferences(custom: Partial<CookiePreferences>) {
    const updated: CookiePreferences = {
      ...preferences.value,
      ...custom,
      necessary: true, // Always enforce essential
      timestamp: new Date().toISOString(),
      version: 1,
    }
    preferences.value = updated
    hasConsented.value = true
    isBannerVisible.value = false
    isModalVisible.value = false
    persistPreferences(updated)
  }

  function openPreferences() {
    isModalVisible.value = true
  }

  function closePreferences() {
    isModalVisible.value = false
  }

  function dismissBanner() {
    isBannerVisible.value = false
  }

  function resetConsent() {
    try {
      localStorage.removeItem(STORAGE_KEY)
    } catch {}
    preferences.value = { ...defaultPreferences }
    hasConsented.value = false
    isBannerVisible.value = true
    isModalVisible.value = false
  }

  return {
    hasConsented,
    isBannerVisible,
    isModalVisible,
    preferences,
    acceptAll,
    acceptEssentialOnly,
    saveCustomPreferences,
    openPreferences,
    closePreferences,
    dismissBanner,
    resetConsent,
  }
}
