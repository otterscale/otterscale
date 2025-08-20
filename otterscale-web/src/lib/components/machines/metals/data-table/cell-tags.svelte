<script lang="ts">
	import { MachineService, type Machine } from '$lib/api/machine/v1/machine_pb';
	import { TagService } from '$lib/api/tag/v1/tag_pb';
	import * as Loading from '$lib/components/custom/loading';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Select from '$lib/components/ui/select/index.js';
	import { cn } from '$lib/utils';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';

	let {
		machine
	}: {
		machine: Machine;
	} = $props();

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	let tags = $state(machine.tags);
	let tagOptions: string[] = $state([]);
	let isTagsLoading = $state(true);
	let isMounted = $state(false);

	const isChanged = $derived(
		!(machine.tags.length === tags.length && machine.tags.every((tag) => tags.includes(tag)))
	);

	const machineClient = createClient(MachineService, transport);
	const tagClient = createClient(TagService, transport);

	let open = $state(false);
	function close() {
		open = false;
	}

	onMount(async () => {
		try {
			tagClient
				.listTags({})
				.then((response) => {
					tagOptions = response.tags.flatMap((tag) => tag.name);
				})
				.finally(() => {
					isTagsLoading = false;
				});

			isMounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if isTagsLoading}
	<Loading.Selection />
{:else}
	<div class="flex w-full justify-end">
		<Select.Root bind:open type="multiple" bind:value={tags}>
			<Select.Trigger class="ring-none m-0 flex-row-reverse border-none p-0 shadow-none">
				{machine.tags.length}
			</Select.Trigger>
			<Select.Content>
				<Select.Group>
					{#each tagOptions as option}
						<Select.Item value={option} label={option}>
							<Icon icon="ph:tag" />{option}
						</Select.Item>
					{/each}
					<div class={cn('grid grid-cols-2 gap-1 border', isChanged ? 'visible' : 'hidden')}>
						<Button
							size="sm"
							onclick={() => {
								toast.promise(
									() =>
										machineClient
											.addMachineTags({
												id: machine.id,
												tags: tags.filter((tag) => !machine.tags.includes(tag))
											})
											.then(() => {
												machineClient.removeMachineTags({
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
							Save
						</Button>
						<Button
							size="sm"
							onclick={() => {
								tags = machine.tags;
							}}
						>
							Reset
						</Button>
					</div>
				</Select.Group>
			</Select.Content>
		</Select.Root>
	</div>
{/if}
