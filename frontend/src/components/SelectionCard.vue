<template>
  <div class="selection-card">
    <div class="icon">
      <font-awesome-icon :icon="`fa-solid fa-${icon}`" />
    </div>
    <span>{{ title }}</span>
  </div>
</template>

<script lang="ts" src="./SelectionCard.ts" />

<style scoped lang="less">
.selection-card {
  width: 250px;
  height: 350px;
  border: 3px solid yellow;
  color: yellow;
  border-radius: 15px;
  display: flex;
  justify-content: space-evenly;
  align-items: center;
  flex-direction: column;

  .icon {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 60%;
    width: 100%;
    font-size: 72px;
  }

  span {
    font-size: 24px;
    text-transform: uppercase;
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  }

  &:hover {}
}
</style>