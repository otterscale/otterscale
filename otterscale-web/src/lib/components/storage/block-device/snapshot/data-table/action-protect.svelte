<script lang="ts" module>
	import type {
		Image,
		Image_Snapshot,
		ProtectImageSnapshotRequest
	} from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { currentCeph } from '$lib/stores';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
</script>

<script lang="ts">
	let {
		snapshot
	}: {
		snapshot: Image_Snapshot;
	} = $props();

	const image: Image = getContext('image');
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	const storageClient = createClient(StorageService, transport);
	const defaults = {
		scopeUuid: $currentCeph?.scopeUuid,
		facilityName: $currentCeph?.name,
		imageName: image.name,
		poolName: image.poolName,
		snapshotName: snapshot.name
	} as ProtectImageSnapshotRequest;
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}
</script>

<button
	class="flex h-full w-full items-center gap-2"
	onclick={() => {
		toast.promise(() => storageClient.protectImageSnapshot(request), {
			loading: `Protecting ${request.snapshotName}...`,
			success: (response) => {
				reloadManager.force();
				return `Protect ${request.snapshotName}`;
			},
			error: (error) => {
				let message = `Fail to protect ${request.snapshotName}`;
				toast.error(message, {
					description: (error as ConnectError).message.toString(),
					duration: Number.POSITIVE_INFINITY
				});
				return message;
			}
		});
		reset();
	}}
>
	<Icon icon="ph:lock-open" />
	Protect
</button>
