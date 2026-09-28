import { defineComponent,  ref, onMounted } from 'vue'
import api from '@/utils/api'
import { useToast } from '@/composables/useToast'
import dayjs from 'dayjs'

export default defineComponent({
  name: 'AdminPaymentsPage',
  setup() {
    const { showToast } = useToast()

    const payments = ref([])

    const meta = ref({ page: 1, limit: 10, total_count: 0 })

    const loading = ref(true)

    const statusFilter = ref('')

    onMounted(() => {
      fetchPayments(1)
    })

    const currencyMap = {
      USD: '$', EUR: '€', GBP: '£', GHS: 'GH₵', NGN: '₦'
    }

    function currencySymbol(code) {
      return currencyMap[code] || code
    }

    async function fetchPayments(page = 1) {
      loading.value = true
      try {
        const res = await api.get(`/admin/payments?page=${page}&limit=${meta.value.limit}&status=${encodeURIComponent(statusFilter.value)}`)
        payments.value = res.data.payments || []
        meta.value = res.data.meta || { page: 1, limit: 10, total_count: payments.value.length }
      } catch (err: any) {
        showToast('Failed to load payments', 'error')
      } finally {
        loading.value = false
      }
    }

    function formatDate(date) {
      if (!date) return '-'
      return dayjs(date).format('MMM DD, YYYY HH:mm')
    }

    function getStatusBadge(status) {
      const base = 'px-2.5 py-1 rounded-full text-xs font-semibold border '
      if (!status) return base + 'bg-surface-container-high text-on-surface-variant border-outline-variant/40'
      switch(String(status).toLowerCase()) {
        case 'successful':
        case 'succeeded':
        case 'completed':
        case 'paid':
          return base + 'bg-emerald-500/10 text-emerald-600 border-emerald-500/20'
        case 'pending':
          return base + 'bg-amber-500/10 text-amber-600 border-amber-500/20'
        case 'failed':
          return base + 'bg-rose-500/10 text-rose-600 border-rose-500/20'
        case 'refunded':
          return base + 'bg-purple-500/10 text-purple-600 border-purple-500/20'
        default:
          return base + 'bg-surface-container-high text-on-surface-variant border-outline-variant/40'
      }
    }
    return {
      showToast,
      payments,
      meta,
      loading,
      statusFilter,
      currencyMap,
      currencySymbol,
      fetchPayments,
      formatDate,
      getStatusBadge
    }
  }
})
