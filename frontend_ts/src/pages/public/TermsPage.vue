<template>
  <PublicPageLayout title="Terms of Service" lastUpdated="August 2026">
    <div v-if="customContent" class="whitespace-pre-line text-on-surface font-body text-base">
      {{ customContent }}
    </div>
    <template v-else>
      <p>Welcome to Teks-Invoice. By using our service, you agree to these terms.</p>
      
      <h3>1. Account Terms</h3>
      <p>You must provide a valid email address and maintain account security. You are responsible for all activity that occurs under your account.</p>
      
      <h3>2. Payment & Credits</h3>
      <p>Credits purchased are non-refundable unless specified by law. Invoice generation consumes credits as per the pricing structure.</p>
      
      <h3>3. Service Availability</h3>
      <p>We strive for 99.9% uptime but do not guarantee uninterrupted access. We reserve the right to perform maintenance as necessary.</p>
      
      <h3>4. Prohibited Uses</h3>
      <p>You may not use our service for any illegal activities or to generate fraudulent invoices.</p>
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
    if (res.data?.settings?.legal_terms) {
      customContent.value = res.data.settings.legal_terms
    }
  } catch (err) {
    // Fallback to default template
  }
})
</script>
