import { defineComponent,  ref, reactive, onMounted } from 'vue'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'
import FlashAlert from '@/components/common/FlashAlert.vue'
import dayjs from 'dayjs'
import { useConfirm } from '@/composables/useConfirm'

export default defineComponent({
  name: 'AdminUsersPage',
  components: {
    FlashAlert
  },
  setup() {
    const { showFlash } = useFlash()

    const users = ref([])

    const meta = ref({ page: 1, limit: 10, total_count: 0 })

    const loading = ref(true)

    const saving = ref(false)

    const searchQuery = ref('')

    const showCreateModal = ref(false)

    const showEditModal = ref(false)

    const showCreditModal = ref(false)

    const formCreate = reactive({ name: '', email: '', password: '', role: 'User', is_activated: true })

    const formEdit = reactive({ id: '', name: '', email: '', role: 'User', is_activated: true })

    const formCredit = reactive({ user_id: '', targetName: '', targetEmail: '', amount: 10, reason: '' })

    onMounted(() => {
      fetchUsers(1)
    })

    async function fetchUsers(page = 1) {
      loading.value = true
      try {
        const res = await api.get(`/admin/users?page=${page}&limit=${meta.value.limit}&search=${encodeURIComponent(searchQuery.value)}`)
        users.value = res.data.users || []
        meta.value = res.data.meta
      } catch (err: any) {
        showFlash('Failed to load users', 'error')
      } finally {
        loading.value = false
      }
    }

    function formatDate(date) {
      return dayjs(date).format('MMM DD, YYYY')
    }

    function openCreateModal() {
      formCreate.name = ''
      formCreate.email = ''
      formCreate.password = ''
      formCreate.role = 'User'
      formCreate.is_activated = true
      showCreateModal.value = true
    }

    async function submitCreateUser() {
      saving.value = true
      try {
        await api.post('/admin/users', formCreate)
        showFlash('User created successfully', 'success')
        showCreateModal.value = false
        fetchUsers(1)
      } catch (err: any) {
        showFlash(err.response?.data?.error || 'Failed to create user', 'error')
      } finally {
        saving.value = false
      }
    }

    function openEditModal(user) {
      formEdit.id = user.id
      formEdit.name = user.name
      formEdit.email = user.email
      formEdit.role = user.role.name
      formEdit.is_activated = user.is_activated
      showEditModal.value = true
    }

    async function submitEditUser() {
      saving.value = true
      try {
        await api.put('/admin/users', formEdit)
        showFlash('User updated successfully', 'success')
        showEditModal.value = false
        fetchUsers(meta.value.page)
      } catch (err: any) {
        showFlash(err.response?.data?.error || 'Failed to update user', 'error')
      } finally {
        saving.value = false
      }
    }

    async function updateUserRole(user) {
      try {
        await api.put('/admin/users', {
          id: user.id,
          name: user.name,
          email: user.email,
          role: user.role.name,
          is_activated: user.is_activated
        })
        showFlash('User role updated', 'success')
      } catch (err: any) {
        showFlash('Failed to update role', 'error')
        fetchUsers(meta.value.page) // revert
      }
    }

    function onRoleSelectChange(user, newRole) {
      if (!user.role) user.role = {}
      user.role.name = newRole
      updateUserRole(user)
    }

    async function toggleUserStatus(user) {
      const newStatus = !user.is_activated
      try {
        await api.put('/admin/users', {
          id: user.id,
          name: user.name,
          email: user.email,
          role: user.role.name,
          is_activated: newStatus
        })
        user.is_activated = newStatus
        showFlash(`User ${newStatus ? 'activated' : 'suspended'} successfully`, 'success')
      } catch (err: any) {
        showFlash('Failed to update status', 'error')
      }
    }

    function openCreditModal(user) {
      formCredit.user_id = user.id
      formCredit.targetName = user.name
      formCredit.targetEmail = user.email
      formCredit.amount = 10
      formCredit.reason = 'Manual Adjustment'
      showCreditModal.value = true
    }

    async function submitCreditAllocation() {
      saving.value = true
      try {
        await api.post('/admin/users/credits', {
          user_id: formCredit.user_id,
          amount: formCredit.amount,
          reason: formCredit.reason
        })
        showFlash('Credits allocated successfully', 'success')
        showCreditModal.value = false
      } catch (err: any) {
        showFlash(err.response?.data?.error || 'Failed to allocate credits', 'error')
      } finally {
        saving.value = false
      }
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
        fetchUsers(meta.value.page)
      } catch (err: any) {
        showFlash('Failed to delete user', 'error')
      }
    }
    return {
      FlashAlert,
      showFlash,
      users,
      meta,
      loading,
      saving,
      searchQuery,
      showCreateModal,
      showEditModal,
      showCreditModal,
      formCreate,
      formEdit,
      formCredit,
      fetchUsers,
      formatDate,
      openCreateModal,
      submitCreateUser,
      openEditModal,
      submitEditUser,
      updateUserRole,
      onRoleSelectChange,
      toggleUserStatus,
      openCreditModal,
      submitCreditAllocation,
      askConfirm,
      deleteUser
    }
  }
})
