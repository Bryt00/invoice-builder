import { defineComponent,  ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'
import FlashAlert from '@/components/common/FlashAlert.vue'
import dayjs from 'dayjs'
import { useConfirm } from '@/composables/useConfirm'

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

    const { askConfirm } = useConfirm()

    async function deleteUser(id) {
      const ok = await askConfirm({
        title: 'Delete User Account',
        message: 'Are you sure you want to permanently delete this user? This action cannot be undone.',
        confirmText: 'Delete Account',
        type: 'danger'
      })
      if (!ok) return

      try {
        await api.delete(`/admin/users/${id}`)
        showFlash('User deleted successfully', 'success')
        router.push('/user/admin/users')
      } catch (err: any) {
        showFlash('Failed to delete user', 'error')
      }
    }

    async function forcePasswordReset(id) {
      const ok = await askConfirm({
        title: 'Force Password Reset',
        message: 'Are you sure you want to force a password reset for this user? They will receive an email with instructions.',
        confirmText: 'Send Reset Email',
        type: 'warning'
      })
      if (!ok) return

      try {
        await api.post(`/admin/users/password-reset`, { user_id: id })
        showFlash('Password reset email sent successfully', 'success')
      } catch (err: any) {
        showFlash('Failed to send password reset', 'error')
      }
    }

    async function impersonateUser(id) {
      const ok = await askConfirm({
        title: 'Impersonate User',
        message: 'Are you sure you want to log in as this user? You will be redirected to their dashboard.',
        confirmText: 'Impersonate',
        type: 'warning'
      })
      if (!ok) return

      try {
        const res = await api.post(`/admin/users/${id}/impersonate`)
        if (res.data?.tokens) {
          localStorage.setItem('access_token', res.data.tokens.access_token)
          localStorage.setItem('refresh_token', res.data.tokens.refresh_token)
          window.location.href = '/user/dashboard'
        }
      } catch (err: any) {
        showFlash('Failed to start impersonation session', 'error')
      }
    }

    return {
      FlashAlert,
      route,
      router,
      showFlash,
      user,
      loading,
      formatDate,
      deleteUser,
      forcePasswordReset,
      impersonateUser
    }
  }
})
