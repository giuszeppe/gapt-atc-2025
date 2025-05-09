<template>
  <div v-if="props.isVisible" class="modal-overlay">
    <div class="modal-container">
      <h2 class="modal-title">Simulation Completed</h2>
      <div class="modal-actions">
        <button @click="goToTranscripts" class="modal-button primary">
          Consult transcript
        </button>
        <button @click="goToIndex" class="modal-button secondary">
          Home
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
const props = defineProps({
  isVisible: Boolean,
  simulationId: Number,
});
const router = useRouter();

function goToTranscripts() {
  router.push({ path: `transcripts/${props.simulationId}` });
}

function goToIndex() {
  router.push({ name: 'index' });
}
</script>

<style scoped lang="less">
@import "@/assets/variables.less";

.modal-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 999;
}

.modal-container {
  background-color: @primary-blue;
  padding: 2rem;
  border-radius: 12px;
  max-width: 400px;
  width: 90%;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
  text-align: center;
}

.modal-title {
  font-size: 1.5rem;
  margin-bottom: 1rem;
}

.primary {
  background-color: @secondary-blue;
  color: white;
}

.modal-actions {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-top: 1.5rem;
}

.modal-button {
  padding: 0.75rem;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.3s;
}


.modal-button.secondary:hover {
  background-color: #d1d5db;
}
</style>
