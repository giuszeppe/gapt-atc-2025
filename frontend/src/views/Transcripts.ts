import { defineComponent, ref } from "vue";
import SelectionContainer from "@/components/SelectionContainer.vue";
import { useRouter } from "vue-router";
import axios from "axios";
import { useStore } from "@/store/store";
import type { SelectionItem } from "@/@types/types";

export default defineComponent({
  name: "Transcripts",
  components: { SelectionContainer },
  setup() {
    const store = useStore();
    const router = useRouter();

    const currentStep = ref(0);
    const selectedCategory = ref("");
    const selectedSimulation = ref("");

    const categorySelection = {
      title: "category",
      options: [
        { title: "takeoff", icon: "plane-departure", tooltip: "Takeoff scenarios" },
        { title: "enroute", icon: "plane", tooltip: "Enroute scenarios" },
        { title: "landing", icon: "plane-arrival", tooltip: "Landing scenarios" },
      ],
    };

    const simulationSelection = ref<SelectionItem>({
      title: "simulation",
      options: [],
    });

    const transcriptSelection = ref<SelectionItem>({
      title: "transcript",
      options: [],
    });

    const rawTranscriptData = ref<any>(null);

    async function handleCategorySelection(value: string) {
      selectedCategory.value = value;
      currentStep.value = 1;

      const response = await axios.get("http://localhost:8080/get-transcripts", {
        headers: {
          Authorization: store.userToken,
        },
      });
      rawTranscriptData.value = response.data.data[value];

      const simNames = Object.keys(rawTranscriptData.value);
      simulationSelection.value.options = simNames.map(name => ({
        title: name,
        icon: "circle-play",
        tooltip: name,
      }));
    }

    function handleSimulationSelection(value: string) {
      selectedSimulation.value = value;
      currentStep.value = 2;

      const simulationData = rawTranscriptData.value[selectedSimulation.value];

      transcriptSelection.value.options = Object.entries(simulationData).map(([id]) => ({
        title: `Transcript ${id}`,
        icon: "folder",
        tooltip: `View transcript #${id}`,
      }));
    }

    function handleTranscriptSelection(value: string) {
      const transcriptId = value.split(" ")[1];
      router.push({ path: `/transcripts/${transcriptId}` });
    }

    return {
      currentStep,
      categorySelection,
      simulationSelection,
      transcriptSelection,
      handleCategorySelection,
      handleSimulationSelection,
      handleTranscriptSelection,
    };
  },
});