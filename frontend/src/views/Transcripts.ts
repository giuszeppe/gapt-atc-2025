import axios from "axios";
import { defineComponent, onMounted } from "vue";


export default defineComponent({
  name: "Transcripts",
  setup() {

    onMounted(async () => {
      const response = await axios.get("http://localhost:8080/get-transcripts");

      console.log("Response from server:", response.data);
    })

    return {}

  }

})