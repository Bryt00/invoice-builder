import { defineComponent,  ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import api from '@/utils/api'

export default defineComponent({
  name: 'InvoicesPage',
  setup() {
    const router = useRouter()

    const route = useRoute()

    const invoices = ref([])

    const loading = ref(true)

    const searchQuery = ref(String(route.query.q || ''))

    const statusFilter = ref(String(route.query.status || 'all'))

    let searchTimeout = null

    function formatDate(dateStr) {
      if (!dateStr) return '-'
      return new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: '2-digit', year: 'numeric' })
    }

    async function fetchInvoices() {
        loading.value = true
        try {
            const params = new URLSearchParams()
            if (searchQuery.value) params.append('q', searchQuery.value)
            if (statusFilter.value && statusFilter.value !== 'all') params.append('status', statusFilter.value)
            
            // Update URL to match filters without full reload
            router.replace({ query: { q: searchQuery.value || undefined, status: statusFilter.value !== 'all' ? statusFilter.value : undefined } })

            const res = await api.get(`/invoices?${params.toString()}`)
            if (res.data?.invoices) {
                invoices.value = res.data.invoices
            } else {
                invoices.value = []
            }
        } catch (err: any) {
            console.error('Failed to fetch invoices:', err)
            invoices.value = []
        } finally {
            loading.value = false
        }
    }

    function debouncedSearch() {
        clearTimeout(searchTimeout)
        searchTimeout = setTimeout(() => {
            fetchInvoices()
        }, 300)
    }

    function setStatusFilter(status) {
        statusFilter.value = status
        fetchInvoices()
    }

    onMounted(() => {
        fetchInvoices()
    })
    return {
      router,
      route,
      invoices,
      loading,
      searchQuery,
      statusFilter,
      searchTimeout,
      formatDate,
      fetchInvoices,
      debouncedSearch,
      setStatusFilter
    }
  }
})
