<template>
    <div
      class="
        d-flex
        flex-grow-1
      "
    >
      <Sidebar
        :title="sidebarTitle"
      />
      <div
          id="canvas"
          ref="canvas"
          class="flex-grow-1"
      >
      </div>
      <div
          class="board-information"
      >
        {{ boardInfo }}
      </div>
    </div>
</template>

<script lang="ts">
import {defineComponent, ref, onMounted, Ref, onBeforeUnmount} from 'vue';
import { setup, stageJSON } from '@/composables/canvas'
import Sidebar from '@/components/board/Sidebar.vue'
import {initBus, getBus, BusEvent} from "@/bus";
import {useRoute} from "vue-router";
import {Board, getBoard, storeBoard} from "@/composables/api";

const bus = getBus('board')

export default defineComponent({
  name: 'Board',
  components: {
    Sidebar
  },
  setup() {
    let autosaver : ReturnType<typeof setTimeout>
    const canvas = ref<HTMLUnknownElement>()
    const board = ref<Board>()

    const sidebarTitle = ref<string>('')

    const route = useRoute()
    const orgSlug:Ref<string> = ref(route.params.orgSlug.toString())
    const boardSlug:Ref<string> = ref(route.params.boardSlug.toString())
    const boardInfo = ref<string>()

    const saveBoard = () => {
      const content = stageJSON()
      storeBoard(orgSlug.value, boardSlug.value, content)
          .then(() => {
            bus.emit('board:autosaved')
            autosaver = setTimeout(saveBoard, 10000)
          })
    }

    onBeforeUnmount(() => {
      clearInterval(autosaver)
    })

    onMounted(() => {
        getBoard(orgSlug.value, boardSlug.value).then(val => {
          if (canvas.value !== undefined) {
            board.value = val
            initBus(orgSlug.value, boardSlug.value)
            sidebarTitle.value = `${orgSlug.value} > ${board.value.name}`
            setup('canvas', canvas.value.clientWidth, canvas.value.clientHeight, val.content)
            autosaver = setTimeout(saveBoard, 10000)
          }
        })
    })

    bus.on('board:autosaved', (ev:BusEvent) => {
      if (bus.isOriginTo(ev)) {
        boardInfo.value = `autosaved @${ev.timestamp}`
      }
    })

    return {
      canvas,
      board,
      orgSlug,
      sidebarTitle,
      boardInfo
    }
  },

});
</script>

<style lang="scss">
  .board-information {
    position: fixed;
    bottom: 1rem;
    right: 1rem;
    font-size: 0.75rem;
    &.fadeout {
      visibility: hidden;
      opacity: 0;
      transition: visibility 0s 2s, opacity 2s linear;
    }
  }
  #canvas {
    outline: none;
    background-color: #FAFAFA;
    overflow: scroll;
    width: 100%;
    height: 100%;

    [contenteditable] {
      outline: 0 solid transparent;
      cursor: text;
    }

    .centered-text-edit-container {
      display:flex;
      align-items: center;
      justify-content: center;
      position: absolute;
      padding: 0;
      margin: 0;
      overflow: hidden;
      background: none;
      outline: none;
      &.text {
        text-align: center;
        word-break: break-all;
      }
    }
  }
</style>
