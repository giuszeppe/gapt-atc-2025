import type { InputType, Role } from '@/@types/types'
import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useStore = defineStore('store', () => {
  const userRole = ref<Role>(localStorage.getItem('userRole') as Role || null!)
  const inputType = ref<InputType>(localStorage.getItem('inputType') as InputType || null!)

  watch(userRole, (newValue) => {
    localStorage.setItem('userRole', newValue)
  })

  watch(inputType, (newValue) => {
    localStorage.setItem('inputType', newValue)
  })

  return { userRole, inputType }
})
