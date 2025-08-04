import { writable, type Writable } from "svelte/store";
import type { Essential } from "$lib/api/essential/v1/essential_pb";
import { PremiumTier } from "$lib/api/premium/v1/premium_pb";
import type { Scope } from "$lib/api/scope/v1/scope_pb";

// Types
interface BreadcrumbState {
    parents: string[];
    current: string;
}

interface AppStores {
    // Navigation
    breadcrumb: Writable<BreadcrumbState>;

    // Premium Tier
    premiumTier: Writable<PremiumTier>;

    // Scope & Essential
    triggerUpdateScopes: Writable<boolean>;
    activeScope: Writable<Scope>;
    currentCeph: Writable<Essential | undefined>;
    currentKubernetes: Writable<Essential | undefined>;
}

// Create stores
const createStores = (): AppStores => ({
    breadcrumb: writable<BreadcrumbState>({ parents: ["/"], current: "/" }),
    premiumTier: writable(PremiumTier.BASIC),
    triggerUpdateScopes: writable(false),
    activeScope: writable<Scope>(),
    currentCeph: writable<Essential | undefined>(undefined),
    currentKubernetes: writable<Essential | undefined>(undefined),
});

// Export individual stores
export const {
    breadcrumb,
    premiumTier,
    triggerUpdateScopes,
    activeScope,
    currentCeph,
    currentKubernetes,
} = createStores();