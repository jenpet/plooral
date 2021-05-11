<template>
  <div
      id="sidebar"
      class="
        d-flex
        flex-column
      "
  >
    <div
      class="flex-grow-1"
    >
      <h4>{{ title }}</h4>
      <ElementSelector/>
    </div>
    <div
      style="position: fixed; bottom: 1rem; left: 1rem;"
    >
      <ActionMenu/>
    </div>
  </div>
</template>

<script>
import {defineComponent, toRefs} from "vue";
import ElementSelector from "@/components/board/sidebar/ElementSelector";
import ActionMenu from "@/components/board/sidebar/ActionMenu";

export default defineComponent({
  name: "Sidebar",
  components: {
    ElementSelector,
    ActionMenu
  },
  props: {
    title: {
      type: String,
      required: true
    }
  }
})
</script>

<style scoped>
  h4 {
    color: #1976D2;
  }
  #sidebar {
    background-color: #F5F5F5;
    width: 30rem;
    height: 100%;
    padding-left: 1rem;
    padding-right: 1rem;
    padding-top: 3rem;
  }
</style>