import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import './style.css'
import { useToast } from './composables/useToast'

const app = createApp(App)

// Global Vue Error Boundary
const { showToast } = useToast()

app.config.errorHandler = (err, _instance, info) => {
  console.error('[Vue Error Boundary]:', err, info)
  const message = err instanceof Error ? err.message : String(err || 'Unknown error occurred')
  showToast(`Application error: ${message}`, 'error', 6000)
}

// Global Unhandled Promise Rejection Handler
window.addEventListener('unhandledrejection', (event) => {
  console.error('[Unhandled Promise Rejection]:', event.reason)
  const rawMsg = event.reason instanceof Error ? event.reason.message : String(event.reason || '')
  if (rawMsg.includes('canceled') || rawMsg.includes('AbortError')) return
  showToast(`Unhandled error: ${rawMsg || 'Something went wrong'}`, 'error', 6000)
})

app.use(createPinia())
app.use(router)

app.mount('#app')

// PWA Service Worker Registration
if ('serviceWorker' in navigator && import.meta.env.PROD) {
  window.addEventListener('load', () => {
    navigator.serviceWorker
      .register('/sw.js')
      .then((reg) => {
        console.log('[SW] ServiceWorker registered with scope: ', reg.scope)
      })
      .catch((err) => {
        console.warn('[SW] ServiceWorker registration failed: ', err)
      })
  })
}

