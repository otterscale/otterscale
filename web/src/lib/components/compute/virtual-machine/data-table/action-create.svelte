<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable, type Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import type {
		CreateVirtualMachineRequest,
		VirtualMachineResources,
		VirtualMachineDisk,
	} from '$lib/api/kubevirt/v1/kubevirt_pb';
	import { KubeVirtService } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import * as Code from '$lib/components/custom/code';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Collapsible from '$lib/components/ui/collapsible';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
	import { cn } from '$lib/utils';

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');
	const kubevirtClient = createClient(KubeVirtService, transport);

	let isAdvancedOpen = $state(false);
	let invalidName: boolean | undefined = $state();
	let invalidNamespace: boolean | undefined = $state();
	let invalidResourceCase: boolean | undefined = $state();
	// Add these new state variables for labels
	let labelKey = $state('');
	let labelValue = $state('');

	export const resourcesCase: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: 'instancetypeName',
			label: 'Instance',
			icon: 'ph:copy-simple',
		},
		{
			value: 'custom',
			label: 'Custom',
			icon: 'ph:scales',
		},
	]);
	export const namespaces: Writable<SingleSelect.OptionType[]> = writable([]);

	async function loadNamespaces() {
		try {
			const response = await kubevirtClient.listNamespaces({
				scopeUuid: $currentKubernetes?.scopeUuid,
				facilityName: $currentKubernetes?.name,
			});

			const namespaceOptions = response.namespaces.map((namespace) => ({
				value: namespace,
				label: namespace,
				icon: 'ph:folder',
			}));

			namespaces.set(namespaceOptions);
		} catch (error) {
			toast.error('Failed to load namespaces', {
				description: (error as ConnectError).message.toString(),
			});
		}
	}

	const DEFAULT_REQUEST = {
		scopeUuid: $currentKubernetes?.scopeUuid,
		facilityName: $currentKubernetes?.name,
		name: '',
		namespace: '',
		networkName: '',
		startupScript: '',
		labels: {},
		disks: [] as VirtualMachineDisk[],
		resources: { case: 'instancetypeName', value: '' },
	} as CreateVirtualMachineRequest;
	const DEFAULT_RESOURCES_CUSTOM = {
		case: 'custom',
		value: {
			cpuCores: 1,
			memoryBytes: 1n * 1024n * 1024n * 1024n, // 1GiB
		} as VirtualMachineResources,
	};
	const DEFAULT_RESOURCES_INSTANCE = {
		case: 'instancetypeName',
		value: '',
	};
	let request: CreateVirtualMachineRequest = $state(DEFAULT_REQUEST);
	let resourcesCustom = $state(DEFAULT_RESOURCES_CUSTOM);
	let resourcesInstance = $state(DEFAULT_RESOURCES_INSTANCE);

	function reset() {
		request = DEFAULT_REQUEST;
		resourcesCustom = DEFAULT_RESOURCES_CUSTOM;
		resourcesInstance = DEFAULT_RESOURCES_INSTANCE;
		isAdvancedOpen = false;
		labelKey = '';
		labelValue = '';
	}
	// Add function to handle adding labels
	function addLabel() {
		if (labelKey.trim() && labelValue.trim()) {
			request.labels = { ...request.labels, [labelKey.trim()]: labelValue.trim() };
			labelKey = '';
			labelValue = '';
		}
	}

	// Add function to remove a label
	function removeLabel(key: string) {
		const { [key]: _, ...rest } = request.labels;
		request.labels = rest;
	}
	let open = $state(false);
	function close() {
		open = false;
	}

	onMount(() => {
		loadNamespaces();
	});
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="default">
		<Icon icon="ph:plus" />
		{m.create()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.create_virtual_machine()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General required type="text" bind:value={request.name} bind:invalid={invalidName} />
				</Form.Field>
				<Form.Field>
					<Form.Label>Namespace</Form.Label>
					<SingleSelect.Root
						required
						options={namespaces}
						bind:value={request.namespace}
						bind:invalid={invalidNamespace}
					>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $namespaces as namespace}
											<SingleSelect.Item option={namespace}>
												<Icon
													icon={namespace.icon ? namespace.icon : 'ph:empty'}
													class={cn('size-5', namespace.icon ? 'visibale' : 'invisible')}
												/>
												{namespace.label}
												<SingleSelect.Check option={namespace} />
											</SingleSelect.Item>
										{/each}
									</SingleSelect.Group>
								</SingleSelect.List>
							</SingleSelect.Options>
						</SingleSelect.Content>
					</SingleSelect.Root>
				</Form.Field>
			</Form.Fieldset>
			<Form.Fieldset>
				<Form.Legend>Resources</Form.Legend>
				<Form.Field>
					<Form.Label>Type</Form.Label>
					<SingleSelect.Root
						required
						options={resourcesCase}
						bind:value={request.resources.case}
						bind:invalid={invalidResourceCase}
					>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $resourcesCase as type}
											<SingleSelect.Item option={type}>
												<Icon
													icon={type.icon ? type.icon : 'ph:empty'}
													class={cn('size-5', type.icon ? 'visibale' : 'invisible')}
												/>
												{type.label}
												<SingleSelect.Check option={type} />
											</SingleSelect.Item>
										{/each}
									</SingleSelect.Group>
								</SingleSelect.List>
							</SingleSelect.Options>
						</SingleSelect.Content>
					</SingleSelect.Root>
				</Form.Field>
				{#if request.resources.case === 'custom'}
					<Form.Field>
						<Form.Label>{m.cpu_cores()}</Form.Label>
						<SingleInput.General required type="number" bind:value={resourcesCustom.value.cpuCores} />
					</Form.Field>
					<Form.Field>
						<Form.Label>{m.memory()}</Form.Label>
						<SingleInput.Measurement
							required
							bind:value={resourcesCustom.value.memoryBytes}
							transformer={(value) => String(value)}
							units={[{ value: 1024 * 1024 * 1024, label: 'GB' } as SingleInput.UnitType]}
						/>
					</Form.Field>
				{:else if request.resources.case === 'instancetypeName'}
					<Form.Field>
						<Form.Label>{m.instance_name()}</Form.Label>
						<SingleInput.General required type="number" bind:value={resourcesInstance.value} />
					</Form.Field>
				{/if}
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Disk</Form.Legend>
			</Form.Fieldset>

			<Collapsible.Root bind:open={isAdvancedOpen} class="py-4">
				<div class="flex items-center justify-between gap-2">
					<p class={cn('text-base font-bold', isAdvancedOpen ? 'invisible' : 'visible')}>Advance</p>
					<Collapsible.Trigger class="bg-muted rounded-full p-1 ">
						<Icon
							icon="ph:caret-left"
							class={cn('transition-all duration-300', isAdvancedOpen ? '-rotate-90' : 'rotate-0')}
						/>
					</Collapsible.Trigger>
				</div>
				<Collapsible.Content>
					<Form.Fieldset>
						<Form.Legend>Advance</Form.Legend>

						<Form.Field>
							<Form.Label>Network Name</Form.Label>
							<SingleInput.General type="text" bind:value={request.networkName} />
						</Form.Field>
						<Form.Field>
							<Form.Label>Labels</Form.Label>
							<div class="space-y-2">
								<!-- Add new label inputs -->
								<div class="flex gap-2">
									<SingleInput.General
										type="text"
										placeholder="Key"
										bind:value={labelKey}
										class="flex-1"
									/>
									<SingleInput.General
										type="text"
										placeholder="Value"
										bind:value={labelValue}
										class="flex-1"
									/>
									<Button
										type="button"
										variant="outline"
										size="sm"
										disabled={!labelKey.trim() || !labelValue.trim()}
										onclick={addLabel}
									>
										<Icon icon="ph:plus" class="size-4" />
										Add
									</Button>
								</div>
								<!-- Display existing labels -->
								{#if Object.keys(request.labels).length > 0}
									<div class="space-y-1">
										{#each Object.entries(request.labels) as [key, value]}
											<div
												class="bg-muted flex items-center justify-between rounded-md px-3 py-2"
											>
												<span class="text-sm">
													<span class="font-medium">{key}</span>: {value}
												</span>
												<Button
													type="button"
													variant="ghost"
													size="sm"
													onclick={() => removeLabel(key)}
												>
													<Icon icon="ph:x" class="size-4" />
												</Button>
											</div>
										{/each}
									</div>
								{/if}
							</div>
						</Form.Field>
						<Form.Field>
							<Form.Label>Startup Script</Form.Label>
							<Code.Root lang="bash" class="w-full" hideLines code={request.startupScript}>
								<Code.CopyButton />
							</Code.Root>
							<SingleInput.Structure preview={false} bind:value={request.startupScript} language="bash" />
						</Form.Field>
					</Form.Fieldset>
				</Collapsible.Content>
			</Collapsible.Root>
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
					disabled={invalidName || invalidNamespace}
					onclick={() => {
						toast.promise(() => kubevirtClient.createVirtualMachine(request), {
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
