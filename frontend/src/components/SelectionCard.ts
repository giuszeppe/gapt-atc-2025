import { defineComponent } from "vue";

export default defineComponent({
  name: "SelectionCard",
  props: {
    title: { type: String, required: true },
    icon: { type: String, required: true },
    isSelected: { type: Boolean, default: false },
  },
  emits: ["select"],
  methods: {
    select() {
      this.$emit("select", this.title);
    },
  },
});
