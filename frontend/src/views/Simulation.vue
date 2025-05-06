<template>
  <div class="simulation-container">
    <div class="left-panel">
      <div class="chat-container">
        <div v-for="(message, index) in leftPanelMessages" :key="'left-' + index"
          :class="['chat-message', message.role === userRole ? 'right' : 'left']">
          <div class="chat-bubble">
            <div v-if="message.role === userRole" v-html="message.content"></div>
            <div v-else>{{ message.text }}</div>
          </div>
        </div>
      </div>

      <VoiceVisualizer :volume="volume" v-if="inputType == 'speech' && isListening" />

      <div class="input-container">
        <input v-model="playerInput" @keyup.enter="handlePlayerInput"
          placeholder="Type your message and press Enter..." />
        <button @click="toggleListening" v-if="inputType == 'speech'">
          <font-awesome-icon icon="fa-solid fa-microphone" />
          {{ isListening ? 'Stop' : 'Speak' }}
        </button>
        <button @click="handlePlayerInput" :disabled="!isUserTurn">Send</button>
      </div>
    </div>

    <div class="right-panel">
      <div class="chat-container">
        <div v-for="(step, index) in rightPanelSteps" :key="'right-' + index"
          :class="['chat-message', step.role === userRole ? 'right' : 'left']">
          <div class="chat-bubble">
            {{ step.text }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>


<script src="./Simulation.ts" lang="ts"></script>

<style scoped>
.speech-section {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(6px);
  padding: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.speech-visualizer {
  display: flex;
  gap: 4px;
  width: 100%;
  height: 100%;
  align-items: flex-end;
}

.bar {
  flex: 1;
  background: linear-gradient(180deg, #4fc3f7, #0288d1);
  border-radius: 4px;
  transition: height 0.1s ease;
  box-shadow: 0 0 6px rgba(0, 183, 255, 0.7);
}

.simulation-container {
  display: flex;
  height: 100vh;
  width: 100vw;
  color: white;
  text-transform: none;
}

.left-panel {
  width: 60%;
  display: flex;
  flex-direction: column;
  background-color: #1e1e1e;
}

.right-panel {
  width: 40%;
  display: flex;
  flex-direction: column;
  padding: 1rem;
  overflow-y: auto;
  background-color: #121212;
}

.chat-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1rem;
  overflow-y: auto;
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

.chat-bubble {
  background-color: rgb(96, 96, 96);
  padding: 0.8rem 1rem;
  border-radius: 1rem;
  word-break: break-word;
}

.chat-message.right .chat-bubble {
  background-color: #003750;
}

.input-container {
  padding: 1rem;
  display: flex;
  gap: 0.5rem;
  background-color: #2c2c2c;
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
  background-color: #007acc;
  color: white;
  cursor: pointer;
}

.input-container button:disabled,
.input-container input:disabled {
  background-color: #333;
  background-color: #555;
  cursor: not-allowed;
}
</style>
