<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { ScopeService, type CreateScopeRequest } from '$lib/api/scope/v1/scope_pb';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import * as Sheet from '$lib/components/ui/sheet';
	import { MachineService, type Machine } from '$lib/api/machine/v1/machine_pb';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select';
	import { IPv4AddressInput } from '$lib/components/custom/ipv4';
	import { IPv4CIDRInput } from '$lib/components/custom/ipv4-cidr';
	import { Separator } from '$lib/components/ui/separator';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import type { Plan } from './plans';

	let { open = $bindable(false), plan = $bindable({} as Plan) }: { open: boolean; plan: Plan } = $props();

	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);
	const scopeClient = createClient(ScopeService, transport);
	const machinesStore = writable<Machine[]>([]);
	const defaultCreateScopeRequest = { name: '' } as CreateScopeRequest;

	let selectedMachine = $state('');
	let createScopeRequest = $state(defaultCreateScopeRequest);

	async function fetchMachines() {
		try {
			const response = await machineClient.listMachines({});
			machinesStore.set(response.machines);
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	async function createScope() {
		try {
			const response = await scopeClient.createScope(createScopeRequest);
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	// TODO
	function handleSubmit() {
		// if (createScopeRequest.name.trim()) {
		// 	scopeClient
		// 		.createScope(createScopeRequest)
		// 		.then((r) => {
		// 			toast.success(m.create_scope_success({ name: r.name }));
		// 			trigger.set(true);
		// 		})
		// 		.catch((e) => {
		// 			toast.error(m.create_scope_error({ name: createScopeRequest.name, error: e.toString() }));
		// 		});
		// 	open = false;
		// 	createScopeRequest = DEFAULT_REQUEST;
		// }
		handleClose();
		// trigger.set(true);
		goto(dynamicPaths.setupScope('zzz').url);
	}

	function handleClose() {
		open = false;
		createScopeRequest = defaultCreateScopeRequest;
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
							<Input id="name" type="text" placeholder="scope-name" required />
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
														class="hover:bg-accent/50 flex items-start gap-x-2 rounded-md border p-2 has-[[aria-checked=true]]:border-slate-600 has-[[aria-checked=true]]:bg-blue-50 dark:has-[[aria-checked=true]]:border-slate-900 dark:has-[[aria-checked=true]]:bg-slate-950"
													>
														<Checkbox
															id={device.name}
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
									<IPv4CIDRInput class="font-sans text-sm font-normal" placeholder="192.168.0.0/16" />
								</div>
								<div class="grid gap-4">
									<div class="grid gap-1">
										<Label for="virtual-ip">{m.create_scope_virtual_ip()} ({m.optional()})</Label>
										<p class="text-muted-foreground text-sm">
											{m.create_scope_virtual_ip_description()}
										</p>
									</div>
									<IPv4AddressInput class="font-sans text-sm font-normal" placeholder="192.168.1.1" />
								</div>
							</div>
						{/if}
					</div>

					<!-- Form Actions -->
					<div class="flex gap-8">
						<Button size="lg" variant="outline" class="flex-1" onclick={handleClose}>
							{m.cancel()}
						</Button>
						<Button type="submit" size="lg" class="flex-1">{m.create()}</Button>
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
