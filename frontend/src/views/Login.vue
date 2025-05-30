<template>
  <div>
    <div class="login-container">
      <h1>Login</h1>
      <form @submit.prevent="handleLogin">
        <input v-model="username" type="text" placeholder="Username" required />
        <input v-model="password" type="password" placeholder="Password" required />
        <button type="submit" class="option-button">Login</button>
      </form>
    </div>
  </div>
</template>

<script src="./Login.ts" lang="ts"></script>

<style lang="less">
@import "@/assets/variables.less";

.login-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  gap: 20px;

  h1 {
    color: @primary-yellow;
    font-size: 32px;
    margin-bottom: 20px;
  }

  form {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 15px;
    width: 100%;
    max-width: 300px;
    font-size: 16px;
  }

  input {
    width: 100%;
    padding: 10px 15px;
    border: 2px solid @primary-yellow;
    border-radius: 5px;
    background-color: transparent;
    color: @primary-yellow;

    &::placeholder {
      color: fade(@primary-yellow, 60%);
    }
  }
}
</style>