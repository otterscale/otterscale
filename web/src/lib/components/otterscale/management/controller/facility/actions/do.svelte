<script lang="ts">
	import Icon from '@iconify/svelte';
	import { toast } from 'svelte-sonner';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import { Button } from '$lib/components/ui/button';

	import { Nexus, type Facility_Action, type DoActionRequest } from '$gen/api/nexus/v1/nexus_pb';

	let { action }: { action: Facility_Action } = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const DEFAULT_REQUEST = {} as DoActionRequest;
	let doActionRequest = $state(DEFAULT_REQUEST);

	function reset() {
		doActionRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger>
		<Button variant="outline" class="size-6" size="icon">
			<Icon icon="ph:caret-double-right" />
		</Button>
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Conduct Action {action.name}</AlertDialog.Title>
			<AlertDialog.Description>
				<p class="mb-4 text-sm text-muted-foreground">
					{#if action.description}
						{action.description}
					{:else}
						No description available for this action.
					{/if}
				</p>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client.doAction(doActionRequest).then((r) => {
						toast.success(`Conduct ${action.name}`);
					});
					reset();
					close();
				}}
			>
				Confirm
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
