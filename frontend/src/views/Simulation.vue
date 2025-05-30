<template>
  <div class="simulation-container">
    <div class="left-panel">
      <div class="chat-container">
        <div v-for="(message, index) in leftPanelMessages" :key="'left-' + index"
          :class="['message-row', message.role === userRole ? 'right' : 'left']">
          <div class="role-avatar">
            <font-awesome-icon
              :icon="message.role === 'aircraft' ? 'fa-solid fa-plane' : 'fa-solid fa-tower-observation'" />
          </div>
          <div class="chat-bubble">
            <div v-html="message.content"></div>
          </div>
        </div>
      </div>

      <VoiceVisualizer :volume="volume" v-if="inputType == 'speech' && isListening" />

      <div class="input-container">
        <template v-if="inputType === 'text' || inputType === 'speech'">
          <input v-model="playerInput" @keyup.enter="handlePlayerInput"
            placeholder="Type your message and press Enter..." />
          <button v-if="inputType === 'speech'" :disabled="!isUserTurn" @mousedown="startListening"
            @mouseup="stopListening" @mouseleave="stopListening">
            <font-awesome-icon icon="fa-solid fa-microphone" />
            {{ isListening ? 'Release to stop' : 'Push to talk' }}
          </button>
          <button @click="handlePlayerInput" :disabled="!isUserTurn">Send</button>
        </template>

        <template v-if="inputType == 'block'">
          <div class="block-mode">
            <div class="selected-blocks">
              <div v-for="(selected, index) in selectedWords" :key="'selected-' + selected.originalIndex"
                class="block selected" draggable="true" @dragstart="onDragStart(index)" @dragover.prevent
                @drop="onDrop(index)" @click="deselectWord(index)">
                {{ selected.word }}
              </div>
            </div>

            <div class="block-container">
              <div v-for="(word, index) in wordBlocks" :key="index" class="block"
                :class="{ invisible: selectedWordsIndexes.includes(index) }" @click="selectWord(index)">
                {{ word }}
              </div>
            </div>
            <button class="submit-button" @click="handlePlayerInput" :disabled="!isUserTurn">Send</button>
          </div>
        </template>

      </div>
    </div>

    <div class="right-panel">
      <div v-if="lobbyCode">LOBBY CODE: {{ lobbyCode }}</div>
      <div class="chat-container">
        <div v-for="(step, index) in rightPanelSteps" :key="'right-' + index" class="message-row"
          :class="step.role === userRole ? 'right' : 'left'">
          <div class="role-avatar">
            <font-awesome-icon
              :icon="step.role === 'aircraft' ? 'fa-solid fa-plane' : 'fa-solid fa-tower-observation'" />
          </div>
          <div class="chat-bubble">{{ step.text }}</div>
        </div>

      </div>
      <button class="option-button mt-10 mb-10" @click="goHome">Go back home</button>
    </div>

    <SimulationEndModal :is-visible="showEndModal" :simulation-id="simulationId" />
  </div>
</template>


<script src="./Simulation.ts" lang="ts"></script>

<style scoped lang="less">
@import "@/assets/variables.less";

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
