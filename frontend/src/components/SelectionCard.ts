import { AdvancementType, InputType, Mode, Role, Scenario } from "@/@types/types";
import { defineComponent, onMounted, type PropType } from "vue";

export default defineComponent({
  name: "SelectionCard",
  props: {
    title: { type: String, required: true },
    icon: { type: String, required: true },
    value: { type: [String, Object] as PropType<Mode | AdvancementType | InputType | Scenario | Role>, required: true },
    isSelected: { type: Boolean, default: false },
  },
  emits: ["select"],
  setup(props, { emit }) {
    function select() {
      emit("select", props.title);
    }

    return {
      select
    }
  },
});
