<template>
  <PublicPageLayout title="Refund Policy" lastUpdated="August 2026">
    <div v-if="customContent" class="whitespace-pre-line text-on-surface font-body text-base">
      {{ customContent }}
    </div>
    <template v-else>
      <p>Thank you for using Teks-Invoice. Please read our refund policy carefully regarding credit purchases.</p>
      
      <h3>1. General Policy</h3>
      <p>Because Teks-Invoice uses a pay-as-you-go credit system, all credit purchases are generally final and non-refundable.</p>
      
      <h3>2. Exceptions</h3>
      <p>We may issue a refund in the following scenarios:</p>
      <ul>
        <li>If there was a billing error or duplicate charge.</li>
        <li>If the platform was completely inaccessible for an extended period, preventing you from using purchased credits.</li>
      </ul>
      
      <h3>3. Contact Us</h3>
      <p>If you believe you are entitled to a refund, please contact support within 7 days of your purchase with your receipt details.</p>
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
    if (res.data?.settings?.legal_refund) {
      customContent.value = res.data.settings.legal_refund
    }
  } catch (err) {
    // Fallback to default template
  }
})
</script>
