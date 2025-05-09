<template>
  <div class="chat-container transcript-container">
    <div v-for="(message, index) in transcript.messages" :key="'left-' + index"
      :class="['chat-message', message.role === 'aircraft' ? 'right' : 'left']">
      <div class="chat-bubble">
        <div v-if="message.role === 'aircraft'" v-html="message.content"></div>
        <div v-else v-html="message.content"></div>
      </div>
    </div>
  </div>
</template>

<script src="./GetTranscript.ts" lang="ts"></script>

<style scoped>
.transcript-container {
  width: 60%;
  height: 100%;
  border-left: 1px solid #003750;
  border-right: 1px solid #003750;
  color: white;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 1rem;
  padding: 0 3rem;
  overflow-y: auto;
}
</style>