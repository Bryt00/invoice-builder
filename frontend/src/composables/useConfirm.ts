import { ref } from 'vue'

export type ConfirmType = 'danger' | 'info' | 'warning'

export interface ConfirmOptions {
  title?: string
  message?: string
  confirmText?: string
  cancelText?: string
  type?: ConfirmType
}

export interface ConfirmState {
  isOpen: boolean
  title: string
  message: string
  confirmText: string
  cancelText: string
  type: ConfirmType
  resolve: ((value: boolean) => void) | null
}

const confirmState = ref<ConfirmState>({
  isOpen: false,
  title: 'Confirm Action',
  message: 'Are you sure you want to proceed?',
  confirmText: 'Confirm',
  cancelText: 'Cancel',
  type: 'danger',
  resolve: null,
})

export function useConfirm() {
  function askConfirm(options: ConfirmOptions = {}): Promise<boolean> {
    return new Promise((resolve) => {
      confirmState.value = {
        isOpen: true,
        title: options.title || 'Confirm Action',
        message: options.message || 'Are you sure you want to proceed?',
        confirmText: options.confirmText || 'Confirm',
        cancelText: options.cancelText || 'Cancel',
        type: options.type || 'danger',
        resolve,
      }
    })
  }

  function handleConfirm() {
    if (confirmState.value.resolve) confirmState.value.resolve(true)
    confirmState.value.isOpen = false
  }

  function handleCancel() {
    if (confirmState.value.resolve) confirmState.value.resolve(false)
    confirmState.value.isOpen = false
  }

  return {
    confirmState,
    askConfirm,
    handleConfirm,
    handleCancel,
  }
}
