import { defineComponent, ref, type PropType } from "vue";
import SelectionCard from "./SelectionCard.vue";
import type { SelectionItem } from "@/@types/types";

export default defineComponent({
  name: "SelectionContainer",
  components: { SelectionCard },
  props: {
    selectionItem: {
      type: Object as PropType<SelectionItem>,
      required: true
    }
  },
  emits: ["selectionConfirmed"],
  setup(_, { emit }) {
    const selected = ref<string>(null!);

    function handleSelection(value: string) {
      selected.value = value;
      emit("selectionConfirmed", value);
    };

    return {
      selected,
      handleSelection
    };
  },
});
