import { defineComponent,  ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'
import FlashAlert from '@/components/common/FlashAlert.vue'

export default defineComponent({
  name: 'ResetPasswordPage',
  components: {
    FlashAlert
  },
  setup() {
    const route = useRoute()

    const router = useRouter()

    const newPassword = ref('')

    const loading = ref(false)

    const showPassword = ref(false)

    const { showFlash } = useFlash()

    async function handleSubmit() {
      const token = route.query.token
      if (!token) {
        showFlash('Missing password reset token', 'error')
        return
      }

      loading.value = true
      try {
        const res = await api.post('/auth/reset-password', {
          token,
          new_password: newPassword.value,
        })
        showFlash(res.data?.message || 'Password reset successfully! Please sign in.', 'success', 8000)
        router.push('/user/login')
      } catch (err: any) {
        showFlash(err.response?.data?.error || 'Failed to reset password', 'error')
      } finally {
        loading.value = false
      }
    }
    return {
      FlashAlert,
      route,
      router,
      newPassword,
      loading,
      showPassword,
      showFlash,
      handleSubmit
    }
  }
})
