import { writable, type Writable } from "svelte/store";

interface NFSStore {
    selectedScopeUuid: Writable<string| undefined>;
    selectedFacilityName: Writable<string| undefined>;
    selectedVolumeName: Writable<string| undefined>;
    selectedSubvolumeGroupName: Writable<string| undefined>;
}

const createNFSStore = (): NFSStore => ({
    selectedScopeUuid: writable<string | undefined>(undefined),
    selectedFacilityName: writable<string | undefined>(undefined),
    selectedVolumeName: writable<string | undefined>(undefined),
    selectedSubvolumeGroupName: writable<string | undefined>(undefined),
});

export { createNFSStore };
export type { NFSStore };

