import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useStore = defineStore('store', () => {
  const userRole = ref<'aircraft' | 'tower'>(null!)
  const inputType = ref<'block' | 'text' | 'speech'>(null!)

  return { userRole, inputType }
})
