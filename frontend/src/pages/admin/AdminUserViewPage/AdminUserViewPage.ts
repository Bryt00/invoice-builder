import { defineComponent,  ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'
import FlashAlert from '@/components/common/FlashAlert.vue'
import dayjs from 'dayjs'

export default defineComponent({
  name: 'AdminUserViewPage',
  components: {
    FlashAlert
  },
  setup() {
    const route = useRoute()

    const router = useRouter()

    const { showFlash } = useFlash()

    const user = ref(null)

    const loading = ref(true)

    onMounted(async () => {
      const id = route.params.id
      if (!id) {
        router.push('/user/admin/users')
        return
      }
      
      try {
        const res = await api.get(`/admin/users/${id}`)
        user.value = res.data.user
      } catch (err: any) {
        showFlash('Failed to load user profile', 'error')
      } finally {
        loading.value = false
      }
    })

    function formatDate(date) {
      return dayjs(date).format('MMMM DD, YYYY at h:mm A')
    }
    return {
      FlashAlert,
      route,
      router,
      showFlash,
      user,
      loading,
      formatDate
    }
  }
})
