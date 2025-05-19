<template>
  <div class="selection-container">
    <span class="title">{{ 'select ' + selectionItem.title }}</span>

    <div class="selections">
      <selection-card v-for="option in selectionItem.options" :key="option.title" :title="option.title" :tooltip="option.tooltip"
        :icon="option.icon" :isSelected="selected === option.title" @select="handleSelection"></selection-card>
    </div>
  </div>
</template>

<script lang="ts" src="./SelectionContainer.ts" />

<style scoped lang="less">
.selection-container {
  text-align: center;
}

.selections {
  display: flex;
  gap: 15px;
  justify-content: center;
  margin-top: 20px;
}
</style>