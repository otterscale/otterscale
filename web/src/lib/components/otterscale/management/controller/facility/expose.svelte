<script lang="ts">
	import Icon from '@iconify/svelte';
	import { toast } from 'svelte-sonner';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import {
		FacilityService,
		type Facility,
		type ExposeFacilityRequest
	} from '$gen/api/facility/v1/facility_pb';

	let {
		scopeUuid,
		facilityByCategory
	}: {
		scopeUuid: string;
		facilityByCategory: Facility;
	} = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(FacilityService, transport);

	const DEFAULT_REQUEST = {
		scopeUuid: scopeUuid,
		name: facilityByCategory.name
	} as ExposeFacilityRequest;
	let exposeFacilityRequest = $state(DEFAULT_REQUEST);

	function reset() {
		exposeFacilityRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class="flex items-center gap-1">
		<Icon icon="ph:eye" />
		Expose
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Expose Facility {facilityByCategory.name}</AlertDialog.Title>
			<AlertDialog.Description class="rounded-lg bg-muted/50 p-4">
				Are you sure you want to expose {facilityByCategory.name}?
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client.exposeFacility(exposeFacilityRequest).then((r) => {
						toast.success(`Expose ${exposeFacilityRequest.name}`);
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
