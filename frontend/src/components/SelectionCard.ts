import { defineComponent } from "vue";

export default defineComponent({
  name: 'SelectionCard',
  props: {
    title: {
      type: String,
      required: true,
    },
    icon: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    return {
    };
  },
});