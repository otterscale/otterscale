<script lang="ts">
	import Icon from '@iconify/svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Nexus, type Facility, type UpdateFacilityRequest } from '$gen/api/nexus/v1/nexus_pb';

	let {
		scopeUuid,
		facilityByCategory
	}: {
		scopeUuid: string;
		facilityByCategory: Facility;
	} = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const DEFAULT_REQUEST = {
		scopeUuid: scopeUuid,
		name: facilityByCategory.name,
		configYaml: facilityByCategory.metadata?.configYaml
	} as UpdateFacilityRequest;
	let updateFacilityRequest = $state(DEFAULT_REQUEST);

	function reset() {
		updateFacilityRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class="flex items-center gap-1">
		<Icon icon="ph:pencil" />
		Edit
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Update Facility {facilityByCategory.name}</AlertDialog.Title>
			<AlertDialog.Description class="grid items-center gap-4 rounded-lg border p-4">
				<span class="grid items-center gap-2">
					<Label class="text-sm font-medium">Scope</Label>
					<Input bind:value={updateFacilityRequest.scopeUuid} />
				</span>
				<span class="grid items-center gap-2">
					<Label class="text-sm font-medium">name</Label>
					<Input bind:value={updateFacilityRequest.name} />
				</span>
				<span class="grid items-center gap-2">
					<Label class="text-sm font-medium">Configuration</Label>
					<Input bind:value={updateFacilityRequest.configYaml} />
				</span>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client.updateFacility(updateFacilityRequest).then((r) => {
						toast.info(`Update ${updateFacilityRequest.name}`);
					});
					// console.log(updateFacilityRequest);
					toast.info(`Update ${updateFacilityRequest.name}`);
					reset();
					close();
				}}
			>
				Confirm
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
