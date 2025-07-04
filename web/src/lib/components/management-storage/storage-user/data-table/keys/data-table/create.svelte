<script lang="ts" module>
	import type { CreateUserKeyRequest, User } from '$gen/api/storage/v1/storage_pb';
	import { StorageService } from '$gen/api/storage/v1/storage_pb';
	import Button from '$lib/components/ui/button/button.svelte';
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
		user,
		data = $bindable()
	}: {
		selectedScope: string;
		selectedFacility: string;
		user: User;
		data: Writable<User[]>;
	} = $props();

	const DEFAULT_REQUEST = {
		scopeUuid: selectedScope,
		facilityName: selectedFacility,
		userId: user.id
	} as CreateUserKeyRequest;
	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);
</script>

<Button
	class="size-sm flex h-full w-full items-center gap-2"
	onclick={() => {
		console.log(request);
		storageClient
			.createUserKey(request)
			.then((r) => {
				toast.success(`Create ${request.userId}`);
				storageClient
					.listUsers({ scopeUuid: selectedScope, facilityName: selectedFacility })
					.then((r) => {
						data.set(r.users);
					});
			})
			.catch((e) => {
				toast.error(`Fail to create key: ${e.toString()}`);
			})
			.finally(() => {
				reset();
			});
	}}
>
	<Icon icon="ph:plus" />
	Create
</Button>
