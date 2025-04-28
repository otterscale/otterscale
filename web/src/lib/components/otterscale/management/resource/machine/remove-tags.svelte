<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import { Nexus, type Machine, type RemoveMachineTagsRequest } from '$gen/api/nexus/v1/nexus_pb';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	let {
		machine
	}: {
		machine: Machine;
	} = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const DEFAULT_REQUEST = {
		id: machine.id,
		tags: [] as string[]
	} as RemoveMachineTagsRequest;

	let removeMachineTagsRequest = $state(DEFAULT_REQUEST);

	function reset() {
		removeMachineTagsRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class="flex items-center gap-1">
		<Icon icon="ph:minus" />
		Remove
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Remove Tags</AlertDialog.Title>
			<AlertDialog.Description class="grid gap-2 rounded-lg  bg-muted/50 p-4">
				<div class="flex flex-wrap gap-2">
					{#each machine.tags as tag}
						<Badge
							variant={removeMachineTagsRequest.tags.includes(tag) ? 'destructive' : 'outline'}
						>
							{tag}
						</Badge>
					{/each}
				</div>
				<Select.Root type="multiple" bind:value={removeMachineTagsRequest.tags}>
					<Select.Trigger>Select</Select.Trigger>
					<Select.Content class="w-fit">
						{#each machine.tags as tag}
							<Select.Item value={tag}>
								{tag}
							</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client.removeMachineTags(removeMachineTagsRequest).then((r) => {
						toast.info(`Remove tags ${removeMachineTagsRequest.tags.join(', ')}`);
					});
					// console.log(removeMachineTagsRequest);
					toast.info(`Remove tags ${removeMachineTagsRequest.tags.join(', ')}`);
					reset();
					close();
				}}
			>
				Remove
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
