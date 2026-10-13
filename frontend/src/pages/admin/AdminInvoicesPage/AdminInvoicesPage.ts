import { defineComponent,  ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'
import FlashAlert from '@/components/common/FlashAlert.vue'
import dayjs from 'dayjs'
import { useConfirm } from '@/composables/useConfirm'

export default defineComponent({
  name: 'AdminInvoicesPage',
  components: {
    FlashAlert
  },
  setup() {
    const router = useRouter()
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
      return dayjs(date).format('MMMM DD, YYYY at h:mm A')
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

    const { askConfirm } = useConfirm()

    async function deleteInvoice(id) {
      const ok = await askConfirm({
        title: 'Delete Invoice',
        message: 'Are you sure you want to permanently delete this invoice? This cannot be undone.',
        confirmText: 'Delete Invoice',
        type: 'danger'
      })
      if (!ok) return

      try {
        await api.delete(`/admin/invoices/${id}`)
        showFlash('Invoice deleted successfully', 'success')
        fetchInvoices(meta.value.page)
      } catch (err: any) {
        showFlash('Failed to delete invoice', 'error')
      }
    }

    async function forceUpdateStatus(id, newStatus) {
      const ok = await askConfirm({
        title: 'Force Update Status',
        message: `Are you sure you want to forcefully mark this invoice as ${newStatus.toUpperCase()}?`,
        confirmText: 'Update Status',
        type: 'warning'
      })
      if (!ok) return

      try {
        await api.put(`/admin/invoices/${id}/status`, { status: newStatus })
        showFlash('Invoice status updated', 'success')
        fetchInvoices(meta.value.page)
      } catch (err: any) {
        showFlash('Failed to update invoice status', 'error')
      }
    }

    async function exportCSV() {
      try {
        const res = await api.get(`/admin/invoices/export?search=${encodeURIComponent(searchQuery.value)}&status=${encodeURIComponent(statusFilter.value)}`, { responseType: 'blob' })
        const url = window.URL.createObjectURL(new Blob([res.data]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', 'invoices_export.csv')
        document.body.appendChild(link)
        link.click()
        link.remove()
      } catch (e) {
        showFlash('Failed to export CSV', 'error')
      }
    }

    return {
      router,
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
      getStatusBadge,
      deleteInvoice,
      forceUpdateStatus,
      exportCSV
    }
  }
})
