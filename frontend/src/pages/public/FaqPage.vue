<template>
  <PublicPageLayout 
    title="Frequently Asked Questions" 
    subtitle="Find answers to common questions about Teks-Invoice, billing, and security."
  >
    <div class="space-y-4 not-prose">
      <div 
        v-for="(faq, index) in faqs" 
        :key="index"
        class="glass-card rounded-2xl border border-outline-variant/60 overflow-hidden transition-all duration-300"
      >
        <button 
          @click="toggle(index)" 
          class="w-full text-left px-6 py-4 flex items-center justify-between focus:outline-none hover:bg-surface-container-lowest/50"
        >
          <h3 class="font-headline font-bold text-on-surface text-base">{{ faq.question }}</h3>
          <span class="material-symbols-outlined text-on-surface-variant transition-transform duration-300" :class="{ 'rotate-180': activeIndex === index }">
            expand_more
          </span>
        </button>
        <div 
          v-show="activeIndex === index" 
          class="px-6 pb-4 pt-1 font-body text-sm text-on-surface-variant leading-relaxed"
        >
          {{ faq.answer }}
        </div>
      </div>
    </div>
  </PublicPageLayout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import PublicPageLayout from '../../layouts/PublicPageLayout.vue'

const activeIndex = ref<number | null>(0)

const toggle = (index: number) => {
  activeIndex.value = activeIndex.value === index ? null : index
}

const faqs = [
  {
    question: 'How do invoice credits work?',
    answer: 'Teks-Invoice operates on a pay-as-you-go credit system. Generating a PDF invoice or sending it directly to a client consumes 1 credit. New accounts receive free starter credits, and you can top up anytime. No monthly subscriptions are required.'
  },
  {
    question: 'Can clients view invoices online?',
    answer: 'Yes! Every invoice you create generates a secure, unique public link. You can share this link directly with your client, allowing them to view and download the invoice without needing an account.'
  },
  {
    question: 'Are my credits going to expire?',
    answer: 'No, purchased credits never expire. You can buy a bundle of credits and use them over the course of several months or years without worrying about losing them.'
  },
  {
    question: 'Is my financial data secure?',
    answer: 'Absolutely. We use industry-standard encryption for all data in transit and at rest. We do not sell your data, and your financial information is strictly private to your account.'
  },
  {
    question: 'Can I customize my invoice branding?',
    answer: 'Depending on your credit package tier, you can add custom logos, change accent colors, and adjust the layout of your invoices to match your business branding.'
  }
]
</script>
