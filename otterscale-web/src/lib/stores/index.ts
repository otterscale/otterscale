import { writable, type Writable } from "svelte/store";
import type { Scope } from "$lib/api/scope/v1/scope_pb";

// Types
interface BreadcrumbState {
    parent: string;
    current: string;
}

// Scope stores
export const scopeLoading: Writable<boolean> = writable(true);
export const activeScope: Writable<Scope> = writable();

// Navigation stores
export const breadcrumb: Writable<BreadcrumbState> = writable({
    parent: "/",
    current: "/"
});