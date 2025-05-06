import type { InputType, Role } from '@/@types/types'
import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useStore = defineStore('store', () => {
  const userRole = ref<Role | null>(localStorage.getItem('userRole') as Role || null)
  const inputType = ref<InputType | null>(localStorage.getItem('inputType') as InputType || null)
  const simulationOutline = ref<any>(JSON.parse(localStorage.getItem('simulationOutline')!))
  const simulationInput = ref<any>(JSON.parse(localStorage.getItem('simulationInput')!))
  const lobbyCode = ref<string | null>(localStorage.getItem('lobbyCode') || null)
  const isMultiplayer = ref<boolean>(JSON.parse(localStorage.getItem('isMultiplayer')!) || false)
  const isPlayerInLobby = ref<boolean>(JSON.parse(localStorage.getItem('isPlayerInLobby')!) || false)

  watch(userRole, (newValue) => {
    localStorage.setItem('userRole', newValue ?? '')
  })

  watch(inputType, (newValue) => {
    localStorage.setItem('inputType', newValue ?? '')
  })

  watch(simulationOutline, (newValue) => {
    localStorage.setItem('simulationOutline', JSON.stringify(newValue))
  })

  watch(simulationInput, (newValue) => {
    localStorage.setItem('simulationInput', JSON.stringify(newValue))
  })

  watch(lobbyCode, (newValue) => {
    localStorage.setItem('lobbyCode', newValue ?? '')
  })

  watch(isMultiplayer, (newValue) => {
    localStorage.setItem('isMultiplayer', JSON.stringify(newValue))
  })

  watch(isPlayerInLobby, (newValue) => {
    localStorage.setItem('isPlayerInLobby', JSON.stringify(newValue))
  })

  return { userRole, inputType, simulationOutline, simulationInput, lobbyCode, isMultiplayer, isPlayerInLobby }
})
