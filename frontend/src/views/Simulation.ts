import { defineComponent, onMounted, onUnmounted, ref, watch } from "vue";
import type { ChatMessage, SimulationStep } from "@/@types/types";
import { useStore } from "@/store/store";
import { useSpeechToText } from "@/composables/useSpeechToText";
import VoiceVisualizer from "@/components/VoiceVisualizer.vue";
import axios from "axios";
import SimulationEndModal from "@/components/SimulationEndModal.vue";

export default defineComponent({
  name: "Simulation",
  components: { VoiceVisualizer, SimulationEndModal },
  setup() {
    const { transcript, isListening, volume, outputBuffer, start, stop, replayAudio } = useSpeechToText();

    const rightPanelSteps = ref<SimulationStep[]>([]);
    const testOutputSteps = ref<SimulationStep[]>([]);
    const leftPanelMessages = ref<ChatMessage[]>([]);
    const playerInput = ref<string>("");
    const stepCount = ref<number>(0);
    const showEndModal = ref<boolean>(false);

    const wordBlocks = ref<string[]>([]);
    const selectedWords = ref<string[]>([]);
    const dragIndex = ref<number | null>(null)

    const store = useStore();
    const userRole = store.userRole;
    const simulationId = store.simulationId;
    const inputType = store.inputType;
    const simulationOutline = store.simulationOutline;
    const simulationInput = store.simulationInput;
    const socket = ref<WebSocket | null>(null);

    const currentStepIndex = ref<number>(0);
    const isUserTurn = ref<boolean>(true);

    const synonyms = ref<Record<string, string[]>>({});

    onMounted(async () => {
      if (store.isMultiplayer && !store.isPlayerInLobby) {
        socket.value = new WebSocket(`ws://localhost:8080/simulation-lobby?lobby=${store.lobbyCode}`);
        store.isPlayerInLobby = true;

        socket.value.onopen = () => {
          console.log("WebSocket connection established.");
          if(socket.value ){
            socket.value.send(store.userToken)
          }
        };

        socket.value.onmessage = async (event) => {
          const text = await event.data.text();
          const jsonData = JSON.parse(text);
          console.log("Message from server:", jsonData);
        };

        socket.value.onerror = (error) => {
          console.error("WebSocket error:", error);
        };

        socket.value.onclose = () => {
          console.log("WebSocket connection closed.");
        };
      }

      if (inputType === "block") {
        const initialStep = testOutputSteps.value[currentStepIndex.value];
        if (initialStep && initialStep.role === userRole) {
          wordBlocks.value = initialStep.text.trim().split(/\s+/);
        }
      }

      const scenarioJson = await fetch("/test.json");
      const scenarioData = await scenarioJson.json();
      rightPanelSteps.value = scenarioData.simulations.takeoff[0].steps;
      synonyms.value = scenarioData.synonyms || {};

      const chatJson = await fetch("/test_output.json");
      const chatData = await chatJson.json();

      testOutputSteps.value = chatData.simulations.takeoff[0].steps;
      stepCount.value = testOutputSteps.value.length;

      autoRespond();
      window.addEventListener("beforeunload", handleBeforeUnload);
    });

    onUnmounted(() => {
      window.removeEventListener("beforeunload", handleBeforeUnload);
    });

    function handleBeforeUnload() {
      store.isPlayerInLobby = false;
    }

    watch(leftPanelMessages.value, async (newVal) => {
      if (newVal.length === stepCount.value) {
        await axios.post(
          "http://localhost:8080/end-simulation", {
          simulation_id: simulationId,
          messages: leftPanelMessages.value,
        }, {
          headers: {
            "Authorization": store.userToken,
          },
        });
        showEndModal.value = true;
      }
    });

    const handlePlayerInput = () => {
      const step = testOutputSteps.value[currentStepIndex.value];
      if (!step || step.role !== userRole) return;

      let inputText = "";

      if (inputType != "block") {
        inputText = playerInput.value.trim();
        if (!inputText) return;
      } else if (inputType === "block") {
        inputText = selectedWords.value.join(" ").trim();
        if (!inputText) return;
      }

      const formattedText = formatUserInput(inputText, step.text);
      const object: ChatMessage = {
        role: userRole,
        type: "text",
        content: formattedText,
      };

      leftPanelMessages.value.push(object);

      if (socket.value) {
        const msg = JSON.stringify(object);
        console.log("Sending message to server:", msg);
        socket.value.send(msg);
      }

      playerInput.value = "";
      selectedWords.value = [];
      wordBlocks.value = [];

      currentStepIndex.value++;

      if (inputType === "block") {
        const nextStep = testOutputSteps.value[currentStepIndex.value];
        if (nextStep && nextStep.role === userRole) {
          wordBlocks.value = nextStep.text.trim().split(/\s+/);
        } else {
          wordBlocks.value = [];
        }
      }
      autoRespond();
    };

    function autoRespond() {
      while (
        testOutputSteps.value[currentStepIndex.value] &&
        testOutputSteps.value[currentStepIndex.value].role !== userRole
      ) {
        const step = testOutputSteps.value[currentStepIndex.value];
        const content = formatUserInput(step.text, step.text);
        leftPanelMessages.value.push({
          role: step.role,
          content: content,
          type: "text",
        });
        currentStepIndex.value++;
      }

      isUserTurn.value = testOutputSteps.value[currentStepIndex.value]?.role === userRole;

      if (isUserTurn.value && inputType === "block") {
        const step = testOutputSteps.value[currentStepIndex.value];
        if (step) {
          wordBlocks.value = shuffleArray(step.text.trim().split(/\s+/));
        }
      }
    }

    function uint8ToBase64(bytes: Uint8Array): string {
      let binary = '';
      const len = bytes.byteLength;
      for (let i = 0; i < len; i++) {
        binary += String.fromCharCode(bytes[i]);
      }
      return btoa(binary);
    }

    async function base64ToAudioBuffer(base64: string): Promise<AudioBuffer> {
      const binary = atob(base64);

      const len = binary.length;
      const bytes = new Uint8Array(len);
      for (let i = 0; i < len; i++) {
        bytes[i] = binary.charCodeAt(i);
      }

      const floatBuffer = new Float32Array(bytes.buffer);
      const audioContext = new AudioContext();

      const audioBuffer = audioContext.createBuffer(1, floatBuffer.length, audioContext.sampleRate);
      audioBuffer.copyToChannel(floatBuffer, 0);

      return audioBuffer;
    }


    function toggleListening() {
      if (isListening.value) {
        stop(async () => {
          playerInput.value = transcript.value.trim();
          if (!outputBuffer.value) return
          if (socket.value && outputBuffer.value) {
            const channelData = outputBuffer.value.getChannelData(0);

            const float32Array = new Float32Array(channelData.length);
            float32Array.set(channelData);
            const uint8Array = new Uint8Array(float32Array.buffer);

            const base64 = uint8ToBase64(uint8Array)
            console.log("BASE&$", base64);
            socket.value.send(JSON.stringify({ type: "audio", content: base64 }));

            replayAudio((await base64ToAudioBuffer(base64)));
          }

        });
      } else {
        start();
      }
    }

    // #region TEXT
    function formatUserInput(userInput: string, expectedInput: string): string {
      const expectedWords = expectedInput.trim().split(/\s+/).map(normalizeWord);
      const synonymsMap = synonyms.value;

      const rawWords = userInput.trim().split(/\s+/);
      const normalizedWords: string[] = [];
      const originalWords: string[] = [];

      for (let i = 0; i < rawWords.length; i++) {
        const oneWord = normalizeWord(rawWords[i]);
        const twoWord =
          i + 1 < rawWords.length ? normalizeWord(`${rawWords[i]} ${rawWords[i + 1]}`) : null;
        const threeWord =
          i + 2 < rawWords.length ? normalizeWord(`${rawWords[i]} ${rawWords[i + 1]} ${rawWords[i + 2]}`) : null;

        if (threeWord && matchesAnySynonym(threeWord, synonymsMap)) {
          normalizedWords.push(threeWord);
          originalWords.push(`${rawWords[i]} ${rawWords[i + 1]} ${rawWords[i + 2]}`);
          i += 2;
          continue;
        }

        if (twoWord && matchesAnySynonym(twoWord, synonymsMap)) {
          normalizedWords.push(twoWord);
          originalWords.push(`${rawWords[i]} ${rawWords[i + 1]}`);
          i += 1;
          continue;
        }

        normalizedWords.push(oneWord);
        originalWords.push(rawWords[i]);
      }

      let formattedText = "";
      let expectedIndex = 0;

      for (let i = 0; i < normalizedWords.length; i++) {
        const userWord = normalizedWords[i];
        const original = originalWords[i];
        let matched = false;

        for (let j = expectedIndex; j < Math.min(expectedWords.length, expectedIndex + 3); j++) {
          const expected = expectedWords[j];
          const expectedSynonyms = synonymsMap[expected] || [];

          if (userWord === expected || expectedSynonyms.includes(userWord) || levenshteinDistance(userWord, expected) <= 1) {
            formattedText += `${original} `;
            expectedIndex = j + 1;
            matched = true;
            break;
          }
        }

        if (!matched) {
          formattedText += `<span class='wrong-word'>${original}</span> `;
        }
      }

      return formattedText.trim();
    }

    function matchesAnySynonym(phrase: string, synonymsMap: Record<string, string[]>) {
      for (const key in synonymsMap) {
        if (synonymsMap[key].includes(phrase)) return true;
      }
      return false;
    }

    function normalizeWord(word: string): string {
      return word.replace(/[.,]/g, "").toLowerCase();
    }

    function levenshteinDistance(a: string, b: string): number {
      const matrix: number[][] = [];

      for (let i = 0; i <= b.length; i++) {
        matrix[i] = [i];
      }
      for (let j = 0; j <= a.length; j++) {
        matrix[0][j] = j;
      }

      for (let i = 1; i <= b.length; i++) {
        for (let j = 1; j <= a.length; j++) {
          if (b.charAt(i - 1) === a.charAt(j - 1)) {
            matrix[i][j] = matrix[i - 1][j - 1];
          } else {
            matrix[i][j] = Math.min(
              matrix[i - 1][j - 1] + 1,
              matrix[i][j - 1] + 1,
              matrix[i - 1][j] + 1
            );
          }
        }
      }
      return matrix[b.length][a.length];
    }

    // #endregion

    // #region BLOCK
    function onDragStart(index: number) {
      dragIndex.value = index
    }

    function onDrop(targetIndex: number) {
      if (dragIndex.value === null || dragIndex.value === targetIndex) return

      const draggedWord = selectedWords.value[dragIndex.value]
      selectedWords.value.splice(dragIndex.value, 1)
      selectedWords.value.splice(targetIndex, 0, draggedWord)

      dragIndex.value = null
    }

    function selectWord(word: string) {
      if (!selectedWords.value.includes(word)) {
        selectedWords.value.push(word)
      }
    }

    function deselectWord(word: string) {
      const index = selectedWords.value.indexOf(word)
      if (index !== -1) {
        selectedWords.value.splice(index, 1)
      }
    }

    function shuffleArray<T>(array: T[]): T[] {
      const shuffled = array.slice();
      for (let i = shuffled.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [shuffled[i], shuffled[j]] = [shuffled[j], shuffled[i]];
      }
      return shuffled;
    }

    const selectedWordIndexes = ref<number[]>([]);

    function toggleWordSelection(index: number) {
      const idx = selectedWordIndexes.value.indexOf(index);
      if (idx !== -1) {
        selectedWordIndexes.value.splice(idx, 1);
      } else {
        selectedWordIndexes.value.push(index);
      }
    }

    // #endregion

    return {
      rightPanelSteps,
      leftPanelMessages,
      playerInput,
      isUserTurn,
      userRole,
      inputType,
      isListening,
      transcript,
      volume,
      simulationId,
      wordBlocks,
      selectedWords,
      showEndModal,
      selectedWordIndexes,
      toggleWordSelection,
      toggleListening,
      handlePlayerInput,
      selectWord,
      deselectWord,
      onDragStart,
      onDrop,
    };
  },
});
