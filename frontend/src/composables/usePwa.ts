import { ref } from 'vue'
import { registerSW } from 'virtual:pwa-register'

export interface BeforeInstallPromptEvent extends Event {
  readonly platforms: string[]
  readonly userChoice: Promise<{
    outcome: 'accepted' | 'dismissed'
    platform: string
  }>
  prompt(): Promise<void>
}

// Global reactive states shared across all components
const needRefresh = ref(false)
const offlineReady = ref(false)
const isOffline = ref(typeof navigator !== 'undefined' ? !navigator.onLine : false)
const canInstall = ref(false)
const isInstalled = ref(false)

let deferredPrompt: BeforeInstallPromptEvent | null = null
let updateServiceWorkerFn: ((reloadPage?: boolean) => Promise<void>) | null = null
let isInitialized = false

export function usePwa() {
  const initPwa = () => {
    if (isInitialized || typeof window === 'undefined') return
    isInitialized = true

    // Check if already in standalone display mode
    if (
      window.matchMedia('(display-mode: standalone)').matches ||
      (window.navigator as unknown as { standalone?: boolean }).standalone
    ) {
      isInstalled.value = true
    }

    // Online / Offline network listeners
    window.addEventListener('online', () => {
      isOffline.value = false
    })
    window.addEventListener('offline', () => {
      isOffline.value = true
    })

    // Capture PWA installation prompt
    window.addEventListener('beforeinstallprompt', (e: Event) => {
      e.preventDefault()
      deferredPrompt = e as BeforeInstallPromptEvent
      canInstall.value = true
    })

    // Listen for successful installation
    window.addEventListener('appinstalled', () => {
      deferredPrompt = null
      canInstall.value = false
      isInstalled.value = true
      console.log('[PWA] Application successfully installed')
    })

    // Register Service Worker in production
    if ('serviceWorker' in navigator && import.meta.env.PROD) {
      updateServiceWorkerFn = registerSW({
        immediate: true,
        onNeedRefresh() {
          console.log('[PWA] New content available; prompt user to refresh.')
          needRefresh.value = true
        },
        onOfflineReady() {
          console.log('[PWA] App is ready to work offline.')
          offlineReady.value = true
        },
        onRegisterError(error) {
          console.warn('[PWA] Service Worker registration failed:', error)
        }
      })
    }
  }

  const promptInstall = async (): Promise<boolean> => {
    if (!deferredPrompt) return false
    try {
      await deferredPrompt.prompt()
      const { outcome } = await deferredPrompt.userChoice
      if (outcome === 'accepted') {
        canInstall.value = false
        deferredPrompt = null
        return true
      }
      return false
    } catch (err) {
      console.error('[PWA] Error during installation prompt:', err)
      return false
    }
  }

  const updateApp = async () => {
    if (updateServiceWorkerFn) {
      await updateServiceWorkerFn(true)
    } else {
      window.location.reload()
    }
  }

  const dismissRefresh = () => {
    needRefresh.value = false
  }

  return {
    needRefresh,
    offlineReady,
    isOffline,
    canInstall,
    isInstalled,
    initPwa,
    promptInstall,
    updateApp,
    dismissRefresh
  }
}
