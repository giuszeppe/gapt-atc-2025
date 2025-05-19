
import { defineComponent, ref, computed, reactive, watch, onMounted, onUnmounted } from "vue";
import SelectionContainer from "../components/SelectionContainer.vue";
import { type SelectionItem, type Role, type InputType, type SimulationItem } from "@/@types/types";
import router from "@/router/router";
import { useStore } from "@/store/store";
import axios from "axios";

export default defineComponent({
  name: "Index",
  components: { SelectionContainer },
  setup() {
    const steps = reactive<SelectionItem[]>([
      {
        title: "input type",
        options: [
          { title: "block", icon: "trowel-bricks", tooltip: "Use a visual block-based editor to construct your input step by step." },
          { title: "text", icon: "comments", tooltip: "Type your instructions or responses as plain text." },
          { title: "speech", icon: "microphone", tooltip: "Speak your input using your device's microphone for voice recognition." },
        ],
      },
      {
        title: "scenario",
        options: [
          { title: "takeoff", icon: "plane-departure", tooltip: "Practice communication and procedures for aircraft takeoff scenarios." },
          { title: "enroute", icon: "plane", tooltip: "Simulate scenarios that occur while the aircraft is enroute between airports." },
          { title: "landing", icon: "plane-arrival", tooltip: "Experience communication and actions required during aircraft landing." },
        ],
      },
      {
        title: "simulation",
        options: [],
      },
      {
        title: "role",
        options: [
          { title: "aircraft", icon: "plane", tooltip: "Take the role of the pilot, responding to ATC instructions and managing the aircraft." },
          { title: "tower", icon: "tower-observation", tooltip: "Act as the tower controller, issuing instructions and managing air traffic." },
        ],
      },
      {
        title: "mode",
        options: [
          { title: "singleplayer", icon: "user", tooltip: "Play alone, the system will act as the counterpart." },
          { title: "multiplayer", icon: "user-group", tooltip: "Create a lobby to play with another user in real time." },
        ],
      },
    ]);

    const store = useStore();
    const selections = ref<string[]>([]);
    const currentStep = ref(0);

    const showSelectionFlow = ref(false);
    const showJoinLobbyInput = ref(false);
    const lobbyCodeInput = ref("");
    const joinLobbyError = ref("");
    const scenariosList = ref<{ id: number, name: string }[]>([]);
    const selectedScenario = ref<number>(null!);

    async function handleSelection(value: string) {
      selections.value[currentStep.value] = value;
      if (steps[currentStep.value].title === "scenario") {
        await requestScenarios(value);
      }
      if (steps[currentStep.value].title === "input type") {
        store.inputType = value as InputType;
      }
      if (steps[currentStep.value].title === "mode") {
        await setupSimulationMode(value);
      }
      if (steps[currentStep.value].title === "role") {
        store.userRole = value as Role;
      }
      if (steps[currentStep.value].title === "simulation") {
        selectedScenario.value = scenariosList.value.find(scenario => scenario.name === value)?.id!;
      }
      if (currentStep.value < steps.length - 1) {
        currentStep.value++;
      }
    }

    async function requestScenarios(type: string) {
      const response = await axios.get("http://localhost:8080/get-scenarios", {
        params: { type },
        headers: {
          "Authorization": store.userToken,
        },
      });
      scenariosList.value = response.data.data
      const simulationStep = steps.find(step => step.title === "simulation");
      if (simulationStep) {
        simulationStep.options = scenariosList.value
          .slice()
          .sort((a, b) => b.name.localeCompare(a.name))
          .map(scenario => ({
            title: scenario.name,
            icon: "circle-play",
            tooltip: `Scenario: ${scenario.name}`,
          }));
      }
    }

    function showUserGuide(){
      router.push({ name: "user-guide" });
    }

    async function setupSimulationMode(mode: string) {
      const response = await axios.post("http://localhost:8080/post-simulation", {
        scenario_id: selectedScenario.value,
        mode,
        role: store.userRole,
        input_type: store.inputType,
      }, {
        headers: {
          "Authorization": store.userToken,
        },
      });
      store.simulationInput = response.data.data.steps[0]
      store.simulationOutline = response.data.data.steps[1]

      store.lobbyCode = response.data.data.lobby_code;
      store.simulationId = response.data.data.simulation.id;
      if (mode === "multiplayer") store.isMultiplayer = true;

      router.push({ name: "simulation" });
    }

    function viewTranscripts() {
      router.push({ name: "transcripts" });
    }

    async function joinLobby() {
      joinLobbyError.value = "";
      if (!lobbyCodeInput.value) return;
      store.isMultiplayer = true;
      store.lobbyCode = lobbyCodeInput.value;
      router.push({ name: "simulation" });
    }

    onMounted(() => {
      store.inputType = null;
      store.userRole = null;
      store.lobbyCode = null;
      store.isMultiplayer = false;
      store.isPlayerInLobby = false;
    });

    return {
      steps,
      selections,
      currentStep,
      showSelectionFlow,
      showJoinLobbyInput,
      lobbyCodeInput,
      joinLobbyError,
      joinLobby,
      handleSelection,
      viewTranscripts,
      showUserGuide,
    };
  }
});