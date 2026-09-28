import { defineComponent,  ref, onMounted } from 'vue'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'
import FlashAlert from '@/components/common/FlashAlert.vue'
import dayjs from 'dayjs'

export default defineComponent({
  name: 'AdminWebhooksPage',
  components: {
    FlashAlert
  },
  setup() {
    const { showFlash } = useFlash()

    const webhooks = ref([])

    const meta = ref({ page: 1, limit: 15, total_count: 0 })

    const loading = ref(true)

    const statusFilter = ref('')

    onMounted(() => {
      fetchWebhooks(1)
    })

    async function fetchWebhooks(page = 1) {
      loading.value = true
      try {
        const res = await api.get(`/admin/webhooks?page=${page}&limit=${meta.value.limit}&status=${encodeURIComponent(statusFilter.value)}`)
        webhooks.value = res.data.webhooks || []
        meta.value = res.data.meta
      } catch (err: any) {
        showFlash('Failed to load webhook logs', 'error')
      } finally {
        loading.value = false
      }
    }

    async function replayWebhook(id) {
      try {
        await api.post(`/admin/webhooks/${id}/replay`)
        showFlash('Webhook queued for replay', 'success')
        fetchWebhooks(meta.value.page)
      } catch (err: any) {
        showFlash('Failed to replay webhook', 'error')
      }
    }

    function formatDate(date) {
      return dayjs(date).format('MMM DD, YYYY HH:mm:ss')
    }

    function getStatusBadge(status) {
      const base = 'px-2.5 py-1 rounded-full text-[10px] uppercase tracking-wider font-bold border '
      switch(status.toLowerCase()) {
        case 'processed':
          return base + 'bg-emerald-500/10 text-emerald-600 border-emerald-500/20'
        case 'pending':
          return base + 'bg-amber-500/10 text-amber-600 border-amber-500/20'
        case 'failed':
          return base + 'bg-rose-500/10 text-rose-600 border-rose-500/20'
        default:
          return base + 'bg-surface-container-high text-on-surface-variant border-outline-variant/40'
      }
    }
    return {
      FlashAlert,
      showFlash,
      webhooks,
      meta,
      loading,
      statusFilter,
      fetchWebhooks,
      replayWebhook,
      formatDate,
      getStatusBadge
    }
  }
})
