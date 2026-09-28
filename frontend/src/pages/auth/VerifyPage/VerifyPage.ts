import { defineComponent,  ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/utils/api'

export default defineComponent({
  name: 'VerifyPage',
  setup() {
    const route = useRoute()

    const loading = ref(true)

    const success = ref(false)

    const errorMessage = ref('')

    onMounted(async () => {
      const token = route.query.token
      if (!token) {
        loading.value = false
        errorMessage.value = 'Invalid or missing activation token.'
        return
      }

      try {
        await api.get(`/auth/verify?token=${token}`)
        success.value = true
      } catch (err: any) {
        errorMessage.value = err.response?.data?.error || 'The activation link is invalid or has expired.'
      } finally {
        loading.value = false
      }
    })
    return {
      route,
      loading,
      success,
      errorMessage
    }
  }
})
