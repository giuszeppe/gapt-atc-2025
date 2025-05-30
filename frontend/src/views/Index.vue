<template>
  <div class="selection-flow">
    <div v-if="!showSelectionFlow" class="welcome-buttons">
      <template v-if="!showJoinLobbyInput">
        <button class="option-button" @click="showUserGuide">User guide</button>
        <button class="option-button" @click="showSelectionFlow = true">Start Simulation</button>
        <button class="option-button" @click="showJoinLobbyInput = true">Join Lobby</button>
        <button class="option-button" @click="viewTranscripts">View Transcripts</button>
      </template>

      <div v-if="showJoinLobbyInput" class="join-lobby-input">
        <input v-model="lobbyCodeInput" placeholder="Enter lobby code" class="option-button" />
        <button @click="joinLobby" class="option-button">Join</button>
        <button @click="showJoinLobbyInput = !showJoinLobbyInput" class="option-button">Go back home</button>
        <div v-if="joinLobbyError" class="error-text">{{ joinLobbyError }}</div>
      </div>
    </div>

    <div v-if="showSelectionFlow">
      <selection-container v-if="currentStep < steps.length" :selection-item="steps[currentStep]"
        @selectionConfirmed="handleSelection"></selection-container>
    </div>
  </div>

</template>


<script lang="ts" src="./Index.ts">
</script>

<style scoped lang="less">
@import "@/assets/variables.less";

.selection-flow {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
}

.welcome-buttons {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 20px;
  width: 100%;
  height: 100%;
}

.join-lobby-input {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 20px;
  width: 100%;
  height: 100%;
}
</style>