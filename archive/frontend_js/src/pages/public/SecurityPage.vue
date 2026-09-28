<template>
  <PublicPageLayout title="Security & Data Integrity" lastUpdated="August 2026">
    <div v-if="customContent" class="whitespace-pre-line text-on-surface font-body text-base">
      {{ customContent }}
    </div>
    <template v-else>
      <p>We take the security of your invoicing data seriously. Here is an overview of our security practices.</p>
      
      <h3>1. Data Encryption</h3>
      <p>All data transmitted between your browser and our servers is encrypted using industry-standard TLS. Sensitive data at rest, such as passwords, are strongly hashed.</p>
      
      <h3>2. Infrastructure Security</h3>
      <p>Our platform is hosted on secure cloud infrastructure with strict access controls and regular security auditing.</p>
      
      <h3>3. Public Invoice Links</h3>
      <p>Public invoices are accessible via unguessable, cryptographically secure URLs. No authentication is required for clients to view their invoice, but the link itself acts as a secure bearer token.</p>
    </template>
  </PublicPageLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import PublicPageLayout from '../../layouts/PublicPageLayout.vue'
import api from '../../utils/api'

const customContent = ref('')

onMounted(async () => {
  try {
    const res = await api.get('/public/settings')
    if (res.data?.settings?.legal_security) {
      customContent.value = res.data.settings.legal_security
    }
  } catch (err) {
    // Fallback to default template
  }
})
</script>
