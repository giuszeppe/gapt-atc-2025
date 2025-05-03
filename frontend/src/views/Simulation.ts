import { defineComponent, onMounted, ref } from "vue";
import type { ChatMessage, SimulationStep } from "@/@types/types";
import { useStore } from "@/store/store";

export default defineComponent({
  name: "Simulation",
  setup() {
    const rightPanelSteps = ref<SimulationStep[]>([]);
    const testOutputSteps = ref<SimulationStep[]>([]);
    const leftPanelMessages = ref<ChatMessage[]>([]);
    const playerInput = ref<string>("");

    const store = useStore();
    const userRole = store.userRole;
    const inputType = store.inputType;

    const currentStepIndex = ref<number>(0);
    const isUserTurn = ref<boolean>(true);

    const synonyms = ref<Record<string, string[]>>({});

    onMounted(async () => {
      const scenarioJson = await fetch("/test.json");
      const scenarioData = await scenarioJson.json();
      rightPanelSteps.value = scenarioData.simulations.takeoff[0].steps;
      synonyms.value = scenarioData.synonyms || {};

      const chatJson = await fetch("/test_output.json");
      const chatData = await chatJson.json();

      testOutputSteps.value = chatData.simulations.takeoff[0].steps;

      autoRespond();
    });

    const handlePlayerInput = () => {
      const step = testOutputSteps.value[currentStepIndex.value];
      if (!step || step.role !== userRole) {
        return;
      }

      if (playerInput.value.trim() === "") return;

      const formattedText = formatUserInput(playerInput.value.trim(), step.text);

      leftPanelMessages.value.push({
        role: userRole,
        text: playerInput.value.trim(),
        formattedText,
      });

      playerInput.value = "";
      currentStepIndex.value++;

      autoRespond();
    };

    function autoRespond() {
      while (
        testOutputSteps.value[currentStepIndex.value] &&
        testOutputSteps.value[currentStepIndex.value].role !== userRole
      ) {
        const step = testOutputSteps.value[currentStepIndex.value];
        leftPanelMessages.value.push({
          role: step.role,
          text: step.text,
        });
        currentStepIndex.value++;
      }

      isUserTurn.value = testOutputSteps.value[currentStepIndex.value]?.role == userRole;
    };

    function formatUserInput(userInput: string, expectedInput: string): string {
      const expectedWords = expectedInput.trim().split(/\s+/).map(normalizeWord);
      const synonymsMap = synonyms.value;

      // Split user input, but allow us to look ahead for multi-word synonyms
      const rawWords = userInput.trim().split(/\s+/);
      const normalizedWords: string[] = [];
      const originalWords: string[] = []; // For preserving original casing

      for (let i = 0; i < rawWords.length; i++) {
        const oneWord = normalizeWord(rawWords[i]);
        const twoWord = i + 1 < rawWords.length
          ? normalizeWord(`${rawWords[i]} ${rawWords[i + 1]}`)
          : null;

        const threeWord = i + 2 < rawWords.length
          ? normalizeWord(`${rawWords[i]} ${rawWords[i + 1]} ${rawWords[i + 2]}`)
          : null;

        let matched = false;

        // Try to match three-word phrases
        if (threeWord && matchesAnySynonym(threeWord, synonymsMap)) {
          normalizedWords.push(threeWord);
          originalWords.push(`${rawWords[i]} ${rawWords[i + 1]} ${rawWords[i + 2]}`);
          i += 2;
          continue;
        }

        // Try to match two-word phrases
        if (twoWord && matchesAnySynonym(twoWord, synonymsMap)) {
          normalizedWords.push(twoWord);
          originalWords.push(`${rawWords[i]} ${rawWords[i + 1]}`);
          i += 1;
          continue;
        }

        // Fallback to single word
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

          if (
            userWord === expected ||
            expectedSynonyms.includes(userWord) ||
            levenshteinDistance(userWord, expected) <= 1
          ) {
            formattedText += `${original} `;
            expectedIndex = j + 1;
            matched = true;
            break;
          }
        }

        if (!matched) {
          formattedText += `<span class="wrong-word">${original}</span> `;
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


    return {
      rightPanelSteps,
      leftPanelMessages,
      playerInput,
      isUserTurn,
      userRole,
      inputType,
      handlePlayerInput,
    };
  },
});
