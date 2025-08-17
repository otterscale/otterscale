import { writable, type Writable } from "svelte/store";

interface nfsSnapshotStore {
    selectedScopeUuid: Writable<string| undefined>;
    selectedFacilityName: Writable<string| undefined>;
    selectedVolumeName: Writable<string| undefined>;
    selectedSubvolumeGroupName: Writable<string| undefined>;
    selectedSubvolumeName: Writable<string| undefined>;
}

const createNFSSnapshotStore = (): nfsSnapshotStore => ({
    selectedScopeUuid: writable<string | undefined>(undefined),
    selectedFacilityName: writable<string | undefined>(undefined),
    selectedVolumeName: writable<string | undefined>(undefined),
    selectedSubvolumeGroupName: writable<string | undefined>(undefined),
    selectedSubvolumeName: writable<string | undefined>(undefined),
});

export { createNFSSnapshotStore };
export type { nfsSnapshotStore };
