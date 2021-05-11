<template>
  <div
    class="
      d-flex
      justify-content-end
    "
  >
    <a
      class="
        d-block
      "
      v-for="button in buttons"
      :key="button.id"
      :href="(typeof button.href !== 'undefined' ? button.href : '#')"
      :title="button.label"
      @click="emitClick(button.id)"
    >
      <component
          :is="button.icon"
      ></component>
    </a>
  </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import { getBus } from '@/bus'
import { BIconGithub, BIconSave2 } from 'bootstrap-icons-vue';

const componentId = 'action-menu'

export default defineComponent({
  name: 'ActionMenu',
  components: {
    BIconGithub,
    BIconSave2
  },
  setup() {
    const bus = getBus(componentId)
    const buttons = [
        {
          id: "export",
          icon: "BIconSave2",
          label: "Export",
          emit: true,
        },
        {
          id: "github",
          icon: "BIconGithub",
          label: "GitHub.com",
          href: "https://github.com/jenpet/plooral",
          emit: false
        }
    ]

    function emitClick(el : string) {
      bus.emit("sidebar:action-menu:" + el)
    }
    return {
      buttons,
      emitClick
    }
  }
})
</script>

<style lang="scss" scoped>
  a {
    font-size: 2rem;
    margin-right: 1rem;
    color: dimgrey;
    &:hover {
      color: #1976D2;
    }
  }
</style>