import type { Scope } from "$lib/api/scope/v1/scope_pb";
import { writable } from "svelte/store";

export const scopeLoading = writable<boolean>(true);
export const activeScope = writable<Scope>();
