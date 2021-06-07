import {computed} from "vue";

export function useModelWrapper(props:any, emit:any, name = 'modelValue') {
    return computed({
        get: () => props[name],
        set: (val:any) => emit(`update:${name}`, val)
    })
}