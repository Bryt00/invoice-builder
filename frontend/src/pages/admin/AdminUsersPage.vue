<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6 animate-fade-in">
    <FlashAlert />

    <!-- Page Header & Search -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
      <div>
        <h1 class="font-headline text-2xl sm:text-3xl font-extrabold text-on-surface flex items-center gap-2.5">
          <span class="material-symbols-outlined text-amber-500 text-[28px]">group</span>
          User Directory
        </h1>
        <p class="font-body text-sm text-on-surface-variant mt-1">Manage user account statuses, roles, and credit balances.</p>
      </div>
      <div class="flex flex-wrap items-center gap-2 w-full sm:w-auto">
        <form @submit.prevent="fetchUsers(1)" class="flex items-center gap-2 w-full sm:w-auto">
          <div class="relative w-full sm:w-64">
            <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-[18px]">search</span>
            <input type="text" v-model="searchQuery" placeholder="Search by name or email..." class="w-full pl-9 pr-4 py-2 bg-surface-container-low border border-outline-variant/60 rounded-xl font-body text-sm text-on-surface focus:outline-none focus:border-amber-500">
          </div>
          <button type="submit" class="bg-amber-500 hover:bg-amber-600 text-on-primary px-4 py-2 rounded-xl font-label text-sm font-bold transition-colors cursor-pointer">Search</button>
        </form>
        <button type="button" @click="openCreateModal" class="inline-flex items-center gap-1.5 px-3 py-2 bg-amber-500 hover:bg-amber-400 text-on-primary rounded-xl font-label text-xs font-bold transition-colors cursor-pointer">
          <span class="material-symbols-outlined text-[16px]">person_add</span>
          <span>Create User</span>
        </button>
      </div>
    </div>

    <!-- Users Table -->
    <div class="glass-card rounded-2xl border border-outline-variant/60 overflow-hidden shadow-xs relative min-h-[300px]">
      
      <!-- Loading Overlay -->
      <div v-if="loading" class="absolute inset-0 bg-surface/50 backdrop-blur-sm z-10 flex items-center justify-center">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-amber-500"></div>
      </div>

      <div class="px-6 py-4 border-b border-outline-variant/40 flex justify-between items-center bg-surface-container-lowest/50">
        <h3 class="font-headline text-base sm:text-lg font-bold text-on-surface">
          Registered Users ({{ meta.total_count || 0 }})
        </h3>
      </div>
      
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="border-b border-outline-variant/40 font-label text-xs uppercase text-on-surface-variant/80 bg-surface-container-low/40">
              <th class="px-6 py-3.5">User</th>
              <th class="px-6 py-3.5">Email</th>
              <th class="px-6 py-3.5">Role</th>
              <th class="px-6 py-3.5">Status</th>
              <th class="px-6 py-3.5">Registered</th>
              <th class="px-6 py-3.5 text-right">Admin Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-outline-variant/30 font-body text-sm text-on-surface">
            <tr v-for="user in users" :key="user.id" class="hover:bg-surface-container-low/40 transition-colors">
              <td class="px-6 py-4 font-semibold text-on-surface flex items-center gap-2.5">
                <span class="w-8 h-8 rounded-full bg-amber-500/20 text-amber-600 border border-amber-500/30 flex items-center justify-center font-bold text-xs shrink-0">
                  {{ user.name.charAt(0).toUpperCase() }}
                </span>
                <span class="hover:text-amber-500 transition-colors">{{ user.name }}</span>
              </td>
              <td class="px-6 py-4 text-on-surface-variant">{{ user.email }}</td>
              <td class="px-6 py-4">
                <select v-model="user.role.name" @change="updateUserRole(user)" class="bg-surface-container-high border border-outline-variant/40 rounded-lg px-2 py-1 font-label text-xs font-bold cursor-pointer focus:outline-none">
                  <option value="User">USER</option>
                  <option value="Admin">ADMIN</option>
                </select>
              </td>
              <td class="px-6 py-4">
                <span v-if="user.is_activated" class="px-2.5 py-1 rounded-full text-xs font-semibold bg-emerald-500/10 text-emerald-600 border border-emerald-500/20 inline-flex items-center gap-1">
                  <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span> Active
                </span>
                <span v-else class="px-2.5 py-1 rounded-full text-xs font-semibold bg-amber-500/10 text-amber-600 border border-amber-500/20 inline-flex items-center gap-1">
                  <span class="w-1.5 h-1.5 rounded-full bg-amber-500"></span> Suspended
                </span>
              </td>
              <td class="px-6 py-4 text-on-surface-variant">{{ formatDate(user.created_at) }}</td>
              <td class="px-6 py-4 text-right flex items-center justify-end gap-1.5">
                <button type="button" @click="openEditModal(user)" class="p-1.5 text-on-surface-variant hover:text-primary hover:bg-surface-container-high rounded-lg transition-colors cursor-pointer" title="Edit User">
                  <span class="material-symbols-outlined text-[18px]">edit</span>
                </button>
                <button type="button" @click="openCreditModal(user)" class="p-1.5 text-on-surface-variant hover:text-amber-500 hover:bg-amber-500/10 rounded-lg transition-colors cursor-pointer" title="Grant/Deduct Credits">
                  <span class="material-symbols-outlined text-[18px]">bolt</span>
                </button>
                <button type="button" @click="toggleUserStatus(user)" :class="user.is_activated ? 'hover:text-amber-500 hover:bg-amber-500/10' : 'hover:text-emerald-500 hover:bg-emerald-500/10'" class="p-1.5 text-on-surface-variant rounded-lg transition-colors cursor-pointer" :title="user.is_activated ? 'Suspend Account' : 'Activate Account'">
                  <span class="material-symbols-outlined text-[18px]">{{ user.is_activated ? 'block' : 'check_circle' }}</span>
                </button>
                <button type="button" @click="deleteUser(user.id)" class="p-1.5 text-on-surface-variant hover:text-rose-500 hover:bg-rose-500/10 rounded-lg transition-colors cursor-pointer" title="Delete User">
                  <span class="material-symbols-outlined text-[18px]">delete</span>
                </button>
              </td>
            </tr>
            <tr v-if="!loading && users.length === 0">
              <td colspan="6" class="px-6 py-12 text-center text-on-surface-variant">No users found.</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="meta.total_count > meta.limit" class="px-6 py-4 border-t border-outline-variant/40 flex items-center justify-between bg-surface-container-lowest/30">
        <span class="text-sm text-on-surface-variant">
          Showing {{ ((meta.page - 1) * meta.limit) + 1 }} to {{ Math.min(meta.page * meta.limit, meta.total_count) }} of {{ meta.total_count }} users
        </span>
        <div class="flex items-center gap-2">
          <button @click="fetchUsers(meta.page - 1)" :disabled="meta.page === 1" class="px-3 py-1.5 rounded-lg border border-outline-variant/40 text-sm font-medium hover:bg-surface-container disabled:opacity-50 disabled:cursor-not-allowed">Previous</button>
          <button @click="fetchUsers(meta.page + 1)" :disabled="meta.page * meta.limit >= meta.total_count" class="px-3 py-1.5 rounded-lg border border-outline-variant/40 text-sm font-medium hover:bg-surface-container disabled:opacity-50 disabled:cursor-not-allowed">Next</button>
        </div>
      </div>
    </div>

    <!-- Modals -->

    <!-- Create User Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
      <div class="bg-surface border border-outline-variant/60 rounded-2xl max-w-md w-full p-6 shadow-2xl space-y-4 animate-fade-in">
        <div class="flex justify-between items-center border-b border-outline-variant/40 pb-3">
          <h3 class="font-headline text-lg font-bold text-on-surface flex items-center gap-2">
            <span class="material-symbols-outlined text-amber-500">person_add</span> Create New User
          </h3>
          <button @click="showCreateModal = false" class="text-on-surface-variant hover:text-on-surface">
            <span class="material-symbols-outlined">close</span>
          </button>
        </div>
        <form @submit.prevent="submitCreateUser" class="space-y-4 font-body">
          <div>
            <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Full Name</label>
            <input type="text" v-model="formCreate.name" required class="w-full px-3.5 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm focus:border-amber-500 focus:outline-none">
          </div>
          <div>
            <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Email Address</label>
            <input type="email" v-model="formCreate.email" required class="w-full px-3.5 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm focus:border-amber-500 focus:outline-none">
          </div>
          <div>
            <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Password</label>
            <input type="password" v-model="formCreate.password" required class="w-full px-3.5 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm focus:border-amber-500 focus:outline-none">
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Role</label>
              <select v-model="formCreate.role" class="w-full px-3 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm font-bold focus:border-amber-500 focus:outline-none">
                <option value="User">User</option>
                <option value="Admin">Admin</option>
              </select>
            </div>
            <div>
              <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Status</label>
              <select v-model="formCreate.is_activated" class="w-full px-3 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm font-bold focus:border-amber-500 focus:outline-none">
                <option :value="true">Active</option>
                <option :value="false">Suspended</option>
              </select>
            </div>
          </div>
          <div class="flex justify-end gap-3 pt-4 border-t border-outline-variant/40">
            <button type="button" @click="showCreateModal = false" class="px-4 py-2 bg-surface-container-high text-on-surface-variant hover:text-on-surface font-label text-xs font-semibold rounded-xl transition-colors">Cancel</button>
            <button type="submit" :disabled="saving" class="px-5 py-2 bg-amber-500 hover:bg-amber-400 text-on-primary font-label text-xs font-bold rounded-xl disabled:opacity-50">Create Account</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Edit User Modal -->
    <div v-if="showEditModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
      <div class="bg-surface border border-outline-variant/60 rounded-2xl max-w-md w-full p-6 shadow-2xl space-y-4 animate-fade-in">
        <div class="flex justify-between items-center border-b border-outline-variant/40 pb-3">
          <h3 class="font-headline text-lg font-bold text-on-surface flex items-center gap-2">
            <span class="material-symbols-outlined text-amber-500">edit</span> Edit User
          </h3>
          <button @click="showEditModal = false" class="text-on-surface-variant hover:text-on-surface">
            <span class="material-symbols-outlined">close</span>
          </button>
        </div>
        <form @submit.prevent="submitEditUser" class="space-y-4 font-body">
          <div>
            <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Full Name</label>
            <input type="text" v-model="formEdit.name" required class="w-full px-3.5 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm focus:border-amber-500 focus:outline-none">
          </div>
          <div>
            <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Email Address</label>
            <input type="email" v-model="formEdit.email" required class="w-full px-3.5 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm focus:border-amber-500 focus:outline-none">
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Role</label>
              <select v-model="formEdit.role" class="w-full px-3 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm font-bold focus:border-amber-500 focus:outline-none">
                <option value="User">User</option>
                <option value="Admin">Admin</option>
              </select>
            </div>
            <div>
              <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1">Status</label>
              <select v-model="formEdit.is_activated" class="w-full px-3 py-2 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm font-bold focus:border-amber-500 focus:outline-none">
                <option :value="true">Active</option>
                <option :value="false">Suspended</option>
              </select>
            </div>
          </div>
          <div class="flex justify-end gap-3 pt-4 border-t border-outline-variant/40">
            <button type="button" @click="showEditModal = false" class="px-4 py-2 bg-surface-container-high text-on-surface-variant hover:text-on-surface font-label text-xs font-semibold rounded-xl transition-colors">Cancel</button>
            <button type="submit" :disabled="saving" class="px-5 py-2 bg-amber-500 hover:bg-amber-400 text-on-primary font-label text-xs font-bold rounded-xl disabled:opacity-50">Save Changes</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Grant Credits Modal -->
    <div v-if="showCreditModal" class="fixed inset-0 z-50 bg-black/60 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="bg-surface rounded-2xl border border-outline-variant/60 p-6 max-w-md w-full shadow-2xl space-y-4 animate-fade-in">
        <div class="flex justify-between items-center border-b border-outline-variant/40 pb-3">
          <h3 class="font-headline text-lg font-bold text-on-surface flex items-center gap-2">
            <span class="material-symbols-outlined text-amber-500">bolt</span> Grant / Deduct Credits
          </h3>
          <button @click="showCreditModal = false" class="text-on-surface-variant hover:text-on-surface">
            <span class="material-symbols-outlined">close</span>
          </button>
        </div>
        <form @submit.prevent="submitCreditAllocation" class="space-y-4 font-body">
          <div>
            <label class="block font-label text-xs font-bold text-on-surface-variant uppercase mb-1">Target User</label>
            <input type="text" :value="`${formCredit.targetName} (${formCredit.targetEmail})`" readonly class="w-full px-3 py-2 bg-surface-container-low border border-outline-variant/40 rounded-xl text-sm font-semibold text-on-surface">
          </div>
          <div>
            <label class="block font-label text-xs font-bold text-on-surface-variant uppercase mb-1">Credit Amount (+ to Add, - to Deduct)</label>
            <input type="number" v-model="formCredit.amount" placeholder="e.g. 50" required class="w-full px-3 py-2 bg-surface border border-outline-variant/60 rounded-xl text-sm font-bold text-on-surface focus:outline-none focus:border-amber-500">
          </div>
          <div>
            <label class="block font-label text-xs font-bold text-on-surface-variant uppercase mb-1">Reason / Note</label>
            <input type="text" v-model="formCredit.reason" placeholder="e.g. Promotional Bonus" required class="w-full px-3 py-2 bg-surface border border-outline-variant/60 rounded-xl text-sm text-on-surface focus:outline-none focus:border-amber-500">
          </div>
          <div class="flex justify-end gap-3 pt-2">
            <button type="button" @click="showCreditModal = false" class="px-4 py-2 rounded-xl border border-outline-variant/60 text-on-surface-variant font-label text-xs font-semibold hover:bg-surface-container">Cancel</button>
            <button type="submit" :disabled="saving" class="px-5 py-2 rounded-xl bg-amber-500 hover:bg-amber-600 text-on-primary font-label text-xs font-bold">Submit Allocation</button>
          </div>
        </form>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import api from '../../utils/api'
import { useFlash } from '../../composables/useFlash'
import FlashAlert from '../../components/common/FlashAlert.vue'
import dayjs from 'dayjs'

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
  } catch (err) {
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
  } catch (err) {
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
  } catch (err) {
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
  } catch (err) {
    showFlash('Failed to update role', 'error')
    fetchUsers(meta.value.page) // revert
  }
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
  } catch (err) {
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
  } catch (err) {
    showFlash(err.response?.data?.error || 'Failed to allocate credits', 'error')
  } finally {
    saving.value = false
  }
}

async function deleteUser(id) {
  if (!confirm('Are you sure you want to permanently delete this user?')) return
  try {
    await api.delete(`/admin/users/${id}`)
    showFlash('User deleted successfully', 'success')
    fetchUsers(meta.value.page)
  } catch (err) {
    showFlash('Failed to delete user', 'error')
  }
}
</script>
