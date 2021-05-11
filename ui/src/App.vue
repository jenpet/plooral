<template>
  <div
    class="
      d-flex
      h-100
      flex-column
    "
  >
    <Header/>
    <router-view
      class="flex-grow-1"
    />
  </div>
</template>

<script type="ts">
import {defineComponent} from "vue";
import Header from '@/components/Header.vue'

export default defineComponent({
  name: 'Main',
  components: {
    Header
  },
});
</script>

<style lang="scss">
  html, body {
    margin: 0;
    height: 100%;
  }

  #app {
    height: 100%;
  }
</style>