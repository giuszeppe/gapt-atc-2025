<template>
  <div class="guide-wrapper">
    <h1 class="title">User Guide</h1>

    <div class="guide-grid">
      <div class="guide-section">
        <h2>🎯 Your Objective</h2>
        <p>
          Play your assigned role in a simulated conversation. You'll see how the simulation has to be carried out and
          be expected to input the correct ATC phraseology
        </p>
      </div>

      <div class="guide-section">
        <h2>🔄 Turn System</h2>
        <p>
          The simulation is turn-based. When it's your turn, you can use blocks, type or speak your response depending
          on the simulation mode you have selected.
        </p>
      </div>

      <div class="guide-section">
        <h2>👥 single & multi player</h2>
        <p>
          You can use the simulator on your own, and the system will automatically respond, or you can choose to carry
          out the simulation with another person, one playing the role of the aircraft, and the other being the tower.
        </p>
      </div>

      <div class="guide-section">
        <h2>❌ Error Handling</h2>
        <p>
          If your response differs from the expected sentence, incorrect words are <span class="error">highlighted in
            red</span>. You must try again until it's correct to proceed.
        </p>
      </div>

      <div class="guide-section">
        <h2>➡️ How to start</h2>
        <p>
          Go back to the home page, click Start Simulation and follow the prompts to set up your scenario.
          If someone else has created a lobby, click Join Lobby and enter their code to join.
        </p>
      </div>

      <div class="guide-section">
        <h2>📄 Final Transcript</h2>
        <p>
          Once the simulation ends, all your inputs and mistakes are saved and can be reviewed in the corresponding
          transcript.
        </p>
      </div>
    </div>

    <button @click="$router.push('/')" class="option-button mt-10">Go back home</button>

  </div>
</template>

<script setup lang="ts"></script>

<style scoped lang="less">
@import "@/assets/variables.less";

.guide-wrapper {
  text-transform: none;
  padding: 2rem;
  color: white;
  font-family: 'Arial', sans-serif;
  background-color: inherit;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.title {
  text-transform: uppercase;
  font-size: 2.2rem;
  color: @primary-yellow;
  margin-bottom: 2rem;
  text-align: center;
}

.guide-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 2rem;
  max-width: 1100px;
  width: 100%;
  padding: 0 1rem;
}

.guide-section {
  background-color: rgba(0, 0, 0, 0.5);
  border: 1px solid @primary-gray;
  border-radius: 1.5rem;
  padding: 1.5rem;
  min-height: 160px;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
}

.guide-section h2 {
  text-transform: uppercase;
  color: @primary-yellow;
  font-size: 1.2rem;
  margin-bottom: 0.75rem;
}

.guide-section p {
  font-size: 0.95rem;
  line-height: 1.5;
  margin: 0;
}

.error {
  color: red;
  font-weight: bold;
}

.back-button {
  margin-top: 3rem;
  background-color: @primary-yellow;
  color: black;
  padding: 0.75rem 1.5rem;
  text-decoration: none;
  border-radius: 10px;
  font-weight: bold;
  transition: background-color 0.2s;
}

.back-button:hover {
  background-color: darken(@primary-yellow, 10%);
}
</style>
