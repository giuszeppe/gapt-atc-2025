<template>
  <div class="tooltip-wrapper">
    <div class="selection-card" @click="select">
      <div class="icon">
        <font-awesome-icon :icon="`fa-solid fa-${icon}`" />
      </div>
      <span class="title">{{ title }}</span>
    </div>
    <div class="tooltip-text">{{ tooltip }}</div>
  </div>
</template>

<script lang="ts" src="./SelectionCard.ts" />

<style scoped lang="less">
@import "@/assets/variables.less";

.tooltip-wrapper {
  position: relative;
  display: inline-block;

  .tooltip-text {
    position: absolute;
    top: 100%;
    left: 50%;
    transform: translateX(-50%);
    margin-top: 12px; 
    background-color: rgba(0, 0, 0, 0.8);
    color: @primary-yellow;
    padding: 8px 12px;
    border-radius: 6px;
    white-space: nowrap;
    font-size: 14px;
    opacity: 0;
    pointer-events: none;
    transition: opacity 0.2s ease;
    z-index: 1000;

    &::after {
      content: "";
      position: absolute;
      bottom: 100%; 
      left: 50%;
      margin-left: -6px;
      border-width: 6px;
      border-style: solid;
      border-color: transparent transparent rgba(0, 0, 0, 0.8) transparent;
    }
  }

  &:hover .tooltip-text {
    opacity: 1;
  }
}

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
