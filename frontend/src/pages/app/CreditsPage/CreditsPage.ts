import { defineComponent,  ref, onMounted } from 'vue'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'

export default defineComponent({
  name: 'CreditsPage',
  setup() {
    const { showFlash } = useFlash()

    const stats = ref({ balance: 0, total_purchased: 0, total_used: 0 })

    const history = ref([])

    const loading = ref(true)

    onMounted(async () => {
        loading.value = true
        try {
            const [balRes, histRes] = await Promise.all([
                api.get('/credits/balance'),
                api.get('/credits/history'),
            ])

            if (balRes.data?.stats) {
                stats.value = balRes.data.stats
            } else if (balRes.data?.balance !== undefined) {
                // Fallback if API returns just balance
                stats.value.balance = balRes.data.balance
            }
            
            if (histRes.data?.history) {
                history.value = histRes.data.history
            } else if (histRes.data?.transactions) {
                history.value = histRes.data.transactions
            }
        } catch (err: any) {
            showFlash('Failed to load credit history', 'error')
        } finally {
            loading.value = false
        }
    })

    function formatDate(dateStr) {
        if (!dateStr) return '-'
        const d = new Date(dateStr)
        if (isNaN(d.getTime())) return dateStr
        
        // Format like "Jan 02, 2006 • 15:04"
        const optsDate: Intl.DateTimeFormatOptions = { month: 'short', day: '2-digit', year: 'numeric' }
        const optsTime: Intl.DateTimeFormatOptions = { hour: '2-digit', minute: '2-digit', hour12: false }
        
        return `${d.toLocaleDateString('en-US', optsDate)} • ${d.toLocaleTimeString('en-US', optsTime)}`
    }
    return {
      showFlash,
      stats,
      history,
      loading,
      formatDate
    }
  }
})
