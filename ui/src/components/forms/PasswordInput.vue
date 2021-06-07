<template>
  <div
    :class="{'d-block': visible, 'd-none': !visible}"
  >
    <TextInput
      v-model="passwordInput"
      name="password"
      label="Password"
      placeholder="Enter a password"
      hint="Enter a password with a certain complexity"
      type="password"
    />
    <TextInput
      v-model="passwordConfirmationInput"
      name="password-confirmation"
      label="Password Confirmation"
      placeholder="Re-enter the password"
      :hint="`Re-enter the same password${hintAddition}`"
      type="password"
    />
  </div>
</template>

<script lang="ts">
import {defineComponent, Ref, ref, watch} from "vue";
import TextInput from "@/components/forms/TextInput.vue";

export default defineComponent({
  name: "PasswordInput",
  components: {
    TextInput
  },
  props: {
    password: String,
    valid: Boolean,
    visible: Boolean
  },
  setup(props, {emit}) {
    const passwordInput:Ref<string> = ref("")
    const passwordConfirmationInput:Ref<string> = ref("")
    const hintAddition:Ref<string> = ref("")

    const handlePasswordChange = () => {
      const valid = passwordConfirmationInput.value !== "" && passwordInput.value === passwordConfirmationInput.value
      emit("update:password", valid ? passwordInput.value : "")
      emit("update:valid", valid)
      if (!valid) {
        hintAddition.value = ". PASSWORDS DO NOT MATCH!"
      } else {
        hintAddition.value = ""
      }
    }

    watch(passwordInput, handlePasswordChange)
    watch(passwordConfirmationInput, handlePasswordChange)

    return {
      passwordInput,
      passwordConfirmationInput,
      hintAddition
    }
  }
})
</script>

<style scoped>

</style>