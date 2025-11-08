import { type Writable, writable } from 'svelte/store';

interface NFSStore {
	selectedScope: Writable<string | undefined>;
	selectedFacility: Writable<string | undefined>;
	selectedVolumeName: Writable<string | undefined>;
	selectedSubvolumeGroupName: Writable<string | undefined>;
}

const createNFSStore = (): NFSStore => ({
	selectedScope: writable<string | undefined>(undefined),
	selectedFacility: writable<string | undefined>(undefined),
	selectedVolumeName: writable<string | undefined>(undefined),
	selectedSubvolumeGroupName: writable<string | undefined>(undefined)
});

export { createNFSStore };
export type { NFSStore };
