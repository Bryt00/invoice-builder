import { defineComponent,  ref, reactive, onMounted, computed } from 'vue'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'
import { useAuthStore } from '@/stores/auth'

export default defineComponent({
  name: 'FinancePage',
  setup() {
    const authStore = useAuthStore()

    const transactions = ref([])

    const summary = ref({ total_income: 0, total_expenses: 0, net_profit: 0 })

    const categories = ref([])

    const loading = ref(true)

    const { showFlash } = useFlash()

    const currencySymbol = computed(() => authStore.currencySymbol)

    const filters = reactive({
        category_id: 'all',
        start_date: '',
        end_date: '',
        search: ''
    })

    const filteredTransactions = computed(() => {
        let result = transactions.value
        
        if (filters.category_id !== 'all') {
            result = result.filter(t => t.category_id === filters.category_id)
        }
        if (filters.start_date) {
            result = result.filter(t => new Date(t.transaction_date) >= new Date(filters.start_date))
        }
        if (filters.end_date) {
            result = result.filter(t => new Date(t.transaction_date) <= new Date(filters.end_date))
        }
        if (filters.search) {
            const s = filters.search.toLowerCase()
            result = result.filter(t => 
                (t.title && t.title.toLowerCase().includes(s)) || 
                (t.description && t.description.toLowerCase().includes(s)) ||
                (t.payee_or_payer && t.payee_or_payer.toLowerCase().includes(s))
            )
        }
        return result
    })

    async function fetchData() {
        loading.value = true
        try {
            const [txnRes, sumRes, catRes] = await Promise.all([
                api.get('/finance/transactions'),
                api.get('/finance/summary'),
                api.get('/finance/categories'),
            ])

            if (txnRes.data?.transactions) transactions.value = txnRes.data.transactions
            if (sumRes.data?.summary) summary.value = sumRes.data.summary
            if (catRes.data?.categories) {
                categories.value = catRes.data.categories
            }
        } catch (err: any) {
            // Handle error silently
        } finally {
            loading.value = false
        }
    }

    onMounted(fetchData)

    function formatDate(dateStr) {
        if (!dateStr) return '-'
        const d = new Date(dateStr)
        if (isNaN(d.getTime())) return dateStr
        return d.toLocaleDateString('en-US', { month: 'short', day: '2-digit', year: 'numeric' })
    }

    async function exportCsv() {
        try {
            const res = await api.get('/finance/export', { responseType: 'blob' })
            const url = window.URL.createObjectURL(new Blob([res.data]))
            const link = document.createElement('a')
            link.href = url
            link.setAttribute('download', 'finance_export.csv')
            document.body.appendChild(link)
            link.click()
            link.parentNode.removeChild(link)
            showFlash('Export downloaded successfully!', 'success')
        } catch (err: any) {
            showFlash('Failed to export CSV', 'error')
        }
    }
    return {
      authStore,
      transactions,
      summary,
      categories,
      loading,
      showFlash,
      currencySymbol,
      filters,
      filteredTransactions,
      fetchData,
      formatDate,
      exportCsv
    }
  }
})
