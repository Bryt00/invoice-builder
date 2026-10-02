import { defineComponent,  reactive, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { useFlash } from '@/composables/useFlash'
import FlashAlert from '@/components/common/FlashAlert.vue'

export default defineComponent({
  name: 'RegisterPage',
  components: {
    FlashAlert
  },
  setup() {
    const authStore = useAuthStore()

    const router = useRouter()

    const { showFlash } = useFlash()

    const form = reactive({
      name: '',
      email: '',
      password: '',
      confirmPassword: '',
    })

    const errors = ref<Record<string, string>>({})

    const showPassword = ref(false)

    const showConfirmPassword = ref(false)

    async function handleSubmit() {
      errors.value = {}
      if (form.password !== form.confirmPassword) {
        errors.value['confirm password'] = 'Passwords do not match'
        return
      }

      try {
        const res = await authStore.register(form.name, form.email, form.password, form.confirmPassword)
        showFlash(res.message || 'Account created! Please check your email for activation link.', 'success', 8000)
        router.push('/user/login')
      } catch (err: any) {
        if (err.response?.status === 422 && typeof err.response.data?.error === 'object') {
          errors.value = err.response.data.error
        } else {
          showFlash(authStore.error || 'Registration failed', 'error')
        }
      }
    }
    return {
      FlashAlert,
      authStore,
      router,
      showFlash,
      form,
      errors,
      showPassword,
      showConfirmPassword,
      handleSubmit
    }
  }
})
