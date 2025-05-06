<template>
  <div class="chat-container">
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
.chat-container {
  width: 50%;
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
</style>