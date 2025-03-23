<template>
  <div>
    <selection-container v-if="currentStep < steps.length" :selection-item="steps[currentStep]"
      @selectionConfirmed="handleSelection"></selection-container>
  </div>
</template>

<script lang="ts" src="./Index.ts">
</script>
