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
		// const userRole = store.userRole;
		const userRole = "aircraft"; // testing
		const inputType = store.inputType;

		const currentStepIndex = ref<number>(0);
		const isUserTurn = ref<boolean>(true);

		onMounted(async () => {

			const scenarioJson = await fetch("/test.json");
			const scenarioData = await scenarioJson.json();
			rightPanelSteps.value = scenarioData.simulations[0].steps;

			const chatJson = await fetch("/test_output.json");
			const chatData = await chatJson.json();
			testOutputSteps.value = chatData.simulations[0].steps;

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
			const userWords = userInput.trim().split(/\s+/);
			const expectedWords = expectedInput.trim().split(/\s+/);

			let formattedText = "";
			let expectedIndex = 0;

			for (let i = 0; i < userWords.length; i++) {
				const userWord = userWords[i];

				if (expectedIndex >= expectedWords.length) {
					formattedText += `<span class="wrong-word">${userWord}</span> `;
					continue;
				}

				const expectedWord = expectedWords[expectedIndex];

				if (normalizeWord(userWord) === normalizeWord(expectedWord)) {
					formattedText += `${userWord} `;
					expectedIndex++;
				} else if (levenshteinDistance(normalizeWord(userWord), normalizeWord(expectedWord)) <= 2) {
					formattedText += `${userWord} `;
					expectedIndex++;
				} else {
					let found = false;
					for (let lookahead = 1; lookahead <= 3; lookahead++) {
						const nextExpected = expectedWords[expectedIndex + lookahead];
						if (nextExpected && normalizeWord(userWord) === normalizeWord(nextExpected)) {
							expectedIndex += lookahead + 1;
							formattedText += `${userWord} `;
							found = true;
							break;
						}
					}

					if (!found) {
						formattedText += `<span class="wrong-word">${userWord}</span> `;
					}
				}
			}

			return formattedText.trim();
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
