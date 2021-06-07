<template>
  <div
    class="
      form-input
      my-2
    "
  >
    <label
      for="slug-input"
    >
      Slug
    </label>
    <input
      id="slug-input"
      type="text"
      class="form-control"
      aria-describedby="slug-input-help"
      :value="modelValue"
      disabled
    />
    <small
      id="slug-input-help"
      class="
          form-text
          text-muted
      "
    >
      The slug will be used as a human readable unique identifier.
    </small>
  </div>
</template>

<script lang="ts">
import {defineComponent, toRefs, watch} from "vue";
import slugify from "@/composables/slug";

export default defineComponent({
  name: "SlugInput",
  props: {
    modelValue: String,
    input: String
  },
  setup(props, {emit}) {
    const { input } = toRefs(props)

    watch(input!, newVal => {
      if (typeof newVal !== "undefined") {
        emit("update:modelValue", slugify(newVal))
      }
    })
  }
})
</script>