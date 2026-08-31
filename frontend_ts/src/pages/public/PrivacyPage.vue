<template>
  <PublicPageLayout title="Privacy Policy" lastUpdated="August 2026">
    <div v-if="customContent" class="whitespace-pre-line text-on-surface font-body text-base">
      {{ customContent }}
    </div>
    <template v-else>
      <p>Your privacy is important to us. This Privacy Policy explains how Teks-Invoice collects, uses, and safeguards your information.</p>
      
      <h3>1. Information We Collect</h3>
      <p>We collect personal information such as your name, email, and business details when you register. We also collect client information you enter for invoices.</p>
      
      <h3>2. How We Use Information</h3>
      <p>Your information is used strictly to provide the invoicing service, process payments, and send important service updates.</p>
      
      <h3>3. Data Sharing</h3>
      <p>We do not sell your data. We may share data with trusted third-party service providers (like payment processors) solely to operate the platform.</p>
      
      <h3>4. Your Rights</h3>
      <p>You have the right to access, correct, or delete your personal data. Contact our support team to exercise these rights.</p>
    </template>
  </PublicPageLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import PublicPageLayout from '../../layouts/PublicPageLayout.vue'
import api from '@/utils/api'

const customContent = ref('')

onMounted(async () => {
  try {
    const res = await api.get('/public/settings')
    if (res.data?.settings?.legal_privacy) {
      customContent.value = res.data.settings.legal_privacy
    }
  } catch (err) {
    // Fallback to default template
  }
})
</script>
