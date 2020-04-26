import {writable} from "svelte/store";

export const qs = writable([])
export const auth_state = writable("pending")