import { defineComponent } from "vue";

export default defineComponent({
  name: "SelectionCard",
  props: {
    title: { type: String, required: true },
    icon: { type: String, required: true },
    tooltip: { type: String, required: true },
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
