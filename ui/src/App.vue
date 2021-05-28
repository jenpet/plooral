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
  @font-face {
    font-family: "Patrick-Hand";
    src: local("Patrick-Hand"),
    url(./assets/fonts/PatrickHand-Regular.ttf) format("truetype");
  }

  @font-face {
    font-family: "Press-Start-2P";
    src: local("Press-Start-2P"),
    url(./assets/fonts/PressStart2P-Regular.ttf) format("truetype");
  }

  html, body {
    margin: 0;
    height: 100%;
  }

  #app {
    height: 100%;
  }
</style>