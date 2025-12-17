<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { User, User_Key } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import CopyButton from '$lib/components/custom/copy-button/copy-button.svelte';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import Button from '$lib/components/ui/button/button.svelte';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		user,
		scope,
		reloadManager
	}: {
		user: User;
		scope: string;
		reloadManager: ReloadManager;
	} = $props();
	const transport: Transport = getContext('transport');

	const storageClient = createClient(StorageService, transport);

	let createdKey = $state<User_Key | undefined>(undefined);
	let resultOpen = $state(false);
</script>

<Button
	class="size-sm flex h-full w-full items-center gap-2 capitalize"
	onclick={() => {
		toast.promise(
			() =>
				storageClient.createUserKey({
					scope: scope,
					userId: user.id
				}),
			{
				loading: `Creating access key...`,
				success: (key) => {
					reloadManager.force();
					createdKey = key;
					resultOpen = true;
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
			}
		);
	}}
>
	<Icon icon="ph:plus" />
	{m.create()}
</Button>

<Modal.Root bind:open={resultOpen}>
	<Modal.Content>
		<Modal.Header>{m.created_access_key()}</Modal.Header>
		{#if createdKey}
			<div class="grid gap-4">
				<div class="group flex w-full items-center gap-2">
					<span
						class="group-hover:bg-text-card rounded-full bg-muted p-2 transition-colors duration-200 group-hover:bg-muted-foreground"
					>
						<Icon icon="ph:key" class="size-5" />
					</span>
					<div class="w-full space-y-1">
						<h3 class="text-xs font-medium text-muted-foreground">{m.access_key()}</h3>
						<span class="flex items-center gap-1">
							<p class="max-w-sm truncate text-sm">{createdKey.accessKey}</p>
							<CopyButton
								class="ml-auto size-4 opacity-0 transition-opacity duration-200 group-hover:opacity-100"
								text={createdKey.accessKey}
							/>
						</span>
					</div>
				</div>

				<div class="group flex w-full items-center gap-2">
					<span
						class="group-hover:bg-text-card rounded-full bg-muted p-2 transition-colors duration-200 group-hover:bg-muted-foreground"
					>
						<Icon icon="ph:lock-key" class="size-5" />
					</span>
					<div class="w-full space-y-1">
						<h3 class="text-xs font-medium text-muted-foreground">{m.secret_key()}</h3>
						<span class="flex items-center gap-1">
							<p class="max-w-sm truncate text-sm">{createdKey.secretKey}</p>
							<CopyButton
								class="ml-auto size-4 opacity-0 transition-opacity duration-200 group-hover:opacity-100"
								text={createdKey.secretKey}
							/>
						</span>
					</div>
				</div>
			</div>
		{/if}
		<Modal.Footer>
			<Modal.Action onclick={() => (resultOpen = false)}>
				{m.close()}
			</Modal.Action>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
