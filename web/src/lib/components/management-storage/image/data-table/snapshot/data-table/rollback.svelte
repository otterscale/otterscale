<script lang="ts" module>
	import type {
		Image,
		Image_Snapshot,
		RollbackImageSnapshotRequest
	} from '$gen/api/storage/v1/storage_pb';
	import { StorageService } from '$gen/api/storage/v1/storage_pb';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import type { Writable } from 'svelte/store';
</script>

<script lang="ts">
	let {
		selectedScope,
		selectedFacility,
		image,
		snapshot,
		data = $bindable()
	}: {
		selectedScope: string;
		selectedFacility: string;
		image: Image;
		snapshot: Image_Snapshot;
		data: Writable<Image[]>;
	} = $props();

	const DEFAULT_REQUEST = {
		scopeUuid: selectedScope,
		facilityName: selectedFacility,
		imageName: image.name,
		poolName: image.poolName,
		snapshotName: snapshot.name
	} as RollbackImageSnapshotRequest;

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
			.rollbackImageSnapshot(request)
			.then((r) => {
				toast.success(`Rollback ${request.snapshotName}`);
				storageClient
					.listImages({ scopeUuid: selectedScope, facilityName: selectedFacility })
					.then((r) => {
						data.set(r.images);
					});
			})
			.catch((e) => {
				toast.error(`Fail to rollback snapshot: ${e.toString()}`);
			})
			.finally(() => {
				reset();
			});
	}}
>
	<Icon icon="ph:lock-open" />
	Rollback
</button>
