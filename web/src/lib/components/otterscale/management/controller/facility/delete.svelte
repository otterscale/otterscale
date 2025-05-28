<script lang="ts">
	import Icon from '@iconify/svelte';
	import { toast } from 'svelte-sonner';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { Switch } from '$lib/components/ui/switch/index.js';

	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';

	import { Nexus, type Facility, type DeleteFacilityRequest } from '$gen/api/nexus/v1/nexus_pb';

	let {
		scopeUuid,
		facilityByCategory
	}: {
		scopeUuid: string;
		facilityByCategory: Facility;
	} = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(Nexus, transport);

	const DEFAULT_REQUEST = {
		scopeUuid: scopeUuid,
		name: facilityByCategory.name,
		destroyStorage: false,
		force: false
	} as DeleteFacilityRequest;
	let deleteFacilityRequest = $state(DEFAULT_REQUEST);

	function reset() {
		deleteFacilityRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class="flex items-center gap-1">
		<Icon icon="ph:trash" />
		Delete
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Delete Facility {facilityByCategory.name}</AlertDialog.Title>
			<AlertDialog.Description>
				<fieldset class="grid items-center gap-4 rounded-lg bg-muted/50 p-4">
					<Label class="text-sm font-medium">Name</Label>
					<Input bind:value={deleteFacilityRequest.name} />
					<div class="flex justify-between">
						<Label class="text-sm font-medium">Destroy Storage</Label>
						<Switch bind:checked={deleteFacilityRequest.destroyStorage} />
					</div>
					<div class="flex justify-between">
						<Label class="text-sm font-medium">Force</Label>
						<Switch bind:checked={deleteFacilityRequest.force} />
					</div>
				</fieldset>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client.deleteFacility(deleteFacilityRequest).then((r) => {
						toast.success(`Delete ${deleteFacilityRequest.name}`);
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
