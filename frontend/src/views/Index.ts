
import { defineComponent, ref, computed, reactive, watch, onMounted } from "vue";
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
          { title: "block", icon: "trowel-bricks" },
          { title: "text", icon: "comments" },
          { title: "speech", icon: "microphone" },
        ],
      },
      {
        title: "scenario",
        options: [
          { title: "takeoff", icon: "plane-departure" },
          { title: "enroute", icon: "plane" },
          { title: "landing", icon: "plane-arrival" },
        ],
      },
      {
        title: "simulation",
        options: [],
      },
      {
        title: "role",
        options: [
          { title: "aircraft", icon: "plane" },
          { title: "tower", icon: "tower-observation" },
        ],
      },
      {
        title: "simulation advancement type",
        options: [
          { title: "continuous", icon: "repeat" },
          { title: "click to step", icon: "forward-step" },
        ],
      },
      {
        title: "mode",
        options: [
          { title: "singleplayer", icon: "user" },
          { title: "multiplayer", icon: "user-group" },
        ],
      },
    ]);

    const store = useStore();
    const selections = ref<string[]>([]);
    const currentStep = ref(0);

    async function handleSelection(value: string) {
      selections.value[currentStep.value] = value;
      if (steps[currentStep.value].title === "scenario") {
        await requestScenarios(value)
      }
      if (steps[currentStep.value].title === "mode") {
        setupSimulationMode(value);
      }
      if (currentStep.value < steps.length - 1) {
        currentStep.value++;
      }
    };

    async function requestScenarios(type: string) {
      const response = await axios.get("http://localhost:8080/get-scenarios", {
        params: { type: type },
      });
      const scenariosList = response.data.data.map((item: { Name: string }) => item.Name) as string[];
      const simulationStep = steps.find(step => step.title === "simulation");
      if (simulationStep) {
        simulationStep.options = scenariosList.map((name: string) => ({ title: name, icon: "circle-play" }));
      }
    }

    async function setupSimulationMode(mode: string) {
      const response = await axios.post("http://localhost:8080/post-simulation", { scenario_id: 1, mode: "multi" })
      store.currentSimulation = response.data.data;
      const lobbyCode = response.data.data.lobby_code;
      if (mode === "multiplayer" && lobbyCode) {
        const socket = new WebSocket(`ws://localhost:8080/simulation-lobby?lobby=${lobbyCode}`);
        debugger;
        socket.onopen = () => {
          debugger;
          console.log("WebSocket connection established.");

          const dummyMessage = new Uint8Array([0x01, 0x02, 0x03, 0x04]);
          socket.send(dummyMessage);
        };

        socket.onmessage = (event) => {
          console.log("Message from server:", event.data);
        };

        socket.onerror = (error) => {
          console.error("WebSocket error:", error);
        };

        socket.onclose = () => {
          console.log("WebSocket connection closed.");
        };

      }
    }

    onMounted(() => {
      store.inputType = null;
      store.userRole = null;
    });

    const isComplete = computed(() => selections.value.length === steps.length);

    watch(isComplete, () => {
      const store = useStore();
      store.inputType = selections.value[0] as InputType;
      store.userRole = selections.value[3] as Role;
      // router.push({ name: "simulation" });
    })

    return {
      steps,
      selections,
      currentStep,
      isComplete,
      handleSelection,
    };
  },
});