import { type Writable, writable } from 'svelte/store';

interface GroupStore {
	selectedScope: Writable<string | undefined>;
	selectedFacility: Writable<string | undefined>;
	selectedVolumeName: Writable<string | undefined>;
}

const createGroupStore = (): GroupStore => ({
	selectedScope: writable<string | undefined>(undefined),
	selectedFacility: writable<string | undefined>(undefined),
	selectedVolumeName: writable<string | undefined>(undefined)
});

export { createGroupStore };
export type { GroupStore };
