<template>
  <div
      class="
      container
      mt-4
    "
  >
    <div
        class="row"
    >
      <div class="col">
        <h1>Boards</h1>
      </div>
    </div>
    <div
        class="row"
    >
      <div class="
        col-12
        col-sm-6
      ">
        <h3>Please pick a board</h3>
        <div
            class="
              list-group
              mt-4
            "
        >
          <a
              v-for="board in boards"
              :key="board.slug"
              :href="`/orgs/${orgSlug}/boards/${board.slug}`"
              class="list-group-item list-group-item-action">
            {{ board.name }}
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent, Ref, ref} from "vue";
import {getBoards, Board} from "@/composables/api";
import {useRoute} from "vue-router";

export default defineComponent({
  name: "BoardList",
  setup() {
    const route = useRoute()

    const orgSlug:Ref<string> = ref(route.params.orgSlug.toString())
    const boards:Ref<Board[]> = ref([])

    getBoards(orgSlug.value).then(val => {
      boards.value = val
    })

    return {
      boards,
      orgSlug
    }
  }
})
</script>

<style lang="scss" scoped>
.container {
  color: #1976D2;
}
</style>