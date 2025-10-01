<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { CreateDataVolumeRequest, DataVolume_Source } from '$lib/api/virtual_machine/v1/virtual_machine_pb';
	import { DataVolume_Source_Type, VirtualMachineService } from '$lib/api/virtual_machine/v1/virtual_machine_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import Button from '$lib/components/ui/button/button.svelte';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';

	// Context dependencies
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');
	const virtualMachineClient = createClient(VirtualMachineService, transport);

	// ==================== State Variables ====================

	// UI state
	let open = $state(false);

	// Form validation state
	let invalidName: boolean | undefined = $state();

	// ==================== Default Values & Constants ====================
	const DEFAULT_REQUEST = {
		scopeUuid: $currentKubernetes?.scopeUuid,
		facilityName: $currentKubernetes?.name,
		name: '',
		namespace: 'default',
		source: { type: DataVolume_Source_Type.HTTP_URL, data: '' } as DataVolume_Source,
		bootImage: true,
		sizeBytes: BigInt(10 * 1024 ** 3),
	} as CreateDataVolumeRequest;

	// ==================== Form State ====================
	let request: CreateDataVolumeRequest = $state(DEFAULT_REQUEST);
	// ==================== Utility Functions ====================
	function reset() {
		request = DEFAULT_REQUEST;
	}
	function close() {
		open = false;
	}

	// ==================== Lifecycle Hooks ====================
	onMount(() => {
		// Initialize form
		reset();
	});
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="default">
		<Icon icon="ph:plus" />
		{m.create()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.create_data_volume()}</Modal.Header>
		<Form.Root>
			<!-- ==================== Basic Configuration ==================== -->
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General required type="text" bind:value={request.name} bind:invalid={invalidName} />
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.namespace()}</Form.Label>
					<SingleInput.General type="text" bind:value={request.namespace} />
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.size()}</Form.Label>
					<SingleInput.Measurement
						required
						bind:value={request.sizeBytes}
						transformer={(value) => String(value)}
						units={[
							{ value: 1024 ** 3, label: 'GB' } as SingleInput.UnitType,
							{ value: 1024 ** 4, label: 'TB' } as SingleInput.UnitType,
						]}
					/>
				</Form.Field>
				{#if request.source}
					<Form.Label>{m.source()}</Form.Label>
					<SingleInput.General
						type="text"
						bind:value={request.source.data}
						placeholder="https://cloud-images.ubuntu.com/xxx/xxx/xxx.img"
					/>
					<div class="flex justify-end gap-2">
						<Button
							variant="outline"
							size="sm"
							href="https://cloud-images.ubuntu.com/"
							target="_blank"
							class="flex items-center gap-1"
						>
							<Icon icon="ph:arrow-square-out" />
							{m.cloud_image()}
						</Button>
					</div>
				{/if}
			</Form.Fieldset>
		</Form.Root>

		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					reset();
				}}
			>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalidName}
					onclick={() => {
						toast.promise(() => virtualMachineClient.createDataVolume(request), {
							loading: `Creating ${request.name}...`,
							success: () => {
								reloadManager.force();
								return `Successfully created ${request.name}`;
							},
							error: (error) => {
								let message = `Failed to create ${request.name}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY,
								});
								return message;
							},
						});
						reset();
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
