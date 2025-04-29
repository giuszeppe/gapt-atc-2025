
import { defineComponent, ref, computed, reactive, watch } from "vue";

import SelectionContainer from "../components/SelectionContainer.vue";
import { AdvancementType, InputType, Mode, Role, Scenario, type SelectionItem } from "@/@types/types";
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
          { title: "block", icon: "trowel-bricks", value: InputType.Block },
          { title: "text", icon: "comments", value: InputType.Text },
          { title: "speech", icon: "microphone", value: InputType.Speech },
        ],
      },
      {
        title: "scenario",
        options: [
          { title: "takeoff", icon: "plane-departure", value: Scenario.Takeoff },
          { title: "enroute", icon: "plane", value: Scenario.Enroute },
          { title: "landing", icon: "plane-arrival", value: Scenario.Landing },
        ],
      },
      {
        title: "role",
        options: [
          { title: "aircraft", icon: "plane", value: Role.Aircraft },
          { title: "tower", icon: "tower-observation", value: Role.Tower },
        ],
      },
      {
        title: "mode",
        options: [
          { title: "singleplayer", icon: "user", value: Mode.Singleplayer },
          { title: "multiplayer", icon: "user-group", value: Mode.Multiplayer },
        ],
      },
      {
        title: "simulation advancement type",
        options: [
          { title: "continuous", icon: "repeat", value: AdvancementType.Continuous },
          { title: "click to step", icon: "forward-step", value: AdvancementType.ClickToStep },
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