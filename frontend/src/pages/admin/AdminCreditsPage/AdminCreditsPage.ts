import { defineComponent,  ref, onMounted } from 'vue'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'
import FlashAlert from '@/components/common/FlashAlert.vue'
import dayjs from 'dayjs'

export default defineComponent({
  name: 'AdminCreditsPage',
  components: {
    FlashAlert
  },
  setup() {
    const { showFlash } = useFlash()

    const credits = ref([])

    const meta = ref({ page: 1, limit: 10, total_count: 0 })

    const loading = ref(true)

    const typeFilter = ref('')

    onMounted(() => {
      fetchCredits(1)
    })

    async function fetchCredits(page = 1) {
      loading.value = true
      try {
        const res = await api.get(`/admin/credits?page=${page}&limit=${meta.value.limit}&type=${encodeURIComponent(typeFilter.value)}`)
        credits.value = res.data.credits || []
        meta.value = res.data.meta
      } catch (err: any) {
        showFlash('Failed to load credit history', 'error')
      } finally {
        loading.value = false
      }
    }

    function formatDate(date) {
      return dayjs(date).format('MMM DD, YYYY HH:mm')
    }

    function getTypeBadge(type) {
      const base = 'px-2.5 py-1 rounded-full text-xs font-semibold border '
      if (!type) return base + 'bg-surface-container-high text-on-surface-variant border-outline-variant/40'
      switch(type.toLowerCase()) {
        case 'purchase':
        case 'grant':
          return base + 'bg-emerald-500/10 text-emerald-600 border-emerald-500/20'
        case 'deduction':
          return base + 'bg-amber-500/10 text-amber-600 border-amber-500/20'
        default:
          return base + 'bg-surface-container-high text-on-surface-variant border-outline-variant/40'
      }
    }
    return {
      FlashAlert,
      showFlash,
      credits,
      meta,
      loading,
      typeFilter,
      fetchCredits,
      formatDate,
      getTypeBadge
    }
  }
})
