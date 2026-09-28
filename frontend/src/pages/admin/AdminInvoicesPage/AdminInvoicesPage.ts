import { defineComponent,  ref, onMounted } from 'vue'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'
import FlashAlert from '@/components/common/FlashAlert.vue'
import dayjs from 'dayjs'

export default defineComponent({
  name: 'AdminInvoicesPage',
  components: {
    FlashAlert
  },
  setup() {
    const { showFlash } = useFlash()

    const invoices = ref([])

    const meta = ref({ page: 1, limit: 10, total_count: 0 })

    const loading = ref(true)

    const searchQuery = ref('')

    const statusFilter = ref('')

    onMounted(() => {
      fetchInvoices(1)
    })

    const currencyMap = {
      USD: '$', EUR: '€', GBP: '£', GHS: 'GH₵', NGN: '₦'
    }

    function currencySymbol(code) {
      return currencyMap[code] || code
    }

    async function fetchInvoices(page = 1) {
      loading.value = true
      try {
        const res = await api.get(`/admin/invoices?page=${page}&limit=${meta.value.limit}&search=${encodeURIComponent(searchQuery.value)}&status=${encodeURIComponent(statusFilter.value)}`)
        invoices.value = res.data.invoices || []
        meta.value = res.data.meta
      } catch (err: any) {
        showFlash('Failed to load invoices', 'error')
      } finally {
        loading.value = false
      }
    }

    function formatDate(date) {
      if (!date || date === '0001-01-01T00:00:00Z') return 'N/A'
      return dayjs(date).format('MMM DD, YYYY')
    }

    function getStatusBadge(status) {
      const base = 'px-2.5 py-1 rounded-full text-xs font-semibold border '
      switch(status.toLowerCase()) {
        case 'paid':
          return base + 'bg-emerald-500/10 text-emerald-600 border-emerald-500/20'
        case 'pending':
          return base + 'bg-amber-500/10 text-amber-600 border-amber-500/20'
        case 'overdue':
          return base + 'bg-rose-500/10 text-rose-600 border-rose-500/20'
        case 'draft':
        default:
          return base + 'bg-surface-container-high text-on-surface-variant border-outline-variant/40'
      }
    }
    return {
      FlashAlert,
      showFlash,
      invoices,
      meta,
      loading,
      searchQuery,
      statusFilter,
      currencyMap,
      currencySymbol,
      fetchInvoices,
      formatDate,
      getStatusBadge
    }
  }
})
