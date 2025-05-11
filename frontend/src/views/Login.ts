import router from "@/router/router";
import { useStore } from "@/store/store";
import axios from "axios";
import { defineComponent, ref } from "vue"


export default defineComponent({
  name: "Login",
  setup() {
    const store = useStore();
    const username = ref<string>("");
    const password = ref<string>("");

    async function handleLogin() {
      const response = await axios.post("http://localhost:8080/login", {
        username: "admin",
        password: "password",
      });
      store.userToken = `Bearer ${response.data.data.token}`;
      router.push({ name: "index" });
    }

    return {
      username,
      password,
      handleLogin,
    }
  }
})