<template>
  <div
    class="
      d-flex
      flex-column
    "
  >
    <div
        class="
          element-row
          my-4
        "
        v-for="(type, typeIdx) in elements"
        :key="typeIdx"
    >
      <div
          v-for="(el, elementIdx) in type"
          :key="elementIdx"
          class="element"
          :class="el.class"
          :style="getElementStyle(elementIdx, el)"
          @dragstart="(ev) => {setElementSpecs(ev,el)}"
          draggable="true"
      ></div>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import colors from '@/composables/canvas/shapes/colors'
import {list, ShapeDefintion} from '@/composables/canvas/shapes';

export default defineComponent({
  name: 'ElementSelector',
  setup() {

    const elements = list()

    function setElementSpecs(ev : DragEvent, el: ShapeDefintion) {
      // TODO: Hacky to set it, definitely needs a rework
      el.x = ev.offsetX
      el.y = ev.offsetY
      ev.dataTransfer!.setData("text/plain", JSON.stringify(el))
    }

    function getElementStyle(idx: number, el: ShapeDefintion) {
      return {
        ...el.style,
        zIndex: colors.length-idx,
        marginLeft: idx*2 + 'rem'
      }
    }

    return {
      colors,
      elements,
      getElementStyle,
      setElementSpecs,
    }
  }
});
</script>

<style lang="scss" scoped>
  .element-row {
    position: relative;
    width: 100%;
    height: 6rem;
    div.element {
      position: absolute;
      border: 1px solid;
      cursor: pointer;
      transition: 0.2s;
      &:hover {
        transform: scale(1.1);
      }
    }
  }
</style>