import { useStore } from "@/store/store"
import axios from "axios"
import { defineComponent, onMounted, ref } from "vue"

export default defineComponent({
  name: "GetTranscript",
  props: {
    id: {
      type: Number,
      required: true,
    },
  },
  setup(props) {
    const store = useStore()
    const transcript = ref<any>("")
    const simulationOutline = ref<[]>([])

    onMounted(async () => {
      console.log(props.id)
      const response = await axios.get(`http://localhost:8080/get-transcripts/${props.id}`, {
        headers: {
          "Authorization": store.userToken,
        },
      })
      console.log(response.data.data)
      transcript.value = response.data.data.transcripts
      simulationOutline.value = response.data.data.steps
    })

    return {
      transcript,
    }
  }
})