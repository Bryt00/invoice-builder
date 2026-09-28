import { defineComponent,  ref } from 'vue'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'
import FlashAlert from '@/components/common/FlashAlert.vue'

export default defineComponent({
  name: 'ForgotPasswordPage',
  components: {
    FlashAlert
  },
  setup() {
    const email = ref('')

    const loading = ref(false)

    const { showFlash } = useFlash()

    async function handleSubmit() {
      loading.value = true
      try {
        const res = await api.post('/auth/forgot-password', { email: email.value })
        showFlash(res.data?.message || 'Password reset link sent to your email.', 'success', 8000)
      } catch (err: any) {
        showFlash('Failed to send reset link', 'error')
      } finally {
        loading.value = false
      }
    }
    return {
      FlashAlert,
      email,
      loading,
      showFlash,
      handleSubmit
    }
  }
})
