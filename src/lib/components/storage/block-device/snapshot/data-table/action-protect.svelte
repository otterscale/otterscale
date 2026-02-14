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
</script>

<button
	class="flex h-full w-full items-center gap-2 capitalize"
	onclick={() => {
		toast.promise(
			() =>
				storageClient.protectImageSnapshot({
					scope: scope,
					imageName: image.name,
					poolName: image.poolName,
					snapshotName: snapshot.name
				} as ProtectImageSnapshotRequest),
			{
				loading: `Protecting ${snapshot.name}...`,
				success: () => {
					reloadManager.force();
					return `Protect ${snapshot.name}`;
				},
				error: (error) => {
					let message = `Fail to protect ${snapshot.name}`;
					toast.error(message, {
						description: (error as ConnectError).message.toString(),
						duration: Number.POSITIVE_INFINITY,
						closeButton: true
					});
					return message;
				},
				finally: () => {
					closeActions();
				}
			}
		);
	}}
>
	<Icon icon="ph:lock-open" />
	{m.protect()}
</button>
