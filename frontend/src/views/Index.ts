
import { defineComponent, ref, computed, reactive, watch } from "vue";
import SelectionContainer from "../components/SelectionContainer.vue";
import { type SelectionItem, type Role, type InputType } from "@/@types/types";
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
        title: "role",
        options: [
          { title: "aircraft", icon: "plane" },
          { title: "tower", icon: "tower-observation" },
        ],
      },
      {
        title: "mode",
        options: [
          { title: "singleplayer", icon: "user" },
          { title: "multiplayer", icon: "user-group" },
        ],
      },
      {
        title: "simulation advancement type",
        options: [
          { title: "continuous", icon: "repeat" },
          { title: "click to step", icon: "forward-step" },
        ],
      },
    ]);

    const selections = ref<string[]>([]);
    const currentStep = ref(0);

    function handleSelection(value: string) {
      selections.value[currentStep.value] = value;
      if (currentStep.value < steps.length - 1) {
        currentStep.value++;
      }
    };

    const isComplete = computed(() => selections.value.length === steps.length);

    watch(isComplete, () => {
      const store = useStore();
      store.inputType = selections.value[0] as InputType;
      store.userRole = selections.value[2] as Role;
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