<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { type Machine, MachineService } from '$lib/api/machine/v1/machine_pb';
	import * as Loading from '$lib/components/custom/loading';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Select from '$lib/components/ui/select/index.js';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let {
		machine,
		reloadManager
	}: {
		machine: Machine;
		reloadManager: ReloadManager;
	} = $props();

	const transport: Transport = getContext('transport');

	let tags = $state(machine.tags);
	let tagOptions: string[] = $state([]);
	let isTagsLoading = $state(false);
	let hasLoadedTags = $state(false);

	const isChanged = $derived(
		!(machine.tags.length === tags.length && machine.tags.every((tag) => tags.includes(tag)))
	);

	const client = createClient(MachineService, transport);

	let open = $state(false);
	function close() {
		open = false;
	}

	async function loadTags() {
		if (hasLoadedTags) return;

		isTagsLoading = true;
		try {
			const response = await client.listTags({});
			tagOptions = response.tags.flatMap((tag) => tag.name);
			hasLoadedTags = true;
		} catch (error) {
			console.error('Error loading tags:', error);
		} finally {
			isTagsLoading = false;
		}
	}
</script>

<div class="flex w-full flex-row-reverse items-center gap-2">
	{machine.tags.length}
	<Select.Root bind:open type="multiple" bind:value={tags} onOpenChange={loadTags}>
		<Select.Trigger class="ring-none m-0 border-none p-0 shadow-none" />
		<Select.Content>
			<Select.Group>
				{#if isTagsLoading}
					<Loading.Selection />
				{:else}
					{#each tagOptions as option}
						<Select.Item value={option} label={option}>
							<Icon icon="ph:tag" />{option}
						</Select.Item>
					{/each}
				{/if}
				<div
					class={cn('grid grid-cols-2 gap-1 border capitalize', isChanged ? 'visible' : 'hidden')}
				>
					<Button
						size="sm"
						onclick={() => {
							toast.promise(
								() =>
									client
										.addMachineTags({
											id: machine.id,
											tags: tags.filter((tag) => !machine.tags.includes(tag))
										})
										.then(() => {
											client.removeMachineTags({
												id: machine.id,
												tags: machine.tags.filter((tag) => !tags.includes(tag))
											});
										}),
								{
									loading: 'Loading...',
									success: () => {
										reloadManager.force();
										return `Update ${machine.fqdn} tags success`;
									},
									error: (error) => {
										let message = `Fail to udpate ${machine.fqdn} tags`;
										toast.error(message, {
											description: (error as ConnectError).message.toString(),
											duration: Number.POSITIVE_INFINITY
										});
										return message;
									}
								}
							);
							close();
						}}
					>
						{m.save()}
					</Button>
					<Button
						size="sm"
						class="capitalize"
						onclick={() => {
							tags = machine.tags;
						}}
					>
						{m.reset()}
					</Button>
				</div>
			</Select.Group>
		</Select.Content>
	</Select.Root>
</div>
