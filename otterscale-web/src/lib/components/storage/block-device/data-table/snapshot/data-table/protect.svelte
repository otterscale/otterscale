<script lang="ts" module>
	import type {
		Image,
		Image_Snapshot,
		ProtectImageSnapshotRequest
	} from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import type { Writable } from 'svelte/store';
</script>

<script lang="ts">
	let {
		selectedScopeUuid,
		selectedFacility,
		image,
		snapshot,
		data = $bindable()
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
		image: Image;
		snapshot: Image_Snapshot;
		data: Writable<Image[]>;
	} = $props();

	const DEFAULT_REQUEST = {
		scopeUuid: selectedScopeUuid,
		facilityName: selectedFacility,
		imageName: image.name,
		poolName: image.poolName,
		snapshotName: snapshot.name
	} as ProtectImageSnapshotRequest;

	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);
</script>

<button
	class="flex h-full w-full items-center gap-2"
	onclick={() => {
		console.log(request);
		storageClient
			.protectImageSnapshot(request)
			.then((r) => {
				toast.success(`Protect ${request.snapshotName}`);
				storageClient
					.listImages({ scopeUuid: selectedScopeUuid, facilityName: selectedFacility })
					.then((r) => {
						data.set(r.images);
					});
			})
			.catch((e) => {
				toast.error(`Fail to unprotect snapshot: ${e.toString()}`);
			})
			.finally(() => {
				reset();
			});
	}}
>
	<Icon icon="ph:lock-open" />
	Protect
</button>
