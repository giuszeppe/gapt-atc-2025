<template>
  <div class="selection-card" @click="select">
    <div class="icon">
      <font-awesome-icon :icon="`fa-solid fa-${icon}`" />
    </div>
    <span class="title">{{ title }}</span>
  </div>
</template>

<script lang="ts" src="./SelectionCard.ts" />

<style scoped lang="less">
@import "@/assets/variables.less";
.selection-card {
  width: 250px;
  height: 350px;
  border: 3px solid @primary-yellow;
  border-radius: 15px;
  display: flex;
  justify-content: space-evenly;
  align-items: center;
  flex-direction: column;
  cursor: pointer;
  transition: all 0.3s ease-in-out;

  .icon {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 60%;
    width: 100%;
    font-size: 72px;
  }

  &:hover {
    border-color: orange;
  }
}
</style>
