<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type {
		Image,
		Image_Snapshot,
		UnprotectImageSnapshotRequest
	} from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
	import { currentCeph } from '$lib/stores';
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
		scope: $currentCeph?.scope,
		facility: $currentCeph?.name,
		imageName: image.name,
		poolName: image.poolName,
		snapshotName: snapshot.name
	} as UnprotectImageSnapshotRequest;
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}
</script>

<button
	class="flex h-full w-full items-center gap-2 capitalize"
	onclick={() => {
		toast.promise(() => storageClient.unprotectImageSnapshot(request), {
			loading: `Unprotecting ${request.snapshotName}...`,
			success: () => {
				reloadManager.force();
				return `Unprotect ${request.snapshotName}`;
			},
			error: (error) => {
				let message = `Fail to unprotect ${request.snapshotName}`;
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
	{m.unprotect()}
</button>
