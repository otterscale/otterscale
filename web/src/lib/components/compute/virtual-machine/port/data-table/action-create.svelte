<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { writable, type Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import type { Application_Service_Port } from '$lib/api/application/v1/application_pb';
	import type {
		CreateVirtualMachineServiceRequest,
		VirtualMachine,
	} from '$lib/api/virtual_machine/v1/virtual_machine_pb';
	import { VirtualMachineService } from '$lib/api/virtual_machine/v1/virtual_machine_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import Button from '$lib/components/ui/button/button.svelte';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
	import { cn } from '$lib/utils';

	// Protocol options
	const protocolOptions = writable([
		{ value: 'TCP', label: 'TCP', icon: 'ph:network' },
		{ value: 'UDP', label: 'UDP', icon: 'ph:network' },
		{ value: 'SCTP', label: 'SCTP', icon: 'ph:network' },
	]);

	let {
		virtualMachine,
	}: {
		virtualMachine: VirtualMachine;
	} = $props();

	// Context dependencies
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');
	const virtualMachineClient = createClient(VirtualMachineService, transport);

	// ==================== State Variables ====================

	// UI state
	let open = $state(false);

	// Form validation state
	let invalidServiceName: boolean | undefined = $state();

	// ==================== Default Values & Constants ====================
	const DEFAULT_REQUEST = {
		scopeUuid: $currentKubernetes?.scopeUuid,
		facilityName: $currentKubernetes?.name,
		namespace: virtualMachine.namespace,
		name: '',
		virtualMachineName: virtualMachine.name,
		ports: [] as Application_Service_Port[],
	} as CreateVirtualMachineServiceRequest;

	// ==================== Form State ====================
	let request: CreateVirtualMachineServiceRequest = $state(DEFAULT_REQUEST);

	// New port configuration state
	let newPort = $state({
		port: undefined as number | undefined,
		nodePort: undefined as number | undefined,
		name: '',
		protocol: 'TCP',
		targetPort: '',
	});

	// ==================== Utility Functions ====================
	function reset() {
		request = DEFAULT_REQUEST;
		newPort = {
			port: undefined,
			nodePort: undefined,
			name: '',
			protocol: 'TCP',
			targetPort: '',
		};
	}

	function close() {
		open = false;
	}

	function addPort() {
		if (newPort.port && newPort.port > 0 && newPort.targetPort.trim()) {
			request.ports = [...request.ports, { ...newPort } as Application_Service_Port];
			// Reset newPort to defaults
			newPort = {
				port: undefined,
				nodePort: undefined,
				name: '',
				protocol: 'TCP',
				targetPort: '',
			};
		}
	}

	function removePort(index: number) {
		request.ports = request.ports.filter((_, i) => i !== index);
	}
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="default">
		<Icon icon="ph:plus" />
		{m.create()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.create_service()}</Modal.Header>
		<Form.Root>
			<!-- ==================== Service Information ==================== -->
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.service_name()}</Form.Label>
					<SingleInput.General
						required
						type="text"
						bind:value={request.name}
						bind:invalid={invalidServiceName}
					/>
				</Form.Field>
			</Form.Fieldset>

			<!-- ==================== Port Configuration ==================== -->
			<Form.Fieldset>
				<Form.Legend>{m.ports()}</Form.Legend>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General type="text" bind:value={newPort.name} />
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.protocol()}</Form.Label>
					<SingleSelect.Root options={protocolOptions} required bind:value={newPort.protocol}>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $protocolOptions as option}
											<SingleSelect.Item {option}>
												<Icon
													icon={option.icon ? option.icon : 'ph:empty'}
													class={cn('size-5', option.icon ? 'visible' : 'invisible')}
												/>
												{option.label}
												<SingleSelect.Check {option} />
											</SingleSelect.Item>
										{/each}
									</SingleSelect.Group>
								</SingleSelect.List>
							</SingleSelect.Options>
						</SingleSelect.Content>
					</SingleSelect.Root>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.port()}</Form.Label>
					<SingleInput.General
						type="number"
						required
						bind:value={newPort.port}
						min="1"
						max="32767"
						placeholder="8080"
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.node_port()}</Form.Label>
					<SingleInput.General
						type="number"
						bind:value={newPort.nodePort}
						min="0"
						max="32767"
						placeholder="30080"
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.target_port()}</Form.Label>
					<SingleInput.General type="text" required bind:value={newPort.targetPort} placeholder="8080" />
				</Form.Field>
				<div class="flex justify-end">
					<Button
						type="button"
						variant="outline"
						size="sm"
						disabled={!newPort.port || !newPort.targetPort.trim()}
						onclick={addPort}
					>
						<Icon icon="ph:plus" class="size-4" />
						{m.add()}
					</Button>
				</div>

				<!-- Display Configured Ports -->
				{#if request.ports.length > 0}
					<div class="space-y-2">
						<h4 class="font-medium">Configured Ports</h4>
						{#each request.ports as port, index}
							<div class="bg-muted flex items-center justify-between rounded-md px-3 py-2">
								<div class="flex-1">
									<div class="flex items-center gap-2">
										<Icon icon="ph:network" class="size-4" />
										<span class="font-medium">{port.name || `Port ${index + 1}`}</span>
									</div>
									<div class="text-muted-foreground text-sm">
										{port.port}{#if port.nodePort && port.nodePort > 0}:{port.nodePort}{/if} → {port.targetPort}
										({port.protocol})
									</div>
								</div>
								<Button type="button" variant="ghost" size="sm" onclick={() => removePort(index)}>
									<Icon icon="ph:x" class="size-4" />
								</Button>
							</div>
						{/each}
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
					disabled={invalidServiceName || !request.name || request.ports.length === 0}
					onclick={() => {
						toast.promise(() => virtualMachineClient.createVirtualMachineService(request), {
							loading: `Creating service ${request.name}...`,
							success: () => {
								reloadManager.force();
								return `Successfully created service ${request.name}`;
							},
							error: (error) => {
								let message = `Failed to create service ${request.name}`;
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
