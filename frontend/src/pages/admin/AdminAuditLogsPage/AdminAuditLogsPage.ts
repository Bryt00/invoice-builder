import { defineComponent,  ref, onMounted } from 'vue'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'
import FlashAlert from '@/components/common/FlashAlert.vue'
import dayjs from 'dayjs'

export default defineComponent({
  name: 'AdminAuditLogsPage',
  components: {
    FlashAlert
  },
  setup() {
    const { showFlash } = useFlash()

    const logs = ref([])

    const meta = ref({ page: 1, limit: 15, total_count: 0 })

    const loading = ref(true)

    const searchQuery = ref('')

    onMounted(() => {
      fetchLogs(1)
    })

    async function fetchLogs(page = 1) {
      loading.value = true
      try {
        const res = await api.get(`/admin/audit-logs?page=${page}&limit=${meta.value.limit}&search=${encodeURIComponent(searchQuery.value)}`)
        logs.value = res.data.logs || []
        meta.value = res.data.meta
      } catch (err: any) {
        showFlash('Failed to load audit logs', 'error')
      } finally {
        loading.value = false
      }
    }

    function formatDate(date) {
      return dayjs(date).format('MMM DD, YYYY HH:mm:ss')
    }
    return {
      FlashAlert,
      showFlash,
      logs,
      meta,
      loading,
      searchQuery,
      fetchLogs,
      formatDate
    }
  }
})
