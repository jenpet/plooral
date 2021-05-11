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
        <h1>Organizations</h1>
      </div>
    </div>
    <div
      class="row"
    >
      <div class="
        col-12
        col-sm-6
      ">
        <h3>Please select your organization</h3>
        <div
            class="
              list-group
              mt-4
            "
        >
          <a
              v-for="org in organizations"
              :key="org.slug"
              :href="`/orgs/${org.slug}/boards`"
              class="list-group-item list-group-item-action">
            {{ org.name }}
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent, Ref, ref} from "vue";
import {getOrganizations, Organization} from "@/composables/api";

export default defineComponent({
  name: "OrganizationList",
  setup() {
    const organizations:Ref<Organization[]> = ref([])

    getOrganizations().then(val => {
      organizations.value = val
    })
    return {
      organizations
    }
  }
})
</script>

<style lang="scss" scoped>
  .container {
    color: #1976D2;
  }
</style>