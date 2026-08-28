<template>
  <div class="min-h-screen bg-[#FAF9F6] py-10 px-4 sm:px-6 lg:px-8">
    <div class="max-w-[700px] mx-auto">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-[32px] font-bold text-[#2D2D2D] font-serif mb-2">Business Profile Setup</h1>
        <p class="text-[#6B6B6B] text-sm">Configure your company profile to generate professional invoices.</p>
      </div>

      <!-- Flash Messages -->
      <div v-if="flashMessage" :class="`mb-6 p-4 rounded-xl text-sm font-medium flex items-center gap-3 ${flashType === 'success' ? 'bg-[#4B7355]/10 text-[#4B7355] border border-[#4B7355]/20' : 'bg-red-50 text-red-600 border border-red-200'}`">
        <span class="material-symbols-outlined text-[20px]">{{ flashType === 'success' ? 'check_circle' : 'error' }}</span>
        {{ flashMessage }}
      </div>

      <!-- Display View -->
      <div v-if="isProfileComplete && !isEditing" class="bg-[#FAF9F6] border border-[#E5E2DC] rounded-[24px] p-8 shadow-sm">
        <div class="flex items-center gap-5 pb-6 border-b border-[#E5E2DC]">
          <img v-if="profile?.logo_url" :src="profile.logo_url" alt="Logo" class="w-16 h-16 rounded-2xl object-contain border border-[#E5E2DC] bg-white p-1 shadow-sm">
          <div v-else class="w-16 h-16 rounded-2xl bg-[#4B7355]/10 text-[#4B7355] border border-[#4B7355]/20 flex items-center justify-center font-bold text-2xl uppercase">
            {{ authStore.user?.name?.charAt(0) || 'B' }}
          </div>
          <div>
            <h3 class="text-xl font-bold text-[#2D2D2D]">{{ authStore.user?.name }}</h3>
            <p class="text-sm text-[#4B7355] font-medium capitalize">{{ profile?.role }}</p>
            <p v-if="profile?.company_name" class="text-xs text-[#6B6B6B] font-medium mt-1">{{ profile.company_name }}</p>
          </div>
        </div>

        <div class="py-6 grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <span class="block text-[11px] font-bold text-[#6B6B6B] uppercase tracking-wider mb-1">Company Name</span>
            <p class="font-medium text-[#2D2D2D]">{{ profile?.company_name || '-' }}</p>
          </div>
          <div>
            <span class="block text-[11px] font-bold text-[#6B6B6B] uppercase tracking-wider mb-1">Business Type</span>
            <p class="font-medium text-[#2D2D2D]">{{ profile?.business_type || '-' }}</p>
          </div>
          <div class="md:col-span-2">
            <span class="block text-[11px] font-bold text-[#6B6B6B] uppercase tracking-wider mb-1">Business Address</span>
            <p class="font-medium text-[#2D2D2D] whitespace-pre-line">{{ profile?.address || '-' }}</p>
          </div>
        </div>

        <div class="pt-6 border-t border-[#E5E2DC] grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <span class="block text-[11px] font-bold text-[#6B6B6B] uppercase tracking-wider mb-1">Tax ID / TIN</span>
            <p class="font-medium text-[#2D2D2D]">{{ profile?.tax_id || '-' }}</p>
          </div>
          <div>
            <span class="block text-[11px] font-bold text-[#6B6B6B] uppercase tracking-wider mb-1">Default Currency</span>
            <p class="font-medium text-[#2D2D2D] inline-flex items-center gap-1.5 px-3 py-1 bg-white rounded-full border border-[#E5E2DC] text-xs">
              <span class="material-symbols-outlined text-[#4B7355] text-[16px]">payments</span>
              {{ profile?.default_currency || 'USD' }}
            </p>
          </div>
          <div>
            <span class="block text-[11px] font-bold text-[#6B6B6B] uppercase tracking-wider mb-1">Registration Number</span>
            <p class="font-medium text-[#2D2D2D]">{{ profile?.registration_number || '-' }}</p>
          </div>
          <div>
            <span class="block text-[11px] font-bold text-[#6B6B6B] uppercase tracking-wider mb-1">Registration Date</span>
            <p class="font-medium text-[#2D2D2D]">{{ formatDate(profile?.registration_date) }}</p>
          </div>
        </div>

        <div class="pt-8 flex justify-end">
          <button @click="isEditing = true" class="bg-[#FAF9F6] border border-[#E5E2DC] text-[#2D2D2D] rounded-full px-6 py-2.5 font-medium hover:bg-[#F2EFE9] transition-colors flex items-center gap-2 text-sm shadow-sm">
            <span class="material-symbols-outlined text-[18px]">edit</span>
            Edit Profile
          </button>
        </div>
      </div>

      <!-- Edit Form -->
      <div v-else class="bg-[#FAF9F6] border border-[#E5E2DC] rounded-[24px] p-8 shadow-sm">
        <form @submit.prevent="handleSubmit">
          
          <!-- Company Details -->
          <div class="mb-8">
            <h2 class="text-[#2D2D2D] text-lg font-bold mb-3">Company Details</h2>
            <hr class="border-[#E5E2DC] mb-6">
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-5 mb-5">
              <div>
                <label class="block text-[11px] font-bold text-[#6B6B6B] uppercase tracking-wider mb-1.5">Company Name <span class="text-red-500">*</span></label>
                <input type="text" v-model="form.company_name" placeholder="Acme Corp" class="w-full bg-white border border-[#E5E2DC] rounded-full px-4 py-2.5 text-[#2D2D2D] text-sm focus:border-[#4B7355] focus:ring-1 focus:ring-[#4B7355] outline-none transition-all" required>
              </div>
              <div>
                <label class="block text-[11px] font-bold text-[#6B6B6B] uppercase tracking-wider mb-1.5">Your Role / Title <span class="text-red-500">*</span></label>
                <input type="text" v-model="form.role" placeholder="CEO / Manager" class="w-full bg-white border border-[#E5E2DC] rounded-full px-4 py-2.5 text-[#2D2D2D] text-sm focus:border-[#4B7355] focus:ring-1 focus:ring-[#4B7355] outline-none transition-all" required>
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-5 mb-5">
              <div>
                <label class="block text-[11px] font-bold text-[#6B6B6B] uppercase tracking-wider mb-1.5">Your Full Name <span class="text-red-500">*</span></label>
                <input type="text" v-model="form.name" placeholder="John Doe" class="w-full bg-white border border-[#E5E2DC] rounded-full px-4 py-2.5 text-[#2D2D2D] text-sm focus:border-[#4B7355] focus:ring-1 focus:ring-[#4B7355] outline-none transition-all" required>
              </div>
              <div>
                <label class="block text-[11px] font-bold text-[#6B6B6B] uppercase tracking-wider mb-1.5">Business Structure</label>
                <input type="text" v-model="form.business_type" placeholder="LLC, Sole Proprietor..." class="w-full bg-white border border-[#E5E2DC] rounded-full px-4 py-2.5 text-[#2D2D2D] text-sm focus:border-[#4B7355] focus:ring-1 focus:ring-[#4B7355] outline-none transition-all">
              </div>
            </div>

            <div class="mb-5">
              <label class="block text-[11px] font-bold text-[#6B6B6B] uppercase tracking-wider mb-1.5">Business Address <span class="text-red-500">*</span></label>
              <textarea v-model="form.address" placeholder="123 Innovation Way, Suite 100..." rows="3" class="w-full bg-white border border-[#E5E2DC] rounded-[20px] px-4 py-3 text-[#2D2D2D] text-sm focus:border-[#4B7355] focus:ring-1 focus:ring-[#4B7355] outline-none transition-all resize-none" required></textarea>
            </div>

            <div class="mb-5">
              <label class="block text-[11px] font-bold text-[#6B6B6B] uppercase tracking-wider mb-1.5">Brand Logo</label>
              <div class="flex items-center gap-4">
                <img v-if="form.logo_url" :src="form.logo_url" alt="Logo" class="w-12 h-12 rounded-full object-contain border border-[#E5E2DC] bg-white p-1 overflow-hidden">
                <div class="flex-1 relative">
                  <input type="file" id="logo_upload" @change="handleLogoUpload" accept="image/png, image/jpeg, image/jpg, image/webp" class="hidden">
                  <label for="logo_upload" :class="['inline-flex items-center gap-2 px-4 py-2 border border-[#E5E2DC] bg-white rounded-full text-sm font-medium text-[#2D2D2D] hover:bg-[#F2EFE9] transition-colors cursor-pointer', uploadingLogo ? 'opacity-50 pointer-events-none' : '']">
                    <span class="material-symbols-outlined text-[18px]">upload</span>
                    {{ uploadingLogo ? 'Uploading...' : 'Choose File' }}
                  </label>
                  <span class="ml-3 text-xs text-[#6B6B6B] inline-block align-middle mt-1.5">JPG, PNG, or WEBP (Max 10MB)</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Registration & Tax Info -->
          <div class="mb-8">
            <h2 class="text-[#2D2D2D] text-lg font-bold mb-3">Registration & Tax Info</h2>
            <hr class="border-[#E5E2DC] mb-6">

            <div class="grid grid-cols-1 md:grid-cols-2 gap-5 mb-5">
              <div>
                <label class="block text-[11px] font-bold text-[#6B6B6B] uppercase tracking-wider mb-1.5">Tax ID / TIN</label>
                <input type="text" v-model="form.tax_id" placeholder="TAX-123456" class="w-full bg-white border border-[#E5E2DC] rounded-full px-4 py-2.5 text-[#2D2D2D] text-sm focus:border-[#4B7355] focus:ring-1 focus:ring-[#4B7355] outline-none transition-all">
              </div>
              <div>
                <label class="block text-[11px] font-bold text-[#6B6B6B] uppercase tracking-wider mb-1.5">Default Currency</label>
                <div class="relative">
                  <select v-model="form.default_currency" class="w-full bg-white border border-[#E5E2DC] rounded-full pl-4 pr-10 py-2.5 text-[#2D2D2D] text-sm focus:border-[#4B7355] focus:ring-1 focus:ring-[#4B7355] outline-none transition-all appearance-none cursor-pointer">
                    <option v-for="curr in currencies" :key="curr.Code || curr" :value="curr.Code || curr">
                      {{ curr.Flag ? curr.Flag + ' ' : '' }}{{ curr.Code || curr }} ({{ curr.Symbol || '$' }})
                    </option>
                  </select>
                  <span class="material-symbols-outlined absolute right-3 top-1/2 -translate-y-1/2 text-[#6B6B6B] pointer-events-none text-[18px]">expand_more</span>
                </div>
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-5 mb-5">
              <div>
                <label class="block text-[11px] font-bold text-[#6B6B6B] uppercase tracking-wider mb-1.5">Registration Number</label>
                <input type="text" v-model="form.registration_number" placeholder="REG-987654321" class="w-full bg-white border border-[#E5E2DC] rounded-full px-4 py-2.5 text-[#2D2D2D] text-sm focus:border-[#4B7355] focus:ring-1 focus:ring-[#4B7355] outline-none transition-all">
              </div>
              <div>
                <label class="block text-[11px] font-bold text-[#6B6B6B] uppercase tracking-wider mb-1.5">Registration Date</label>
                <input type="date" v-model="form.registration_date" class="w-full bg-white border border-[#E5E2DC] rounded-full px-4 py-2.5 text-[#2D2D2D] text-sm focus:border-[#4B7355] focus:ring-1 focus:ring-[#4B7355] outline-none transition-all">
              </div>
            </div>

            <div class="mb-5">
              <label class="block text-[11px] font-bold text-[#6B6B6B] uppercase tracking-wider mb-1.5">Registered Office Address (If different)</label>
              <textarea v-model="form.registered_address" placeholder="Official Registered Office Address..." rows="2" class="w-full bg-white border border-[#E5E2DC] rounded-[20px] px-4 py-3 text-[#2D2D2D] text-sm focus:border-[#4B7355] focus:ring-1 focus:ring-[#4B7355] outline-none transition-all resize-none"></textarea>
            </div>
          </div>

          <!-- Actions -->
          <div class="flex justify-end gap-3 pt-4">
            <button v-if="isProfileComplete" type="button" @click="isEditing = false" class="bg-transparent text-[#6B6B6B] rounded-full px-6 py-2.5 font-medium hover:bg-[#E5E2DC]/50 transition-colors text-sm">
              Cancel
            </button>
            <button type="submit" :disabled="saving" class="bg-[#4B7355] text-white rounded-full px-8 py-2.5 font-medium hover:bg-[#3D5E45] transition-colors shadow-sm disabled:opacity-50 text-sm">
              {{ saving ? 'Saving...' : 'Save Business Profile' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'
import api from '../../utils/api'

const router = useRouter()
const authStore = useAuthStore()

const saving = ref(false)
const uploadingLogo = ref(false)
const currencies = ref([])
const profile = ref(null)
const isEditing = ref(false)

const flashMessage = ref('')
const flashType = ref('success')

const showFlash = (msg, type = 'success') => {
  flashMessage.value = msg
  flashType.value = type
  setTimeout(() => { flashMessage.value = '' }, 5000)
}

const isProfileComplete = computed(() => {
  return !!authStore.user?.is_profile_complete
})

const form = reactive({
  name: '',
  company_name: '',
  role: '',
  address: '',
  tax_id: '',
  default_currency: 'USD',
  registration_number: '',
  registration_date: '',
  business_type: '',
  registered_address: '',
  logo_url: ''
})

onMounted(async () => {
  try {
    const res = await api.get('/profile')
    if (res.data?.profile) {
      profile.value = res.data.profile
      populateForm(profile.value)
    } else {
      // Default name to user's name if profile not setup
      form.name = authStore.user?.name || ''
    }
    if (res.data?.currencies) {
      currencies.value = res.data.currencies
    }
  } catch (err) {
    console.error(err)
  }
})

function populateForm(p) {
  form.name = authStore.user?.name || ''
  form.company_name = p.company_name || ''
  form.role = p.role || ''
  form.address = p.address || ''
  form.tax_id = p.tax_id || ''
  form.default_currency = p.default_currency || 'USD'
  form.registration_number = p.registration_number || ''
  form.registration_date = p.registration_date ? p.registration_date.substring(0, 10) : ''
  form.business_type = p.business_type || ''
  form.registered_address = p.registered_address || ''
  form.logo_url = p.logo_url || ''
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return dateStr
  return d.toLocaleDateString('en-US', { month: 'short', day: '2-digit', year: 'numeric' })
}

async function handleLogoUpload(event) {
  const file = event.target.files[0]
  if (!file) return

  const formData = new FormData()
  formData.append('logo', file)

  uploadingLogo.value = true
  try {
    const res = await api.post('/profile/logo', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    if (res.data?.logo_url) {
      form.logo_url = res.data.logo_url
      showFlash('Logo uploaded successfully!', 'success')
    }
  } catch (err) {
    showFlash(err.response?.data?.error || 'Failed to upload logo', 'error')
  } finally {
    uploadingLogo.value = false
    // Clear the input so the same file can be uploaded again if needed
    event.target.value = ''
  }
}

async function handleSubmit() {
  saving.value = true
  try {
    await api.put('/profile', form)
    await authStore.fetchCurrentUser() // Update user name/state
    
    showFlash('Business profile updated successfully!', 'success')
    
    const wasEditing = isEditing.value
    isEditing.value = false
    
    // Refetch profile to update display view
    const res = await api.get('/profile')
    if (res.data?.profile) {
      profile.value = res.data.profile
      populateForm(profile.value)
    }

    if (!wasEditing) {
      // If it was their first time setting it up, route them to dashboard
      router.push({ name: 'dashboard' })
    }
  } catch (err) {
    let errorMsg = 'Failed to update profile'
    if (err.response?.data?.error) {
      const apiErr = err.response.data.error
      if (typeof apiErr === 'object') {
        errorMsg = Object.values(apiErr)[0] // Extract the first validation error string
      } else {
        errorMsg = apiErr
      }
    }
    showFlash(errorMsg, 'error')
  } finally {
    saving.value = false
  }
}
</script>
