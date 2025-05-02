<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import { writable } from 'svelte/store';
	import { Input } from '$lib/components/ui/input';
	import {
		Nexus,
		type Machine,
		type AddMachineTagsRequest,
		type Tag
	} from '$gen/api/nexus/v1/nexus_pb';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';

	let {
		machine
	}: {
		machine: Machine;
	} = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const tagsStore = writable<Tag[]>();
	const tagsLoading = writable(true);
	async function fetchTags() {
		try {
			const response = await client.listTags({});
			tagsStore.set(response.tags);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			tagsLoading.set(false);
		}
	}

	const DEFAULT_REQUEST = {
		id: machine.id,
		tags: [] as string[]
	} as AddMachineTagsRequest;

	let addMachineTagsRequest = $state(DEFAULT_REQUEST);

	function reset() {
		addMachineTagsRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	let newTags = $state([] as string[]);
	let newTag = $state('');

	onMount(async () => {
		try {
			await fetchTags();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class="flex items-center gap-1">
		<Icon icon="ph:plus" />
		Add
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Add Tags</AlertDialog.Title>
			<AlertDialog.Description class="grid gap-2 rounded-lg  bg-muted/50 p-4">
				<div class="flex flex-wrap gap-2">
					{#each machine.tags as tag}
						<Badge variant="outline">
							{tag}
						</Badge>
					{/each}
					{#each addMachineTagsRequest.tags as tag}
						<Badge>
							{tag}
						</Badge>
					{/each}
				</div>
				<Select.Root type="multiple" bind:value={addMachineTagsRequest.tags}>
					<Select.Trigger>Select</Select.Trigger>
					<Select.Content class="w-fit">
						{#each $tagsStore as tag}
							<Select.Item
								value={tag.name}
								class="flex items-center gap-1"
								disabled={machine.tags.includes(tag.name)}
							>
								{tag.name}
								{#if tag.comment}
									<HoverCard.Root openDelay={13}>
										<HoverCard.Trigger>
											<Icon icon="ph:info" class="size-4 text-blue-800" />
										</HoverCard.Trigger>
										<HoverCard.Content>
											{tag.comment}
										</HoverCard.Content>
									</HoverCard.Root>
								{/if}
							</Select.Item>
						{/each}
						<!-- <Select.Item
							value={''}
							onclick={() => {
								toast.info('new tag');
							}}
						>
							New Tag
						</Select.Item> -->
					</Select.Content>
				</Select.Root>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client
						.addMachineTags(addMachineTagsRequest)
						.then((r) => {
							toast.info(`Add tags ${addMachineTagsRequest.tags.join(', ')}`);
						})
						.catch((e) => {
							toast.error(`Add tags ${addMachineTagsRequest.tags.join(', ')} fail`);
						});
					// toast.info(`Add tags ${addMachineTagsRequest.tags.join(', ')}`);
					console.log(addMachineTagsRequest);
					reset();
					close();
				}}
			>
				Add
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
