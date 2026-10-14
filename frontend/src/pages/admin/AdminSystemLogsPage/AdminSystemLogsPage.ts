import { defineComponent, ref, reactive, onMounted, onUnmounted } from 'vue'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'
import { useToast } from '@/composables/useToast'
import FlashAlert from '@/components/common/FlashAlert.vue'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'

dayjs.extend(relativeTime)

export interface SystemLog {
  id: number
  timestamp: string
  level: 'INFO' | 'WARN' | 'ERROR' | 'FATAL'
  source: string
  message: string
  method?: string
  path?: string
  status_code?: number
  latency_ms?: number
  ip_address?: string
  user_id?: string
  details?: string
}

export interface LogStats {
  total_logs: number
  error_count_24h: number
  warn_count_24h: number
  info_count_24h: number
  avg_latency_ms: number
  db_size_bytes: number
}

export default defineComponent({
  name: 'AdminSystemLogsPage',
  components: {
    FlashAlert
  },
  setup() {
    const { showFlash } = useFlash()
    const { showToast } = useToast()

    const logs = ref<SystemLog[]>([])
    const stats = ref<LogStats>({
      total_logs: 0,
      error_count_24h: 0,
      warn_count_24h: 0,
      info_count_24h: 0,
      avg_latency_ms: 0,
      db_size_bytes: 0
    })

    const meta = ref({ page: 1, limit: 50, total_count: 0 })
    const loading = ref(true)
    const isClearing = ref(false)
    const showClearConfirm = ref(false)
    const selectedLog = ref<SystemLog | null>(null)
    const expandedLogId = ref<number | null>(null)

    // Filters
    const selectedLevel = ref('ALL')
    const selectedSource = ref('ALL')
    const selectedStatusType = ref('all')
    const searchQuery = ref('')
    const isLiveTail = ref(true)

    let liveInterval: any = null

    onMounted(() => {
      fetchLogs(1)
      fetchStats()
      startLiveTail()
    })

    onUnmounted(() => {
      stopLiveTail()
    })

    function toggleLiveTail() {
      isLiveTail.value = !isLiveTail.value
      if (isLiveTail.value) {
        startLiveTail()
        showToast('Live log tail resumed (polling every 3s)', 'info')
      } else {
        stopLiveTail()
        showToast('Live log tail paused', 'info')
      }
    }

    function startLiveTail() {
      stopLiveTail()
      liveInterval = setInterval(() => {
        if (!loading.value && meta.value.page === 1) {
          fetchLogs(1, true)
          fetchStats()
        }
      }, 3000)
    }

    function stopLiveTail() {
      if (liveInterval) {
        clearInterval(liveInterval)
        liveInterval = null
      }
    }

    async function fetchLogs(page = 1, silent = false) {
      if (!silent) loading.value = true
      try {
        const params = new URLSearchParams({
          page: String(page),
          limit: String(meta.value.limit)
        })

        if (selectedLevel.value !== 'ALL') params.set('level', selectedLevel.value)
        if (selectedSource.value !== 'ALL') params.set('source', selectedSource.value)
        if (selectedStatusType.value !== 'all') params.set('status_type', selectedStatusType.value)
        if (searchQuery.value.trim()) params.set('search', searchQuery.value.trim())

        const res = await api.get(`/admin/system-logs?${params.toString()}`)
        logs.value = res.data.logs || []
        meta.value = res.data.meta || { page, limit: 50, total_count: 0 }
      } catch (err: any) {
        if (!silent) {
          showFlash('Failed to load system logs', 'error')
        }
      } finally {
        if (!silent) loading.value = false
      }
    }

    async function fetchStats() {
      try {
        const res = await api.get('/admin/system-logs/stats')
        if (res.data?.stats) {
          stats.value = res.data.stats
        }
      } catch (e) {
        console.warn('Failed to fetch log stats', e)
      }
    }

    async function clearAllLogs() {
      isClearing.value = true
      try {
        await api.delete('/admin/system-logs')
        showToast('System logs cleared successfully', 'success')
        showFlash('All system logs have been purged.', 'success')
        showClearConfirm.value = false
        logs.value = []
        meta.value.total_count = 0
        fetchStats()
      } catch (err: any) {
        showToast(err.response?.data?.error || 'Failed to clear system logs', 'error')
      } finally {
        isClearing.value = false
      }
    }

    function toggleExpandRow(id: number) {
      expandedLogId.value = expandedLogId.value === id ? null : id
    }

    function viewDetailsModal(logItem: SystemLog) {
      selectedLog.value = logItem
    }

    function copyToClipboard(text: string) {
      if (!text) return
      navigator.clipboard.writeText(text).then(() => {
        showToast('Copied details to clipboard', 'success')
      }).catch(() => {
        showToast('Failed to copy', 'error')
      })
    }

    function formatTime(timestamp: string) {
      if (!timestamp) return '-'
      return dayjs(timestamp).format('HH:mm:ss.SSS')
    }

    function formatDate(timestamp: string) {
      if (!timestamp) return '-'
      return dayjs(timestamp).format('MMM D, YYYY HH:mm:ss')
    }

    function formatRelative(timestamp: string) {
      if (!timestamp) return ''
      return dayjs(timestamp).fromNow()
    }

    function formatFileSize(bytes: number) {
      if (!bytes || bytes === 0) return '0 KB'
      if (bytes < 1024) return `${bytes} B`
      if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
      return `${(bytes / (1024 * 1024)).toFixed(2)} MB`
    }

    function getLevelBadgeClass(level: string) {
      switch (level) {
        case 'ERROR':
        case 'FATAL':
          return 'bg-error/15 text-error border-error/30 ring-1 ring-error/20 font-bold'
        case 'WARN':
          return 'bg-amber-500/15 text-amber-600 dark:text-amber-400 border-amber-500/30 font-semibold'
        case 'INFO':
        default:
          return 'bg-primary/15 text-primary border-primary/30 font-semibold'
      }
    }

    function getStatusBadgeClass(code?: number) {
      if (!code) return 'text-on-surface-variant'
      if (code >= 500) return 'bg-error/20 text-error font-bold border border-error/40'
      if (code >= 400) return 'bg-amber-500/20 text-amber-600 dark:text-amber-400 font-bold border border-amber-500/40'
      if (code >= 300) return 'bg-sky-500/20 text-sky-600 dark:text-sky-400 font-semibold border border-sky-500/40'
      if (code >= 200) return 'bg-emerald-500/20 text-emerald-600 dark:text-emerald-400 font-bold border border-emerald-500/40'
      return 'text-on-surface-variant'
    }

    function getMethodBadgeClass(method?: string) {
      switch (method?.toUpperCase()) {
        case 'GET':
          return 'text-emerald-600 dark:text-emerald-400 font-mono font-bold'
        case 'POST':
          return 'text-blue-600 dark:text-blue-400 font-mono font-bold'
        case 'PUT':
        case 'PATCH':
          return 'text-amber-600 dark:text-amber-400 font-mono font-bold'
        case 'DELETE':
          return 'text-error font-mono font-bold'
        default:
          return 'text-on-surface-variant font-mono font-bold'
      }
    }

    return {
      logs,
      stats,
      meta,
      loading,
      isClearing,
      showClearConfirm,
      selectedLog,
      expandedLogId,
      selectedLevel,
      selectedSource,
      selectedStatusType,
      searchQuery,
      isLiveTail,
      toggleLiveTail,
      fetchLogs,
      fetchStats,
      clearAllLogs,
      toggleExpandRow,
      viewDetailsModal,
      copyToClipboard,
      formatTime,
      formatDate,
      formatRelative,
      formatFileSize,
      getLevelBadgeClass,
      getStatusBadgeClass,
      getMethodBadgeClass
    }
  }
})
