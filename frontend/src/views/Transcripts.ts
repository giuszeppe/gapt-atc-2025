import { useStore } from "@/store/store";
import axios from "axios";
import { defineComponent, onMounted } from "vue";


export default defineComponent({
  name: "Transcripts",
  setup() {
    const store = useStore();

    onMounted(async () => {
      const response = await axios.get("http://localhost:8080/get-transcripts", {
        headers: {
          'Authorization': `${store.userToken}`,
        },
      });

      console.log("Response from server:", response.data);
    })

    return {}

  }

})