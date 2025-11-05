<script lang="ts">
	import { create } from '@bufbuild/protobuf';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import type { Plan } from './plans';

	import { page } from '$app/state';
	import {
		MachineService,
		type Machine,
		CreateMachineRequestSchema,
		AddMachineTagsRequestSchema,
		CommissionMachineRequestSchema,
		GetMachineRequestSchema,
	} from '$lib/api/machine/v1/machine_pb';
	import { OrchestratorService, CreateNodeRequestSchema } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import { ScopeService, CreateScopeRequestSchema } from '$lib/api/scope/v1/scope_pb';
	import { IPv4AddressInput } from '$lib/components/custom/ipv4';
	import { IPv4CIDRInput } from '$lib/components/custom/ipv4-cidr';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import { Separator } from '$lib/components/ui/separator';
	import * as Sheet from '$lib/components/ui/sheet';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';

	let { open = $bindable(false), plan = $bindable({} as Plan) }: { open: boolean; plan: Plan } = $props();

	// Constants
	const TOAST_DURATION_MS = 1200000; // 20 minutes
	const POLLING_INTERVAL_MS = 5000; // 5 seconds
	const MAX_POLL_ATTEMPTS = 60 * 4; // 20 minutes with 5 second intervals
	const DEVICE_PATH_PREFIX = '/dev/';
	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);
	const scopeClient = createClient(ScopeService, transport);
	const orchestratorClient = createClient(OrchestratorService, transport);
	const machinesStore = writable<Machine[]>([]);

	// Form state
	let scopeName = $state('');
	let selectedMachine = $state('');
	let selectedDevices = $state<string[]>([]);
	let calicoCidr = $state('');
	let virtualIp = $state('');
	let prefixName = $state('');
	let isSubmitting = $state(false);

	async function fetchMachines() {
		try {
			const response = await machineClient.listMachines({});
			machinesStore.set(response.machines);
		} catch (error) {
			console.error('Error fetching machines:', error);
			toast.error('Failed to fetch machines');
		}
	}

	function handleSubmit(event: Event) {
		event.preventDefault();

		if (isSubmitting) return;
		isSubmitting = true;

		// Create a loading toast that we can update
		const toastId = toast.loading(`Creating scope ${scopeName}...`, { duration: TOAST_DURATION_MS });
		handleClose();
		(async () => {
			try {
				// Step 1: Create Scope
				const createScopeRequest = create(CreateScopeRequestSchema, {
					name: scopeName,
				});
				const scopeResponse = await scopeClient.createScope(createScopeRequest);

				// Step 2: Add Machine Tags
				toast.loading('Adding machine tags...', { id: toastId, duration: TOAST_DURATION_MS });
				const addMachineTagsRequest = create(AddMachineTagsRequestSchema, {
					id: selectedMachine,
					tags: [],
				});
				await machineClient.addMachineTags(addMachineTagsRequest);

				// Step 3: Commission Machine
				toast.loading('Commissioning machine...', { id: toastId, duration: TOAST_DURATION_MS });
				const commissionMachineRequest = create(CommissionMachineRequestSchema, {
					id: selectedMachine,
					enableSsh: true,
					skipBmcConfig: false,
					skipNetworking: false,
					skipStorage: false,
				});
				await machineClient.commissionMachine(commissionMachineRequest);

				// Step 4: Wait for machine status to be Ready
				toast.loading('Waiting for machine to be ready...', { id: toastId, duration: TOAST_DURATION_MS });
				const getMachineRequest = create(GetMachineRequestSchema, {
					id: selectedMachine,
				});

				// Poll until machine status is Ready
				let machineReady = false;
				let retryCount = 0;
				while (!machineReady && retryCount < MAX_POLL_ATTEMPTS) {
					await new Promise((resolve) => setTimeout(resolve, POLLING_INTERVAL_MS)); // Wait 5 seconds between checks
					const machineResponse = await machineClient.getMachine(getMachineRequest);
					if (machineResponse.status.toLowerCase() === 'ready') {
						machineReady = true;
					}
					retryCount++;
				}

				// Step 5: Create Machine
				toast.loading('Creating machine...', { id: toastId, duration: TOAST_DURATION_MS });
				const createMachineRequest = create(CreateMachineRequestSchema, {
					id: selectedMachine,
					scope: scopeResponse.name,
				});
				const machineResponse = await machineClient.createMachine(createMachineRequest);

				// Step 6: Wait for agent status to be Started
				toast.loading('Waiting for agent to start...', { id: toastId, duration: TOAST_DURATION_MS });

				// Poll until agent status is Started
				let agentStarted = false;
				retryCount = 0;
				while (!agentStarted && retryCount < MAX_POLL_ATTEMPTS) {
					await new Promise((resolve) => setTimeout(resolve, POLLING_INTERVAL_MS)); // Wait 5 seconds between checks
					const machineResponse = await machineClient.getMachine(getMachineRequest);
					if (machineResponse.agentStatus.toLowerCase() === 'started') {
						agentStarted = true;
					}
					retryCount++;
				}

				// Step 7: Create Node
				toast.loading('Creating node...', { id: toastId, duration: TOAST_DURATION_MS });
				const createNodeRequest = create(CreateNodeRequestSchema, {
					scope: scopeResponse.name,
					machineId: machineResponse.id,
					prefixName: prefixName || scopeName,
					virtualIps: virtualIp ? [virtualIp] : [],
					calicoCidr: calicoCidr,
					osdDevices: selectedDevices.map((device) => `${DEVICE_PATH_PREFIX}${device}`),
				});
				await orchestratorClient.createNode(createNodeRequest);

				// Success
				toast.success(m.create_scope_success({ name: scopeResponse.name }), { id: toastId, duration: 5000 });
				isSubmitting = false;
				resetForm();
			} catch (error) {
				isSubmitting = false;
				const message = `Failed to create scope: ${scopeName}`;
				toast.error(message, {
					id: toastId,
					description: (error as ConnectError).message.toString(),
					duration: Number.POSITIVE_INFINITY,
				});
			}
		})();
	}

	function handleClose() {
		open = false;
	}

	function resetForm() {
		scopeName = '';
		selectedMachine = '';
		selectedDevices = [];
		calicoCidr = '';
		virtualIp = '';
		prefixName = '';
		isSubmitting = false;
	}

	onMount(async () => {
		await fetchMachines();
	});
