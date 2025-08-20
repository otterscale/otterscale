<script lang="ts">
	import type { CreateUserKeyRequest, User } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import Button from '$lib/components/ui/button/button.svelte';
	import { currentCeph } from '$lib/stores';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	const user: User = getContext('user');
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	const storageClient = createClient(StorageService, transport);

	const defaults = {
		scopeUuid: $currentCeph?.scopeUuid,
		facilityName: $currentCeph?.name,
		userId: user.id
	} as CreateUserKeyRequest;
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}
</script>

<Button
	class="size-sm flex h-full w-full items-center gap-2"
	onclick={() => {
		toast.promise(() => storageClient.createUserKey(request), {
			loading: `Creating access key...`,
			success: (response) => {
				reloadManager.force();
				return `Create access key`;
			},
			error: (error) => {
				let message = `Fail to create access key`;
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
	<Icon icon="ph:plus" />
	Create
</Button>
