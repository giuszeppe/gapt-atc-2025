<template>
  <div class="simulation-container">
    <div class="left-panel">
      <div class="chat-container">
        <div v-for="(message, index) in transcript.messages" :key="'left-' + index"
          :class="['message-row', message.role === 'aircraft' ? 'right' : 'left']">
          <div class="role-avatar">
            <font-awesome-icon
              :icon="message.role === 'aircraft' ? 'fa-solid fa-plane' : 'fa-solid fa-tower-observation'" />
          </div>
          <div class="chat-bubble">
            <div v-html="message.content"></div>
          </div>
        </div>
      </div>
    </div>
    <div class="right-panel">
      <div class="chat-container">
        <div v-for="(step, index) in simulationOutline" :key="'right-' + index" class="message-row"
          :class="step.role === 'aircraft' ? 'right' : 'left'">
          <div class="role-avatar">
            <font-awesome-icon
              :icon="step.role === 'aircraft' ? 'fa-solid fa-plane' : 'fa-solid fa-tower-observation'" />
          </div>
          <div class="chat-bubble">{{ step.text }}</div>
        </div>
      </div>
      <button class="option-button mt-10 mb-10" @click="goHome">Go back home</button>
    </div>
  </div>
</template>

<script src="./GetTranscript.ts" lang="ts"></script>

<style scoped lang="less">
@import "@/assets/variables.less";

.go-back-button {
  width: 200px;
  height: 50px;
  background-color: transparent;
  border: 3px solid @primary-yellow;
  color: @primary-yellow;
  border-radius: 5px;
  font-size: 18px;
  cursor: pointer;
  margin-bottom: 10px;
}

.simulation-container {
  display: flex;
  height: 100vh;
  width: 100vw;
  color: white;
  text-transform: none;
}

.chat-message {
  display: flex;
  max-width: 85%;
}

.chat-message.left {
  justify-content: flex-start;
}

.chat-message.right {
  justify-content: flex-end;
  margin-left: auto;
}

.input-container {
  width: 90%;
  padding: 1rem;
  display: flex;
  gap: 0.5rem;
}

.input-container input {
  flex: 1;
  padding: 0.5rem;
  border-radius: 8px;
  border: none;
  background-color: #444;
  color: white;
}

.input-container button {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 8px;
  background-color: @secondary-blue;
  color: white;
  cursor: pointer;
}

.input-container button:disabled,
.input-container input:disabled {
  background-color: #333;
  background-color: #555;
  cursor: not-allowed;
}

.block-container,
.selected-blocks {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.block {
  background-color: #1e3a8a;
  color: white;
  padding: 0.5rem 1rem;
  cursor: pointer;
  user-select: none;
  text-align: center;
}

.block-mode {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  gap: 1rem;
}


.block.selected {
  background-color: #2563eb;
}

.role-icon {
  margin-right: 0.5rem;
  color: #ccc;
}

.invisible {
  visibility: hidden;
}
</style>