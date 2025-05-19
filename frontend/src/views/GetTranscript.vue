<template>
  <div class="chat-container transcript-container">
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
</template>

<script src="./GetTranscript.ts" lang="ts"></script>

<style scoped lang="less">
@import "@/assets/variables.less";

.transcript-container {
  width: 60%;
  height: 100%;
  border-left: 1px solid @primary-blue;
  border-right: 1px solid @primary-blue;
  color: white;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 1rem;
  padding: 0 3rem;
  overflow-y: auto;
}

.message-row {
  display: flex;
  align-items: flex-end;
  gap: 0.75rem;
  max-width: 100%;

  &.right {
    flex-direction: row-reverse;

    .role-avatar {
      background-color: @primary-blue;
      margin-left: 0.5rem;
    }

    .chat-bubble {
      background-color: @primary-blue;
    }
  }

  &.left {
    flex-direction: row;

    .role-avatar {
      background-color: @primary-gray;
      margin-right: 0.5rem;
    }

    .chat-bubble {
      background-color: @primary-gray;
    }
  }

  .chat-bubble {
    padding: 0.8rem 1rem;
    border-radius: 1rem;
    word-break: break-word;
    max-width: 80%;
  }
}

.role-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  flex-shrink: 0;
}
</style>