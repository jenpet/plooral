<template>
  <div
    id="board-header"
    class="shadow-sm"
  >
    <h1 class="d-block pixel-style">PLOORAL</h1> <div style="margin-left: 1rem">avoid processes. be productive.</div>
  </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";

export default defineComponent({
  name: "BoardHeader"
})
</script>

<style lang="scss" scoped>
  #board-header {
    min-height: 4rem;
    color: #1976D2;
    display: flex;
    align-items: center;
    padding-left: 2rem;
    font-size: 1.5rem;
    background-color: #F5F5F5;
  }

  .pixel-style {
    font-family: 'Press-Start-2P';
  }
</style>