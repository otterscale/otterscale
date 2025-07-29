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
    // Navigation
    breadcrumb: Writable<BreadcrumbState>;

    // Edition
    edition: Writable<string>;

    // Scope & Essential
    triggerUpdateScopes: Writable<boolean>;
    loadingScopes: Writable<boolean>;
    activeScope: Writable<Scope>;
    currentCeph: Writable<Essential | undefined>;
    currentKubernetes: Writable<Essential | undefined>;
}

// Create stores
const createStores = (): AppStores => ({
    breadcrumb: writable<BreadcrumbState>({ parent: "/", current: "/" }),
    edition: writable(m.basic_edition()),
    triggerUpdateScopes: writable(false),
    loadingScopes: writable(true),
    activeScope: writable<Scope>(),
    currentCeph: writable<Essential | undefined>(undefined),
    currentKubernetes: writable<Essential | undefined>(undefined),
});

// Export individual stores
export const {
    breadcrumb,
    edition,
    triggerUpdateScopes,
    loadingScopes,
    activeScope,
    currentCeph,
    currentKubernetes,
} = createStores();