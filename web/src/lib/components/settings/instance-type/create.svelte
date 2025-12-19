<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { CreateInstanceTypeRequest } from '$lib/api/instance/v1/instance_pb';
	import { InstanceService } from '$lib/api/instance/v1/instance_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	let { scope, reloadManager }: { scope: string; reloadManager: ReloadManager } = $props();

	const transport: Transport = getContext('transport');
	const virtualMachineClient = createClient(InstanceService, transport);

	let request: CreateInstanceTypeRequest = $state({} as CreateInstanceTypeRequest);
	let invalid: boolean | undefined = $state();
	let open = $state(false);

	function init() {
		request = {
			scope: scope,
			name: '',
			cpuCores: 1,
			memoryBytes: BigInt(1024 ** 3) // 1GB default
		} as CreateInstanceTypeRequest;
	}

	function close() {
		open = false;
	}
</script>

<Modal.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (isOpen) {
			init();
		}
	}}
>
	<Modal.Trigger variant="default">
		<Icon icon="ph:plus" />
		{m.create()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.create_instance_type()}</Modal.Header>
		<Form.Root>
			<!-- ==================== Basic Configuration ==================== -->
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.GeneralRule
						required
						type="text"
						bind:value={request.name}
						bind:invalid
						validateRule="rfc1123"
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.cpu_cores()}</Form.Label>
					<SingleInput.General type="number" bind:value={request.cpuCores} min="1" />
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.memory()}</Form.Label>
					<SingleInput.Measurement
						bind:value={request.memoryBytes}
						transformer={(value) => String(value)}
						units={[{ value: 1024 ** 3, label: 'GB' } as SingleInput.UnitType]}
					/>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>

		<Modal.Footer>
			<Modal.Cancel>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
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
									duration: Number.POSITIVE_INFINITY
								});
								return message;
							}
						});
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
