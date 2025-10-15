<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { CreateInstanceTypeRequest } from '$lib/api/instance/v1/instance_pb';
	import { InstanceService } from '$lib/api/instance/v1/instance_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';

	// Context dependencies
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');
	const virtualMachineClient = createClient(InstanceService, transport);

	// ==================== State Variables ====================

	// UI state
	let open = $state(false);

	// Form validation state
	let invalidName: boolean | undefined = $state();

	// ==================== Default Values & Constants ====================
	const DEFAULT_REQUEST = {
		scope: $currentKubernetes?.scope,
		facility: $currentKubernetes?.name,
		name: '',
		namespace: 'default',
		cpuCores: 1,
		memoryBytes: BigInt(1024 ** 3), // 1GB default
	} as CreateInstanceTypeRequest;

	// ==================== Form State ====================
	let request: CreateInstanceTypeRequest = $state({ ...DEFAULT_REQUEST });
	// ==================== Utility Functions ====================
	function reset() {
		request = { ...DEFAULT_REQUEST };
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
		<Modal.Header>{m.create_instance_type()}Create Instance Type</Modal.Header>
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
					<Form.Label>{m.cpu_cores()}</Form.Label>
					<SingleInput.General required type="number" bind:value={request.cpuCores} min="1" />
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.memory()}</Form.Label>
					<SingleInput.Measurement
						required
						bind:value={request.memoryBytes}
						transformer={(value) => String(value)}
						units={[{ value: 1024 ** 3, label: 'GB' } as SingleInput.UnitType]}
					/>
				</Form.Field>
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
						toast.promise(() => virtualMachineClient.createInstanceType(request), {
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
