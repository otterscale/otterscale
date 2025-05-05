<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import { getContext } from 'svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { Nexus, type CreateScopeRequest } from '$gen/api/nexus/v1/nexus_pb';
	import { toast } from 'svelte-sonner';
	import { Button } from '$lib/components/ui/button';

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const DEFAULT_REQUEST = {} as CreateScopeRequest;

	let createScopeRequest = $state(DEFAULT_REQUEST);

	function reset() {
		createScopeRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class="flex items-center gap-1">
		<Button variant="ghost" class="flex items-center gap-1">
			<Icon icon="ph:plus" /> Scope
		</Button>
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Create Scope</AlertDialog.Title>
			<AlertDialog.Description class="flex items-center justify-between gap-2">
				<Label>Name</Label>
				<Input bind:value={createScopeRequest.name} />
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client
						.createScope(createScopeRequest)
						.then((r) => {
							toast.info(`Create ${r.name} success`);
						})
						.catch((e) => {
							toast.error(`Fail to create ${createScopeRequest.name}: ${e.toString()}`);
						});

					reset();
					close();
				}}
			>
				Create
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
