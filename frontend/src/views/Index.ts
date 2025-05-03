
import { defineComponent, ref, computed, reactive, watch, onMounted } from "vue";
import SelectionContainer from "../components/SelectionContainer.vue";
import { type SelectionItem, type Role, type InputType, type SimulationItem } from "@/@types/types";
import router from "@/router/router";
import { useStore } from "@/store/store";

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

    const selections = ref<string[]>([]);
    const currentStep = ref(0);

    async function handleSelection(value: string) {
      selections.value[currentStep.value] = value;
      if (steps[currentStep.value].title === "scenario") {
        const scenariosList = await loadScenarios(value);
        const simulationStep = steps.find(step => step.title === "simulation");
        if (simulationStep) {
          simulationStep.options = scenariosList.map((name: string) => ({ title: name, icon: "circle-play" }));
        }
      }
      if (currentStep.value < steps.length - 1) {
        currentStep.value++;
      }
    };

    async function loadScenarios(value: string) {
      const scenarioJson = await fetch("/test.json");
      const scenarioData = await scenarioJson.json();
      return scenarioData.simulations[value].map((simulation: SimulationItem) => simulation.name);
    }

    onMounted(() => {
      const store = useStore();
      store.inputType = null;
      store.userRole = null;
    });

    const isComplete = computed(() => selections.value.length === steps.length);

    watch(isComplete, () => {
      const store = useStore();
      store.inputType = selections.value[0] as InputType;
      store.userRole = selections.value[3] as Role;
      router.push({ name: "simulation" });
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