</script>

<Sheet.Root bind:open onOpenChange={handleClose}>
	<Sheet.Content class="inset-y-auto bottom-0 h-9/10 rounded-tl-lg sm:max-w-4/5">
		<Sheet.Header class="h-full p-0">
			<div class="flex h-full flex-col p-12 lg:max-w-3/5">
				<!-- Plan Header -->
				<div class="flex flex-col space-y-4">
					<Badge variant="secondary" class="bg-primary/10 text-primary flex items-center uppercase">
						{#if plan.star}
							<Icon icon="ph:star-fill" class="text-yellow-500" />
						{/if}
						<span>{plan.tier}</span>
					</Badge>

					<h2 class="text-3xl font-semibold tracking-tight">{plan.name}</h2>
					<p class="text-accent-foreground/80 text-md">{plan.description}</p>

					<div class="flex flex-wrap gap-2">
						{#each plan.tags as tag}
							<Badge variant="outline" class="bg-background/50 backdrop-blur-sm">
								{tag}
							</Badge>
						{/each}
					</div>
					<Separator class="my-2" />
				</div>

				<!-- Form -->
				<form class="flex h-full flex-col justify-between pt-4" onsubmit={handleSubmit}>
					<div class="grid gap-6">
						<!-- Scope Name -->
						<div class="grid gap-4">
							<div class="grid gap-1">
								<Label for="name">
									{m.create_scope_name()}
									<div class="h-2 w-2 rounded-full bg-yellow-500"></div>
								</Label>
								<p class="text-muted-foreground text-sm">
									{m.create_scope_name_description()}
								</p>
							</div>
							<Input id="name" type="text" placeholder="scope-name" bind:value={scopeName} required />
						</div>

						<!-- Machine Selection -->
						<div class="grid gap-4">
							<div class="grid gap-1">
								<Label for="machine">
									{m.create_scope_machine()}
									<div class="h-2 w-2 rounded-full bg-yellow-500"></div>
								</Label>
								<p class="text-muted-foreground text-sm">
									{m.create_scope_machine_description()}
								</p>
							</div>
							<Select.Root type="single" bind:value={selectedMachine} required>
								<Select.Trigger id="storage-devices" class="w-full text-left">
									{#if selectedMachine}
										{@const machine = $machinesStore.find((m) => m.id === selectedMachine)!}
										{@render machineSelectItem(machine)}
									{:else}
										<span>{m.create_scope_machine_select()}</span>
									{/if}
								</Select.Trigger>
								<Select.Content>
									{#each $machinesStore.filter((m) => m.status === 'Ready') as machine}
										<Select.Item value={machine.id}>
											{@render machineSelectItem(machine)}
										</Select.Item>
									{/each}
								</Select.Content>
							</Select.Root>
						</div>

						<!-- Storage Devices -->
						{#if selectedMachine}
							{@const selectedMachineData = $machinesStore.find((m) => m.id === selectedMachine)}
							{@const availableDevices =
								selectedMachineData?.blockDevices?.filter((device) => !device.bootDisk) ?? []}

							<div class="grid gap-4">
								<div class="grid gap-1">
									<Label for="storage-devices">
										{m.create_scope_storage_devices()}
										<div class="h-2 w-2 rounded-full bg-yellow-500"></div>
									</Label>
									<p class="text-muted-foreground text-sm">
										{m.create_scope_storage_devices_description()}
									</p>
								</div>
								<div class="flex gap-4 overflow-x-auto">
									{#each availableDevices as device}
										<Tooltip.Provider>
											<Tooltip.Root>
												<Tooltip.Trigger>
													<Label
														class="hover:bg-accent/50 flex items-start gap-x-2 rounded-md border p-2 has-aria-checked:border-slate-600 has-aria-checked:bg-blue-50 dark:has-aria-checked:border-slate-900 dark:has-aria-checked:bg-slate-950"
													>
														<Checkbox
															id={device.name}
															checked={selectedDevices.includes(device.name)}
															onCheckedChange={(checked) => {
																if (checked) {
																	selectedDevices = [...selectedDevices, device.name];
																} else {
																	selectedDevices = selectedDevices.filter(
																		(d) => d !== device.name,
																	);
																}
															}}
															class="data-[state=checked]:border-slate-600 data-[state=checked]:bg-slate-600 data-[state=checked]:text-white dark:data-[state=checked]:border-slate-700 dark:data-[state=checked]:bg-slate-700"
														/>
														<p class="text-sm leading-none">{device.name}</p>
													</Label>
												</Tooltip.Trigger>
												<Tooltip.Content>
													<p>
														[{device.firmwareVersion}] {device.model}
														{formatCapacity(device.storageMb).value}
														{formatCapacity(device.storageMb).unit}
													</p>
												</Tooltip.Content>
											</Tooltip.Root>
										</Tooltip.Provider>
									{/each}
								</div>
							</div>

							<!-- Network Configuration -->
							<div class="grid grid-cols-2 gap-4">
								<div class="grid gap-4">
									<div class="grid gap-1">
										<Label for="calico-cidr">{m.create_scope_calico_cidr()} ({m.optional()})</Label>
										<p class="text-muted-foreground text-sm">
											{m.create_scope_calico_cidr_description()}
										</p>
									</div>
									<IPv4CIDRInput
										class="font-sans text-sm font-normal"
										placeholder="192.168.0.0/16"
										bind:value={calicoCidr}
									/>
								</div>
								<div class="grid gap-4">
									<div class="grid gap-1">
										<Label for="virtual-ip">{m.create_scope_virtual_ip()} ({m.optional()})</Label>
										<p class="text-muted-foreground text-sm">
											{m.create_scope_virtual_ip_description()}
										</p>
									</div>
									<IPv4AddressInput
										class="font-sans text-sm font-normal"
										placeholder="192.168.1.1"
										bind:value={virtualIp}
									/>
								</div>
							</div>

							<!-- Prefix Name (optional) -->
							<div class="grid gap-4">
								<div class="grid gap-1">
									<Label for="prefix-name">Prefix Name ({m.optional()})</Label>
									<p class="text-muted-foreground text-sm">
										Prefix for node naming (defaults to scope name if not provided)
									</p>
								</div>
								<Input id="prefix-name" type="text" placeholder="node-prefix" bind:value={prefixName} />
							</div>
						{/if}
					</div>

					<!-- Form Actions -->
					<div class="flex gap-8">
						<Button
							size="lg"
							variant="outline"
							class="flex-1"
							onclick={() => {
								handleClose();
								resetForm();
							}}
							disabled={isSubmitting}
						>
							{m.cancel()}
						</Button>
						<Button type="submit" size="lg" class="flex-1" disabled={isSubmitting}>
							{#if isSubmitting}
								<Icon icon="ph:spinner" class="mr-2 size-4 animate-spin" />
								Creating...
							{:else}
								{m.create()}
							{/if}
						</Button>
					</div>
				</form>
			</div>

			<!-- Plan Image -->
			<div class="relative lg:absolute lg:inset-y-0 lg:right-0 lg:w-2/5">
				<img src={plan.image} alt={plan.name} class="absolute inset-0 size-full object-cover" />
			</div>
		</Sheet.Header>
	</Sheet.Content>
</Sheet.Root>

{#snippet machineSelectItem(machine: Machine)}
	<div class="flex items-center space-x-2">
		<HoverCard.Root>
			<HoverCard.Trigger
				href={dynamicPaths.machinesMetal(page.params.scope).url + '/' + machine.id}
				target="_blank"
				rel="noreferrer noopener"
				class="flex items-center space-x-1 rounded-sm underline-offset-4 group-hover:underline focus-visible:outline-2 focus-visible:outline-offset-8 focus-visible:outline-black"
			>
				<Icon icon="ph:info" class="size-4" />
			</HoverCard.Trigger>
			<HoverCard.Content class="w-180">
				<div class="grid grid-cols-6 gap-x-6 gap-y-2">
					<!-- General Information -->
					<h4 class="text-md col-span-6 font-semibold">{m.general()}</h4>
					<span class="text-muted-foreground text-sm">{m.architecture()}</span>
					<span class="text-sm">{machine.architecture}</span>
					<span class="text-muted-foreground text-sm">{m.cpu()}</span>
					<span class="col-span-3 text-sm">{machine.hardwareInformation.cpu_model}</span>
					<span class="text-muted-foreground text-sm">{m.memory()}</span>
					<span class="text-sm">
						{formatCapacity(machine.memoryMb).value}
						{formatCapacity(machine.memoryMb).unit}
					</span>
					<span class="text-muted-foreground text-sm">{m.storage()}</span>
					<span class="text-sm">
						{formatCapacity(machine.storageMb).value}
						{formatCapacity(machine.storageMb).unit}
					</span>

					<Separator class="col-span-6 my-2" />

					<!-- System Information -->
					<h4 class="text-md col-span-6 font-semibold">{m.system()}</h4>
					<span class="text-muted-foreground text-sm">{m.vendor()}</span>
					<span class="text-sm">{machine.hardwareInformation.system_vendor}</span>
					<span class="text-muted-foreground text-sm">{m.product()}</span>
					<span class="col-span-3 text-sm">{machine.hardwareInformation.system_product}</span>
					<span class="text-muted-foreground text-sm">{m.version()}</span>
					<span class="text-sm">{machine.hardwareInformation.system_version}</span>
					<span class="text-muted-foreground text-sm">{m.serial()}</span>
					<span class="col-span-3 text-sm">{machine.hardwareInformation.system_serial}</span>

					<Separator class="col-span-6 my-2" />

					<!-- Mainboard Information -->
					<h4 class="text-md col-span-6 font-semibold">{m.mainboard()}</h4>
					<span class="text-muted-foreground text-sm">{m.vendor()}</span>
					<span class="text-sm">{machine.hardwareInformation.mainboard_vendor}</span>
					<span class="text-muted-foreground text-sm">{m.product()}</span>
					<span class="col-span-3 text-sm">{machine.hardwareInformation.mainboard_product}</span>
					<span class="text-muted-foreground text-sm">{m.firmware()}</span>
					<span class="text-sm">{machine.hardwareInformation.mainboard_firmware_vendor}</span>
					<span class="text-muted-foreground text-sm">{m.version()}</span>
					<span class="col-span-3 text-sm">{machine.hardwareInformation.mainboard_firmware_version}</span>
					<span class="text-muted-foreground text-sm">{m.boot_mode()}</span>
					<span class="text-sm">{machine.biosBootMethod}</span>
					<span class="text-muted-foreground text-sm">{m.date()}</span>
					<span class="col-span-3 text-sm">{machine.hardwareInformation.mainboard_firmware_date}</span>

					<Separator class="col-span-6 my-2" />

					<!-- Network Information -->
					<h4 class="text-md col-span-6 font-semibold">{m.networking()}</h4>
					{#each machine.networkInterfaces as network}
						<span class="text-muted-foreground text-sm">{m.name()}</span>
						<span class="text-sm">{network.name}</span>
						<span class="text-muted-foreground text-sm">{m.mac_address()}</span>
						<span class="col-span-3 text-sm">{network.macAddress}</span>
					{/each}
				</div>
			</HoverCard.Content>
		</HoverCard.Root>
		<span>{machine.hostname}</span>
	</div>
{/snippet}
