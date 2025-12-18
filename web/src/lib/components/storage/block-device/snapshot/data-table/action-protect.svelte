<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type {
		Image,
		Image_Snapshot,
		ProtectImageSnapshotRequest
	} from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		snapshot,
		image,
		scope,
		reloadManager,
		closeActions
	}: {
		snapshot: Image_Snapshot;
		image: Image;
		scope: string;
		reloadManager: ReloadManager;
		closeActions: () => void;
	} = $props();

	const transport: Transport = getContext('transport');

	const storageClient = createClient(StorageService, transport);
	const defaults = {
		scope: scope,
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
	class="flex h-full w-full items-center gap-2 capitalize"
	onclick={() => {
		toast.promise(() => storageClient.protectImageSnapshot(request), {
			loading: `Protecting ${request.snapshotName}...`,
			success: () => {
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
			},
			finally: () => {
				closeActions();
			}
		});
		reset();
	}}
>
	<Icon icon="ph:lock-open" />
	{m.protect()}
</button>
