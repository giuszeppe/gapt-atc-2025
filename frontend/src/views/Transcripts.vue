<template>
  <div v-if="currentStep === 0">
    <selection-container :selection-item="categorySelection" @selectionConfirmed="handleCategorySelection" />
  </div>

  <div v-else-if="currentStep === 1">
    <selection-container :selection-item="simulationSelection" @selectionConfirmed="handleSimulationSelection" />
  </div>

  <div v-else-if="currentStep === 2">
    <selection-container :selection-item="transcriptSelection" @selectionConfirmed="handleTranscriptSelection" />
  </div>
</template>


<script lang="ts" src="./Transcripts.ts">

</script>

<style scoped lang="less">
@import "@/assets/variables.less";

.transcripts-selection {
  padding: 2rem;
  color: white;

  .section {
    margin-bottom: 2rem;

    h2 {
      margin-bottom: 1rem;
    }

    .button-group {
      display: flex;
      gap: 1rem;

      button {
        padding: 0.75rem 1.25rem;
        background-color: @primary-blue;
        border: none;
        border-radius: 10px;
        color: white;
        cursor: pointer;

        &:hover {
          background-color: lighten(@primary-blue, 10%);
        }
      }
    }

    ul {
      list-style: none;
      padding: 0;

      .transcript-item {
        cursor: pointer;
        padding: 0.5rem 1rem;
        border: 1px solid @primary-yellow;
        border-radius: 6px;
        margin: 0.25rem 0;

        &:hover {
          background-color: @primary-yellow;
          color: black;
        }
      }
    }
  }
}
</style>
