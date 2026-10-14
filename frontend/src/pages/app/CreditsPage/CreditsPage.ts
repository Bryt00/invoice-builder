import { defineComponent, ref, computed, onMounted } from 'vue'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'

export default defineComponent({
  name: 'CreditsPage',
  setup() {
    const { showFlash } = useFlash()

    const stats = ref({ balance: 0, total_purchased: 0, total_used: 0 })
    const history = ref<any[]>([])
    const loading = ref(true)

    // Pagination (5 records per page)
    const pageSize = 5
    const currentPage = ref(1)

    const totalRecords = computed(() => history.value.length)
    const totalPages = computed(() => Math.max(1, Math.ceil(totalRecords.value / pageSize)))
    const startIndex = computed(() => (totalRecords.value === 0 ? 0 : (currentPage.value - 1) * pageSize + 1))
    const endIndex = computed(() => Math.min(currentPage.value * pageSize, totalRecords.value))

    const paginatedHistory = computed(() => {
      const start = (currentPage.value - 1) * pageSize
      return history.value.slice(start, start + pageSize)
    })

    const visiblePages = computed(() => {
      const total = totalPages.value
      const current = currentPage.value
      if (total <= 5) {
        return Array.from({ length: total }, (_, i) => i + 1)
      }
      let start = Math.max(1, current - 2)
      let end = Math.min(total, start + 4)
      if (end - start < 4) {
        start = Math.max(1, end - 4)
      }
      const pages: number[] = []
      for (let i = start; i <= end; i++) {
        pages.push(i)
      }
      return pages
    })

    function prevPage() {
      if (currentPage.value > 1) {
        currentPage.value--
      }
    }

    function nextPage() {
      if (currentPage.value < totalPages.value) {
        currentPage.value++
      }
    }

    function goToPage(page: number) {
      if (page >= 1 && page <= totalPages.value) {
        currentPage.value = page
      }
    }

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

    function formatDate(dateStr: any) {
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
      paginatedHistory,
      loading,
      currentPage,
      pageSize,
      totalRecords,
      totalPages,
      startIndex,
      endIndex,
      visiblePages,
      prevPage,
      nextPage,
      goToPage,
      formatDate,
    }
  },
})
