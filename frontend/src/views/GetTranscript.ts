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

    const transcript = ref<any>("")

    onMounted(async () => {
      console.log(props.id)
      const response = await axios.get(`http://localhost:8080/get-transcripts/${props.id}`)
      transcript.value = response.data.data
      console.log(transcript.value)
    })

    return {
      transcript,
    }
  }
})