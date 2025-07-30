import { writable, type Writable } from "svelte/store";
import type { Essential } from "$lib/api/essential/v1/essential_pb";
import type { Scope } from "$lib/api/scope/v1/scope_pb";
import { m } from "$lib/paraglide/messages";

// Types
interface BreadcrumbState {
    parents: string[];
    current: string;
}

interface AppStores {
    // Navigation
    breadcrumb: Writable<BreadcrumbState>;

    // Edition
    edition: Writable<string>;

    // Scope & Essential
    triggerUpdateScopes: Writable<boolean>;
    activeScope: Writable<Scope | undefined>;
    currentCeph: Writable<Essential | undefined>;
    currentKubernetes: Writable<Essential | undefined>;
}

// Create stores
const createStores = (): AppStores => ({
    breadcrumb: writable<BreadcrumbState>({ parents: ["/"], current: "/" }),
    edition: writable(m.basic_edition()),
    triggerUpdateScopes: writable(false),
    activeScope: writable<Scope | undefined>(undefined),
    currentCeph: writable<Essential | undefined>(undefined),
    currentKubernetes: writable<Essential | undefined>(undefined),
});

// Export individual stores
export const {
    breadcrumb,
    edition,
    triggerUpdateScopes,
    activeScope,
    currentCeph,
    currentKubernetes,
} = createStores();