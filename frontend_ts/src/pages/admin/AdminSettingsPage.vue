<template>
  <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6 animate-fade-in">
    <FlashAlert />

    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 border-b border-outline-variant/40 pb-6">
      <div>
        <h1 class="font-headline text-2xl sm:text-3xl font-extrabold text-on-surface flex items-center gap-2.5">
          <span class="material-symbols-outlined text-amber-500 text-[28px]">settings</span>
          System Settings
        </h1>
        <p class="font-body text-sm text-on-surface-variant mt-1">Configure global application parameters, policies, and maintenance mode.</p>
      </div>
      <button @click="saveSettings" :disabled="saving" class="inline-flex items-center gap-2 px-6 py-2.5 bg-amber-500 hover:bg-amber-400 text-on-primary font-label text-sm font-bold rounded-xl shadow-sm transition-all cursor-pointer disabled:opacity-50">
        <span v-if="saving" class="animate-spin material-symbols-outlined text-[18px]">progress_activity</span>
        <span v-else class="material-symbols-outlined text-[18px]">save</span>
        <span>Save Changes</span>
      </button>
    </div>

    <div v-if="loading" class="flex items-center justify-center min-h-[300px]">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-amber-500"></div>
    </div>

    <div v-else class="space-y-8">
      
      <!-- Core Configuration -->
      <section class="glass-card rounded-2xl border border-outline-variant/60 p-6 shadow-xs">
        <h3 class="font-headline text-lg font-bold text-on-surface border-b border-outline-variant/40 pb-3 mb-5 flex items-center gap-2">
          <span class="material-symbols-outlined text-amber-500">manufacturing</span> Core Configuration
        </h3>
        
        <div class="space-y-6">
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 p-4 rounded-xl bg-surface-container-lowest/50 border border-outline-variant/30">
            <div>
              <p class="font-label text-sm font-bold text-on-surface">Maintenance Mode</p>
              <p class="text-xs text-on-surface-variant mt-0.5">Disables access for all non-admin users across the platform.</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input type="checkbox" v-model="form.maintenance_mode" class="sr-only peer">
              <div class="w-11 h-6 bg-surface-container-high peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-amber-500"></div>
            </label>
          </div>

          <div>
            <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1.5">Support Contact Email</label>
            <input type="email" v-model="form.support_email" class="w-full max-w-md px-4 py-2.5 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm focus:border-amber-500 focus:outline-none">
            <p class="text-[11px] text-on-surface-variant mt-1.5">This email will be displayed on support pages and used for system notifications.</p>
          </div>

          <div>
            <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1.5">Default Signup Credits (Bonus)</label>
            <input type="number" v-model="form.default_signup_credits" class="w-full max-w-[120px] px-4 py-2.5 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm font-bold focus:border-amber-500 focus:outline-none">
            <p class="text-[11px] text-on-surface-variant mt-1.5">Number of free credits automatically granted to users upon registration.</p>
          </div>
        </div>
      </section>

      <!-- Legal & Policies -->
      <section class="glass-card rounded-2xl border border-outline-variant/60 p-6 shadow-xs">
        <h3 class="font-headline text-lg font-bold text-on-surface border-b border-outline-variant/40 pb-3 mb-5 flex items-center gap-2">
          <span class="material-symbols-outlined text-amber-500">policy</span> Legal & Policies (Markdown)
        </h3>
        
        <div class="space-y-6">
          <div>
            <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1.5">Terms of Service</label>
            <textarea v-model="form.legal_terms" rows="4" class="w-full px-4 py-3 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm font-mono focus:border-amber-500 focus:outline-none"></textarea>
          </div>
          <div>
            <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1.5">Privacy Policy</label>
            <textarea v-model="form.legal_privacy" rows="4" class="w-full px-4 py-3 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm font-mono focus:border-amber-500 focus:outline-none"></textarea>
          </div>
          <div>
            <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1.5">Refund Policy</label>
            <textarea v-model="form.legal_refund" rows="4" class="w-full px-4 py-3 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm font-mono focus:border-amber-500 focus:outline-none"></textarea>
          </div>
          <div>
            <label class="block font-label text-xs font-bold uppercase text-on-surface-variant mb-1.5">Security Policy</label>
            <textarea v-model="form.legal_security" rows="4" class="w-full px-4 py-3 rounded-xl bg-surface-container border border-outline-variant/60 text-on-surface text-sm font-mono focus:border-amber-500 focus:outline-none"></textarea>
          </div>
        </div>
      </section>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import api from '@/utils/api'
import { useFlash } from '@/composables/useFlash'
import FlashAlert from '../../components/common/FlashAlert.vue'

const { showFlash } = useFlash()

const loading = ref(true)
const saving = ref(false)

const form = reactive({
  maintenance_mode: false,
  support_email: '',
  default_signup_credits: '0',
  legal_terms: '',
  legal_privacy: '',
  legal_refund: '',
  legal_security: ''
})

onMounted(() => {
  fetchSettings()
})

async function fetchSettings() {
  loading.value = true
  try {
    const res = await api.get('/admin/settings')
    const s = res.data.settings
    if (s) {
      form.maintenance_mode = s.MaintenanceMode === true || s.MaintenanceMode === 'true'
      form.support_email = s.SupportContactEmail || ''
      form.default_signup_credits = s.DefaultSignupBonus || '0'
      form.legal_terms = s.LegalTerms || ''
      form.legal_privacy = s.LegalPrivacy || ''
      form.legal_refund = s.LegalRefund || ''
      form.legal_security = s.LegalSecurity || ''
    }
  } catch (err: any) {
    showFlash('Failed to load system settings', 'error')
  } finally {
    loading.value = false
  }
}

async function saveSettings() {
  saving.value = true
  try {
    await api.put('/admin/settings', {
      maintenance_mode: form.maintenance_mode ? 'true' : 'false',
      support_email: form.support_email,
      default_signup_credits: String(form.default_signup_credits),
      legal_terms: form.legal_terms,
      legal_privacy: form.legal_privacy,
      legal_refund: form.legal_refund,
      legal_security: form.legal_security
    })
    showFlash('System settings updated successfully', 'success')
  } catch (err: any) {
    showFlash('Failed to update system settings', 'error')
  } finally {
    saving.value = false
  }
}
</script>
