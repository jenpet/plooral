<template>
  <div
    class="mt-4"
  >
    <div
      id="organization-add-header"
      @click="visible = !visible"
    >
      <h3>Add a new organization</h3>
      <BIconCaretDown
          :class="{ 'd-none': visible }"
          class="mt-2"
      />
      <BIconCaretUp
        :class="{ 'd-none': !visible }"
        class="mt-2"
      />
    </div>
    <div
      :class="{ 'd-none': !visible }"
    >
      <TextInput
          v-model="name"
          name="name"
          label="Organization Name"
          placeholder="Enter an organization name"
          hint="Provide a name for your organization"
      />
      <SlugInput
        v-model="slug"
        :input="name"
      />
      <TextInput
          v-model="description"
          name="description"
          label="Description"
          placeholder="Enter an organization description"
          hint="Provide a short description for your organization"
      />
      <CheckboxInput
        v-model="protect"
        name="protected"
        label="make organization password protected"
      />
      <CheckboxInput
          v-model="hidden"
          name="hidden"
          label="publicly hide organization (requires password protection)"
      />
      <PasswordInput
        v-model:password="password"
        v-model:valid="passwordValid"
        :visible="protect"
      />
      <div
        class="mt-4"
        style="display: flex; justify-content: space-between"
      >
        <button
          type="button"
          class="
            btn
            btn-primary
          "
          @click="handleSubmit"
        >
          Create
        </button>
        <button
          type="button"
          class="
            btn
            btn-primary
          "
        >
          Reset
        </button>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent, Ref, ref, watch} from "vue";
import SlugInput from "@/components/forms/SlugInput.vue";
import TextInput from "@/components/forms/TextInput.vue";
import CheckboxInput from "@/components/forms/CheckboxInput.vue";
import PasswordInput from "@/components/forms/PasswordInput.vue";
import {createOrganization} from "@/composables/api";
import { getBus } from '@/bus'
import { BIconCaretDown, BIconCaretUp } from 'bootstrap-icons-vue';

const bus = getBus("organizations-form")

export default defineComponent({
  name: "AddOrganizationForm",
  components: {
    PasswordInput,
    SlugInput,
    TextInput,
    CheckboxInput,
    BIconCaretDown,
    BIconCaretUp
  },
  setup() {
    const name:Ref<string> = ref("")
    const slug:Ref<string> = ref("")
    const description:Ref<string> = ref("")
    const protect:Ref<boolean> = ref(false)
    const hidden:Ref<boolean> = ref(false)
    const password:Ref<string> = ref("")
    const passwordValid:Ref<boolean> = ref(false)

    const visible:Ref<boolean> = ref(false)

    const reset = () => {
      name.value = ""
      slug.value = ""
      description.value = ""
      protect.value = false
      hidden.value = false
      password.value = ""
      passwordValid.value = false
    }

    const handleSubmit = () => {
      createOrganization({
        name: name.value,
        slug: slug.value,
        description: description.value,
        hidden: hidden.value,
        password: password.value,
        password_confirmation: password.value,
      }).then(() => {
        bus.emit("organizations:organization:submitted")
        reset()
      })
    }

    // enabling hidden also enables protected
    watch(hidden, newVal => {
      if (newVal && !protect.value) {
        protect.value = true
        return
      }
    })

    // disabling protected also disables hidden
    watch(protect, newVal => {
      if (!newVal && hidden.value) {
        hidden.value = false
      }
    })



    return {
      visible,
      name,
      slug,
      description,
      protect,
      hidden,
      password,
      passwordValid,
      handleSubmit
    }
  }
})
</script>

<style scoped lang="scss">
  #organization-add-header {
    &:hover {
      background-color: rgba(#007bff, 0.2);
    }
    display: flex;
    justify-content: space-between;
    cursor: pointer;
  }

</style>