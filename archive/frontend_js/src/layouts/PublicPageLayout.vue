<template>
  <div class="min-h-screen bg-background text-on-surface flex flex-col selection:bg-primary-container/40 selection:text-primary">
    <!-- Navbar -->
    <AppNavbar />

    <main class="flex-1 flex flex-col items-center">
      <!-- Header Section -->
      <section class="w-full bg-surface-container-lowest/70 border-b border-outline-variant/40 pt-24 pb-16 px-4 sm:px-6 lg:px-8">
        <div class="max-w-4xl mx-auto text-center space-y-4">
          <h1 class="font-headline text-3xl sm:text-4xl md:text-5xl font-extrabold text-on-surface">{{ title }}</h1>
          <p v-if="subtitle" class="font-body text-base sm:text-lg text-on-surface-variant max-w-2xl mx-auto">
            {{ subtitle }}
          </p>
          <p v-if="lastUpdated" class="font-body text-sm text-outline font-medium mt-4">
            Last Updated: {{ lastUpdated }}
          </p>
        </div>
      </section>

      <!-- Content Section -->
      <section class="w-full max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12 md:py-16">
        <div class="prose prose-sm sm:prose-base prose-green max-w-none text-on-surface-variant 
          prose-headings:text-on-surface prose-headings:font-headline prose-headings:font-bold
          prose-a:text-primary prose-a:no-underline hover:prose-a:underline
          prose-strong:text-on-surface prose-strong:font-bold">
          <slot></slot>
        </div>
      </section>
    </main>

    <!-- Footer -->
    <AppFooter />
  </div>
</template>

<script setup lang="ts">
import AppNavbar from '../components/common/AppNavbar.vue'
import AppFooter from '../components/common/AppFooter.vue'

defineProps<{
  title: string
  subtitle?: string
  lastUpdated?: string
}>()
</script>

<style>
/* Add some basic prose styling rules here just in case tailwindcss/typography plugin isn't active */
.prose h1 { @apply text-3xl mb-6 mt-8; }
.prose h2 { @apply text-2xl mb-4 mt-8; }
.prose h3 { @apply text-xl mb-4 mt-6; }
.prose p { @apply mb-4 leading-relaxed; }
.prose ul { @apply list-disc list-inside mb-4 space-y-2; }
.prose ol { @apply list-decimal list-inside mb-4 space-y-2; }
.prose li { @apply leading-relaxed; }
.prose strong { @apply font-bold; }
</style>
