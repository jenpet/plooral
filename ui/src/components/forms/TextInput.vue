<template>
  <div
    class="
      form-input
      my-2
    "
  >
    <label
        :for="`${name}-input`"
    >
      {{ label }}
    </label>
    <input
        :id="`${name}-input`"
        :type="type"
        class="form-control"
        :aria-describedby="`${name}-input-help`"
        :placeholder="placeholder"
        v-model="value"
    />
    <small
        :id="`${name}-input-help`"
        class="
          form-text
          text-muted
      "
    >
      {{ hint }}
    </small>
  </div>
</template>

<script lang="ts">
import {defineComponent, ref, watch} from "vue";

export default defineComponent({
  name: "TextInput",
  props: {
    modelValue: String,
    name: String,
    label: String,
    placeholder: String,
    type: {
      type: String,
      default: "text"
    },
    hint: String
  },
  setup(props, {emit}) {
    // somehow simply utilizing :modelValue="modelValue" and @update:modeValue="()" does not work properly
    // so this is a workaround with an internal variable of the component
    const value = ref("")
    watch(value, newVal => {
      emit("update:modelValue", newVal)
    })
    return {
      value
    }
  }
})
</script>