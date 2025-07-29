import { writable, type Writable } from "svelte/store";
import type { Essential } from "$lib/api/essential/v1/essential_pb";
import type { Scope } from "$lib/api/scope/v1/scope_pb";
import { m } from "$lib/paraglide/messages";

// Types
interface BreadcrumbState {
    parent: string;
    current: string;
}

interface AppStores {
    // Edition
    edition: Writable<string>;

    // Scope & Essential
    triggerUpdateScopes: Writable<boolean>;
    loadingScopes: Writable<boolean>;
    activeScope: Writable<Scope>;
    currentEssentials: Writable<Essential[]>;

    // Navigation
    breadcrumb: Writable<BreadcrumbState>;
}

// Create stores
const createStores = (): AppStores => ({
    edition: writable(m.basic_edition()),
    triggerUpdateScopes: writable(false),
    loadingScopes: writable(true),
    activeScope: writable<Scope>(),
    currentEssentials: writable<Essential[]>([]),
    breadcrumb: writable<BreadcrumbState>({ parent: "/", current: "/" })
});

// Export individual stores
export const {
    edition,
    triggerUpdateScopes,
    loadingScopes,
    activeScope,
    currentEssentials,
    breadcrumb
} = createStores();