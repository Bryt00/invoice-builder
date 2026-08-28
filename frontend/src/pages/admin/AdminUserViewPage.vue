<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6 animate-fade-in">
    <FlashAlert />

    <div v-if="loading" class="flex items-center justify-center min-h-[400px]">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-amber-500"></div>
    </div>

    <template v-else-if="user">
      <!-- Page Header & Back Button -->
      <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <router-link to="/user/admin/users" class="inline-flex items-center text-sm font-label text-on-surface-variant hover:text-amber-500 transition-colors mb-2">
            <span class="material-symbols-outlined text-[18px] mr-1">arrow_back</span>
            Back to Directory
          </router-link>
          <h1 class="font-headline text-2xl sm:text-3xl font-extrabold text-on-surface flex items-center gap-3">
            <span class="w-10 h-10 rounded-full bg-amber-500/20 text-amber-600 border border-amber-500/30 flex items-center justify-center font-bold text-lg shrink-0">
              {{ user.name.charAt(0).toUpperCase() }}
            </span>
            {{ user.name }}'s Profile
          </h1>
        </div>
        <div class="flex items-center gap-2">
          <span v-if="user.is_activated" class="px-3 py-1.5 rounded-full text-xs font-semibold bg-emerald-500/10 text-emerald-600 border border-emerald-500/20 inline-flex items-center gap-1.5">
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span> Active Account
          </span>
          <span v-else class="px-3 py-1.5 rounded-full text-xs font-semibold bg-amber-500/10 text-amber-600 border border-amber-500/20 inline-flex items-center gap-1.5">
            <span class="w-1.5 h-1.5 rounded-full bg-amber-500"></span> Suspended Account
          </span>
        </div>
      </div>

      <!-- Bento Grid for User Stats -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mt-6">
        <!-- Profile Info -->
        <div class="glass-card rounded-2xl p-6 border border-outline-variant/60 shadow-xs">
          <h3 class="font-label text-xs uppercase tracking-wider font-bold text-on-surface-variant mb-4">Account Details</h3>
          <div class="space-y-4 font-body text-sm">
            <div>
              <p class="text-on-surface-variant text-xs mb-0.5">Email Address</p>
              <p class="text-on-surface font-semibold">{{ user.email }}</p>
            </div>
            <div>
              <p class="text-on-surface-variant text-xs mb-0.5">Role</p>
              <p class="text-on-surface font-semibold">{{ user.role.name }}</p>
            </div>
            <div>
              <p class="text-on-surface-variant text-xs mb-0.5">Registered On</p>
              <p class="text-on-surface font-semibold">{{ formatDate(user.created_at) }}</p>
            </div>
            <div>
              <p class="text-on-surface-variant text-xs mb-0.5">Profile Completed</p>
              <p class="text-on-surface font-semibold">{{ user.is_profile_complete ? 'Yes' : 'No' }}</p>
            </div>
          </div>
        </div>

        <!-- Credits Balance -->
        <div class="glass-card rounded-2xl p-6 border border-outline-variant/60 shadow-xs flex flex-col justify-between">
          <div>
            <div class="w-10 h-10 rounded-xl bg-amber-500/10 text-amber-600 border border-amber-500/20 flex items-center justify-center mb-4">
              <span class="material-symbols-outlined text-2xl">bolt</span>
            </div>
            <p class="text-xs text-on-surface-variant font-label uppercase tracking-wider font-bold mb-1">Available Credits</p>
            <h3 class="font-headline text-4xl font-black text-on-surface">{{ user.credits || 0 }}</h3>
          </div>
          <button class="mt-4 w-full py-2 bg-surface-container-high hover:bg-amber-500 hover:text-on-primary text-on-surface transition-colors rounded-xl font-label text-xs font-bold border border-outline-variant/40">
            View Credit History
          </button>
        </div>

        <!-- Actions -->
        <div class="glass-card rounded-2xl p-6 border border-outline-variant/60 shadow-xs flex flex-col justify-between">
          <h3 class="font-label text-xs uppercase tracking-wider font-bold text-on-surface-variant mb-4">Quick Actions</h3>
          <div class="space-y-2">
             <button class="w-full flex items-center justify-between px-4 py-2.5 bg-surface-container hover:bg-surface-container-high border border-outline-variant/40 rounded-xl font-label text-sm text-on-surface transition-colors">
               <div class="flex items-center gap-2">
                 <span class="material-symbols-outlined text-[18px] text-amber-500">lock_reset</span>
                 Force Password Reset
               </div>
               <span class="material-symbols-outlined text-[16px] text-on-surface-variant">chevron_right</span>
             </button>
             <button class="w-full flex items-center justify-between px-4 py-2.5 bg-surface-container hover:bg-surface-container-high border border-outline-variant/40 rounded-xl font-label text-sm text-on-surface transition-colors">
               <div class="flex items-center gap-2">
                 <span class="material-symbols-outlined text-[18px] text-amber-500">mail</span>
                 Send Notification Email
               </div>
               <span class="material-symbols-outlined text-[16px] text-on-surface-variant">chevron_right</span>
             </button>
             <button class="w-full flex items-center justify-between px-4 py-2.5 bg-rose-500/10 hover:bg-rose-500/20 border border-rose-500/20 rounded-xl font-label text-sm text-rose-600 transition-colors">
               <div class="flex items-center gap-2">
                 <span class="material-symbols-outlined text-[18px]">delete_forever</span>
                 Delete Account Data
               </div>
               <span class="material-symbols-outlined text-[16px]">chevron_right</span>
             </button>
          </div>
        </div>
      </div>
    </template>

    <div v-else class="text-center py-12">
      <p class="text-on-surface-variant">User not found.</p>
      <router-link to="/user/admin/users" class="text-amber-500 hover:underline mt-2 inline-block">Return to Directory</router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../../utils/api'
import { useFlash } from '../../composables/useFlash'
import FlashAlert from '../../components/common/FlashAlert.vue'
import dayjs from 'dayjs'

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
  } catch (err) {
    showFlash('Failed to load user profile', 'error')
  } finally {
    loading.value = false
  }
})

function formatDate(date) {
  return dayjs(date).format('MMMM DD, YYYY at h:mm A')
}
</script>
