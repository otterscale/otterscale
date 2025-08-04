import { writable, type Writable } from "svelte/store";
import type { Essential } from "$lib/api/essential/v1/essential_pb";
import { PremiumTier } from "$lib/api/premium/v1/premium_pb";
import type { Scope } from "$lib/api/scope/v1/scope_pb";
import { staticPaths, type Path } from "$lib/path";

// Types
interface BreadcrumbState {
    parents: Path[];
    current: Path;
}

interface AppStores {
    // Navigation
    breadcrumb: Writable<BreadcrumbState>;

    // Premium Tier
    premiumTier: Writable<PremiumTier>;

    // Scope & Essential
    activeScope: Writable<Scope>;
    currentCeph: Writable<Essential | undefined>;
    currentKubernetes: Writable<Essential | undefined>;
}

// Create stores
const createStores = (): AppStores => ({
    breadcrumb: writable<BreadcrumbState>({ parents: [], current: staticPaths.home }),
    premiumTier: writable(PremiumTier.BASIC),
    activeScope: writable<Scope>(),
    currentCeph: writable<Essential | undefined>(undefined),
    currentKubernetes: writable<Essential | undefined>(undefined),
});

// Export individual stores
export const {
    breadcrumb,
    premiumTier,
    activeScope,
    currentCeph,
    currentKubernetes,
} = createStores();