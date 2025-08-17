import { writable, type Writable } from "svelte/store";

interface GroupStore {
    selectedScopeUuid: Writable<string| undefined>;
    selectedFacilityName: Writable<string| undefined>;
    selectedVolumeName: Writable<string| undefined>;
}

const createGroupStore = (): GroupStore => ({
    selectedScopeUuid: writable<string | undefined>(undefined),
    selectedFacilityName: writable<string | undefined>(undefined),
    selectedVolumeName: writable<string | undefined>(undefined),
});

export { createGroupStore };
export type { GroupStore };
