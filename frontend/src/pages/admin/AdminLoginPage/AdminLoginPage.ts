import { defineComponent,  ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

export default defineComponent({
  name: 'AdminLoginPage',
  setup() {
    const email = ref('')

    const password = ref('')

    const showPassword = ref(false)

    const loading = ref(false)

    const error = ref('')

    const authStore = useAuthStore()

    const router = useRouter()

    async function handleAdminLogin() {
      loading.value = true
      error.value = ''

      try {
        await authStore.login(email.value, password.value)
        if (authStore.user?.role?.name === 'Admin') {
          router.push('/user/admin/dashboard')
        } else {
          error.value = 'Access Denied: Administrator credentials required.'
        }
      } catch (err: any) {
        error.value = err.response?.data?.error || 'Authentication failed'
      } finally {
        loading.value = false
      }
    }
    return {
      email,
      password,
      showPassword,
      loading,
      error,
      authStore,
      router,
      handleAdminLogin
    }
  }
})
