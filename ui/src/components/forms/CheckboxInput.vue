<template>
  <div
      class="
        from-check
        my-2
      "
  >
    <input
      class="
        form-check-input
        me-2
      "
      type="checkbox"
      v-model="value"
      :id="`${name}-checkbox`"
    />
    <label
      class="form-check-label"
      :for="`${name}-checkbox`"
    >
      {{ label }}
    </label>
  </div>
</template>

<script lang="ts">
import {defineComponent, ref, watch} from "vue";
import {useModelWrapper} from "@/composables/modelWrapper.js";

export default defineComponent({
  name: "CheckboxInput",
  props: {
    modelValue: Boolean,
    name: String,
    label: String
  },
  setup(props, {emit}) {
    return {
      value: useModelWrapper(props, emit, "modelValue")
    }
  }
})
</script>