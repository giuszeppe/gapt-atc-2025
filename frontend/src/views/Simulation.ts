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
    const { transcript, isListening, volume, start, stop } = useSpeechToText();

    const rightPanelSteps = ref<SimulationStep[]>([]);
    const testOutputSteps = ref<SimulationStep[]>([]);
    const leftPanelMessages = ref<ChatMessage[]>([]);
    const playerInput = ref<string>("");
    const stepCount = ref<number>(0);
    const showEndModal = ref<boolean>(false);

    const store = useStore();
    const userRole = store.userRole;
    const inputType = "speech"; // store.inputType;
    const simulationOutline = store.simulationOutline;
    const simulationInput = store.simulationInput;
    const isMultiplayer = store.isMultiplayer;
    const lobbyCode = store.lobbyCode;
    const isPlayerInLobby = store.isPlayerInLobby;
    const socket = ref<WebSocket | null>(null);

    const currentStepIndex = ref<number>(0);
    const isUserTurn = ref<boolean>(true);

    const synonyms = ref<Record<string, string[]>>({});

    onMounted(async () => {
      if (isMultiplayer && !isPlayerInLobby) {
        socket.value = new WebSocket(`ws://localhost:8080/simulation-lobby?lobby=${store.lobbyCode}`);
        store.isPlayerInLobby = true;
        socket.value.onopen = () => {
          console.log("WebSocket connection established.");
        };

        socket.value.onmessage = async (event) => {
          const text = await event.data.text()
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
      if (newVal.length == stepCount.value) {
        await axios.post("http://localhost:8080/end-simulation", { simulation_id: 1, messages: leftPanelMessages.value });
        showEndModal.value = true;
      }
    })

    const handlePlayerInput = () => {
      const step = testOutputSteps.value[currentStepIndex.value];
      if (!step || step.role !== userRole) {
        return;
      }

      if (playerInput.value.trim() === "") return;

      const formattedText = formatUserInput(playerInput.value.trim(), step.text);
      const object: ChatMessage = {
        role: userRole,
        type: "text",
        content: formattedText,
      }

      leftPanelMessages.value.push(object);

      if (socket.value) {
        const a = JSON.stringify(object);
        console.log("Sending message to server:", a);
        socket.value.send(a)
      }

      playerInput.value = "";
      currentStepIndex.value++;

      autoRespond();
    };

    function autoRespond() {
      while (testOutputSteps.value[currentStepIndex.value] && testOutputSteps.value[currentStepIndex.value].role !== userRole) {
        const step = testOutputSteps.value[currentStepIndex.value];
        const content = formatUserInput(step.text, step.text);
        leftPanelMessages.value.push({
          role: step.role,
          content: content,
          type: 'text',
        });
        currentStepIndex.value++;
      }

      isUserTurn.value = testOutputSteps.value[currentStepIndex.value]?.role == userRole;
    };

    function formatUserInput(userInput: string, expectedInput: string): string {
      const expectedWords = expectedInput.trim().split(/\s+/).map(normalizeWord);
      const synonymsMap = synonyms.value;

      const rawWords = userInput.trim().split(/\s+/);
      const normalizedWords: string[] = [];
      const originalWords: string[] = [];

      for (let i = 0; i < rawWords.length; i++) {
        const oneWord = normalizeWord(rawWords[i]);
        const twoWord = i + 1 < rawWords.length ? normalizeWord(`${rawWords[i]} ${rawWords[i + 1]}`) : null;

        const threeWord = i + 2 < rawWords.length ? normalizeWord(`${rawWords[i]} ${rawWords[i + 1]} ${rawWords[i + 2]}`) : null;

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

      let formattedText = '';
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

    function toggleListening() {
      if (isListening.value) {
        stop(() => {
          playerInput.value = transcript.value.trim();
        });
      } else {
        start();
      }
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
      showEndModal,
      toggleListening,
      handlePlayerInput,
    };
  },
});
