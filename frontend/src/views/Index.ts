
import { defineComponent, ref, computed } from "vue";

import SelectionContainer from "../components/SelectionContainer.vue";
import type { SelectionItem } from "@/@types/types";

export default defineComponent({
  name: "Index",
  components: { SelectionContainer },
  setup() {
    const steps = ref<SelectionItem[]>([
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
    ]);

    const selections = ref<string[]>([]);
    const currentStep = ref(0);

    function handleSelection(value: string) {
      selections.value[currentStep.value] = value;
      if (currentStep.value < steps.value.length - 1) {
        currentStep.value++;
      }
    };

    const isComplete = computed(() => selections.value.length === steps.value.length);

    return {
      steps,
      selections,
      currentStep,
      isComplete,
      handleSelection,
    };
  },
});