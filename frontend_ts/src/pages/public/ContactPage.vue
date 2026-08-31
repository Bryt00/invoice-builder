<template>
  <PublicPageLayout 
    title="Contact Support" 
    subtitle="Need help with your invoices or have a question about credits? We're here to help."
  >
    <div class="flex flex-col md:flex-row gap-8 items-start">
      
      <div class="w-full md:w-1/3 space-y-6">
        <div class="glass-card rounded-2xl p-6 border border-outline-variant/60">
          <h3 class="font-headline font-bold text-lg text-on-surface mb-2 flex items-center gap-2">
            <span class="material-symbols-outlined text-primary">mail</span> Email Us
          </h3>
          <p class="text-sm text-on-surface-variant font-medium">
            <a :href="`mailto:${supportEmail}`" class="text-primary hover:underline">{{ supportEmail }}</a>
          </p>
          <p class="text-xs text-on-surface-variant mt-2">We typically reply within 24 hours.</p>
        </div>
        
        <div class="glass-card rounded-2xl p-6 border border-outline-variant/60">
          <h3 class="font-headline font-bold text-lg text-on-surface mb-2 flex items-center gap-2">
            <span class="material-symbols-outlined text-primary">forum</span> FAQ
          </h3>
          <p class="text-sm text-on-surface-variant mb-3">Many common questions are answered in our FAQ section.</p>
          <router-link to="/faq" class="text-sm font-semibold text-primary hover:underline flex items-center gap-1">
            Read the FAQ <span class="material-symbols-outlined text-[16px]">arrow_right_alt</span>
          </router-link>
        </div>
      </div>

      <div class="w-full md:w-2/3 glass-card rounded-2xl p-6 md:p-8 border border-outline-variant/60">
        <h3 class="font-headline font-bold text-xl text-on-surface mb-6">Send us a message</h3>
        <form @submit.prevent="submitForm" class="space-y-4">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="flex flex-col gap-1.5">
              <label class="font-label text-xs font-semibold text-on-surface">Your Name</label>
              <input type="text" class="px-3.5 py-2.5 rounded-xl bg-surface-container-low border border-outline-variant/60 text-sm focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-all text-on-surface" required placeholder="John Doe">
            </div>
            <div class="flex flex-col gap-1.5">
              <label class="font-label text-xs font-semibold text-on-surface">Email Address</label>
              <input type="email" class="px-3.5 py-2.5 rounded-xl bg-surface-container-low border border-outline-variant/60 text-sm focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-all text-on-surface" required placeholder="john@example.com">
            </div>
          </div>
          <div class="flex flex-col gap-1.5">
            <label class="font-label text-xs font-semibold text-on-surface">Subject</label>
            <input type="text" class="px-3.5 py-2.5 rounded-xl bg-surface-container-low border border-outline-variant/60 text-sm focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-all text-on-surface" required placeholder="How can we help?">
          </div>
          <div class="flex flex-col gap-1.5">
            <label class="font-label text-xs font-semibold text-on-surface">Message</label>
            <textarea rows="5" class="px-3.5 py-2.5 rounded-xl bg-surface-container-low border border-outline-variant/60 text-sm focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-all text-on-surface resize-y" required placeholder="Your message..."></textarea>
          </div>
          <button type="submit" class="w-full py-3 rounded-xl font-label text-sm font-bold bg-primary text-on-primary hover:bg-on-primary-fixed-variant transition-colors shadow-md">
            Send Message
          </button>
        </form>
      </div>

    </div>
  </PublicPageLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import PublicPageLayout from '../../layouts/PublicPageLayout.vue'
import api from '@/utils/api'

const supportEmail = ref('support@teks-invoice.com')

onMounted(async () => {
  try {
    const res = await api.get('/public/settings')
    if (res.data?.settings?.support_email) {
      supportEmail.value = res.data.settings.support_email
    }
  } catch (err) {
    // Keep default
  }
})

const submitForm = () => {
  alert('Thanks for your message! This is a demo form, but in a real app it would submit to our backend.')
}
</script>
