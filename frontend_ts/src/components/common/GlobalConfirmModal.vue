<template>
  <Teleport to="body">
    <transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="confirmState.isOpen"
        class="fixed inset-0 z-[10000] flex items-center justify-center p-4 bg-slate-950/60 backdrop-blur-sm"
        @click.self="handleCancel"
      >
        <div
          class="w-full max-w-md bg-surface-container-lowest border border-outline-variant/40 rounded-3xl p-6 shadow-2xl space-y-5 animate-fade-in font-body text-on-surface"
        >
          <div class="flex items-start gap-4">
            <div
              class="w-12 h-12 rounded-2xl flex items-center justify-center shrink-0"
              :class="[
                confirmState.type === 'danger'
                  ? 'bg-rose-500/15 text-rose-600 border border-rose-500/20'
                  : 'bg-amber-500/15 text-amber-600 border border-amber-500/20'
              ]"
            >
              <span class="material-symbols-outlined text-[24px]">
                {{ confirmState.type === 'danger' ? 'warning' : 'help' }}
              </span>
            </div>
            <div class="space-y-1">
              <h3 class="font-headline text-lg font-bold text-on-surface">
                {{ confirmState.title }}
              </h3>
              <p class="font-body text-sm text-on-surface-variant leading-relaxed">
                {{ confirmState.message }}
              </p>
            </div>
          </div>

          <div class="flex items-center justify-end gap-3 pt-2">
            <button
              type="button"
              @click="handleCancel"
              class="px-4 py-2.5 rounded-xl font-label text-sm font-bold text-on-surface-variant hover:bg-surface-variant/40 transition-all cursor-pointer"
            >
              {{ confirmState.cancelText }}
            </button>
            <button
              type="button"
              @click="handleConfirm"
              class="px-5 py-2.5 rounded-xl font-label text-sm font-bold shadow-md transition-all cursor-pointer"
              :class="[
                confirmState.type === 'danger'
                  ? 'bg-rose-600 text-white hover:bg-rose-700 shadow-rose-600/20'
                  : 'bg-primary text-on-primary hover:bg-primary/90 shadow-primary/20'
              ]"
            >
              {{ confirmState.confirmText }}
            </button>
          </div>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<script setup lang="ts">
import { useConfirm } from '@/composables/useConfirm'

const { confirmState, handleConfirm, handleCancel } = useConfirm()
</script>
