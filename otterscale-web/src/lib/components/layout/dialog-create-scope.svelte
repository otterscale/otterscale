<script lang="ts">
	import { mode } from 'mode-watcher';
	import { getContext, onMount } from 'svelte';
	import { writable, type Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { PremiumTier } from '$lib/api/premium/v1/premium_pb';
	import { ScopeService, type CreateScopeRequest } from '$lib/api/scope/v1/scope_pb';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Card from '$lib/components/ui/card';
	import * as Carousel from '$lib/components/ui/carousel';
	import type { CarouselAPI } from '$lib/components/ui/carousel/context.js';
	import { Label } from '$lib/components/ui/label';
	import * as Sheet from '$lib/components/ui/sheet';
	import { m } from '$lib/paraglide/messages';
	import { premiumTier } from '$lib/stores';
	import AdvancedTierImage from '$lib/assets/advanced-tier.jpg';
	import BasicTierImage from '$lib/assets/basic-tier.jpg';
	import EnterpriseTierImage from '$lib/assets/enterprise-tier.jpg';
	import { MachineService, type Machine } from '$lib/api/machine/v1/machine_pb';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select';
	import { IPv4AddressInput } from '$lib/components/custom/ipv4';
	import { IPv4CIDRInput } from '$lib/components/custom/ipv4-cidr';
	import { Separator } from '$lib/components/ui/separator';
	import * as RadioGroup from '$lib/components/ui/radio-group';

	import { Checkbox } from '$lib/components/ui/checkbox';

	import CalendarDaysIcon from '@lucide/svelte/icons/calendar-days';
	import * as Avatar from '$lib/components/ui/avatar';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import { dynamicPaths } from '$lib/path';
	import { page } from '$app/state';
	import { formatCapacity } from '$lib/formatter';

	let {
		open = $bindable(false),
		trigger = $bindable(writable(false))
	}: { open: boolean; trigger: Writable<boolean> } = $props();

	let api = $state<CarouselAPI>();
	let current = $state(0);
	let openSheet = $state(false);
	let selected = $state(0);
	let selectedMachine = $state('');

	$effect(() => {
		if (api) {
			current = api.selectedScrollSnap() + 1;
			api.on('select', () => {
				current = api!.selectedScrollSnap() + 1;
			});
		}
	});

	let mounted = false;
	onMount(async () => {
		try {
			await fetchMachines();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
		mounted = true;
	});

	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);
	const scopeClient = createClient(ScopeService, transport);
	const machinesStore = writable<Machine[]>([]);

	async function fetchMachines() {
		try {
			const response = await machineClient.listMachines({});
			machinesStore.set(response.machines);
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	const DEFAULT_REQUEST = { name: '' } as CreateScopeRequest;

	let createScopeRequest = $state(DEFAULT_REQUEST);

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
	}

	function handleClose() {
		open = false;
		createScopeRequest = DEFAULT_REQUEST;
	}

	interface Plan {
		tier: string;
		star: boolean;
		name: string;
		description: string;
		tags: string[];
		image: string;
		disabled: boolean;
	}

	const plans: Plan[] = [
		{
			tier: m.basic_tier(),
			star: false,
			name: m.basic_tier_name(),
			description: m.basic_tier_description(),
			tags: ['Ceph', 'Kubernetes', m.single_node()],
			image: BasicTierImage,
			disabled: $premiumTier < PremiumTier.BASIC
		},
		{
			tier: m.advanced_tier(),
			star: true,
			name: m.advanced_tier_name(),
			description: m.advanced_tier_description(),
			tags: ['Ceph', 'Multi-Node', m.multi_node(), m.cluster()],
			image: AdvancedTierImage,
			disabled: $premiumTier < PremiumTier.ADVANCED
		},
		{
			tier: m.enterprise_tier(),
			star: true,
			name: m.enterprise_tier_name(),
			description: m.enterprise_tier_description(),
			tags: ['Ceph', 'Kubernetes', m.multi_node(), m.cluster()],
			image: EnterpriseTierImage,
			disabled: $premiumTier < PremiumTier.ENTERPRISE
		}
	];
</script>

<Sheet.Root bind:open={openSheet}>
	<Sheet.Trigger>Open</Sheet.Trigger>
	<Sheet.Content class="inset-y-auto bottom-0 h-9/10 rounded-tl-lg sm:max-w-4/5">
		{@const plan = plans[selected]}
		<Sheet.Header class="h-full p-0">
			<div class="flex h-full flex-col p-12 lg:max-w-3/5">
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

					<Separator class="my-4" />
				</div>

				<!-- <div class="flex flex-col gap-6">
					<div class="flex items-center gap-3">
						<Checkbox class="h-8 w-8 rounded-full bg-slate-200" id="terms" />
					</div>
					<div class="flex items-start gap-3">
						<Checkbox id="terms-2" checked />
						<div class="grid gap-2">
							<Label for="terms-2">Accept terms and conditions</Label>
							<p class="text-muted-foreground text-sm">
								By clicking this checkbox, you agree to the terms and conditions.
							</p>
						</div>
					</div>
					<div class="flex items-start gap-3">
						<Checkbox id="toggle" disabled />
						<Label for="toggle">Enable notifications</Label>
					</div>
					<Label
						class="hover:bg-accent/50 flex items-start gap-3 rounded-lg border p-3 has-[[aria-checked=true]]:border-blue-600 has-[[aria-checked=true]]:bg-blue-50 dark:has-[[aria-checked=true]]:border-blue-900 dark:has-[[aria-checked=true]]:bg-blue-950"
					>
						<Checkbox
							id="toggle-2"
							checked
							class="data-[state=checked]:border-blue-600 data-[state=checked]:bg-blue-600 data-[state=checked]:text-white dark:data-[state=checked]:border-blue-700 dark:data-[state=checked]:bg-blue-700"
						/>
						<div class="grid gap-1.5 font-normal">
							<p class="text-sm leading-none font-medium">Enable notifications</p>
							<p class="text-muted-foreground text-sm">
								You can enable or disable notifications at any time.
							</p>
						</div>
					</Label>
				</div>

				<RadioGroup.Root value="option-one">
					<div class="flex items-center space-x-2">
						<RadioGroup.Item
							class="h-8 w-8 rounded-full bg-slate-200"
							value="option-one"
							id="option-one"
						/>
					</div>
					<div class="flex items-center space-x-2">
						<RadioGroup.Item value="option-two" id="option-two" />
						<Label for="option-two">Option Two</Label>
					</div>
				</RadioGroup.Root> -->

				<form class="flex h-full flex-col justify-between pt-4" onsubmit={handleSubmit}>
					<div class="grid gap-6">
						<div class="grid gap-4">
							<div class="grid gap-1">
								<Label for="name">
									Name
									<div class="h-2 w-2 rounded-full bg-yellow-500"></div>
								</Label>
								<p class="text-muted-foreground text-sm">
									Enter a unique identifier for your scope and prefix for facility naming.
								</p>
							</div>
							<Input id="name" type="text" placeholder="scope-name" required />
						</div>
						<div class="grid gap-4">
							<div class="grid gap-1">
								<Label for="machine">
									Machine
									<div class="h-2 w-2 rounded-full bg-yellow-500"></div>
								</Label>
								<p class="text-muted-foreground text-sm">
									Choose a machine for your Ceph and Kubernetes deployment.
								</p>
							</div>
							<Select.Root type="single" bind:value={selectedMachine} required>
								<Select.Trigger id="storage-devices" class="w-full">
									{selectedMachine
										? ($machinesStore.find((m) => m.id === selectedMachine)?.hostname ??
											'Unknown machine')
										: 'Select a machine'}
								</Select.Trigger>
								<Select.Content>
									{#each $machinesStore as machine}
										{console.log(machine)}
										{#if machine.status == 'Ready'}
											<Select.Item class="group" value={machine.id}>
												<HoverCard.Root>
													<HoverCard.Trigger
														href={dynamicPaths.machinesMetal(page.params.scope).url +
															'/' +
															machine.id}
														target="_blank"
														rel="noreferrer noopener"
														class="flex items-center space-x-1 rounded-sm underline-offset-4 group-hover:underline focus-visible:outline-2 focus-visible:outline-offset-8 focus-visible:outline-black"
													>
														<Icon icon="ph:info" class="size-4" />
													</HoverCard.Trigger>
													<HoverCard.Content align="end" class="w-120">
														<div class="grid grid-cols-4 gap-2">
															<h4 class="text-md col-span-4 font-semibold">General</h4>
															<span class="text-muted-foreground text-sm">Architecture</span>
															<span class="col-span-3 text-sm">
																{machine.architecture}
															</span>
															<span class="text-muted-foreground text-sm">CPU</span>
															<span class="col-span-3 text-sm">
																{machine.hardwareInformation.cpu_model}
															</span>
															<span class="text-muted-foreground text-sm">Memory</span>
															<span class="text-sm">
																{formatCapacity(machine.memoryMb).value}
																{formatCapacity(machine.memoryMb).unit}
															</span>
															<span class="text-muted-foreground text-sm">Storage</span>
															<span class="text-sm">
																{formatCapacity(machine.storageMb).value}
																{formatCapacity(machine.storageMb).unit}
															</span>
															<Separator class="col-span-4 my-2" />
															<h4 class="text-md col-span-4 font-semibold">System</h4>
															<span class="text-muted-foreground text-sm">Vendor</span>
															<span class="col-span-3 text-sm"
																>{machine.hardwareInformation.system_vendor}</span
															>
															<span class="text-muted-foreground text-sm">Product</span>
															<span class="col-span-3 text-sm"
																>{machine.hardwareInformation.system_product}
															</span>
															<span class="text-muted-foreground text-sm">Version</span>
															<span class="text-sm">
																{machine.hardwareInformation.system_version}
															</span>
															<span class="text-muted-foreground text-sm">Serial</span>
															<span class="text-sm">
																{machine.hardwareInformation.system_serial}
															</span>
															<Separator class="col-span-4 my-2" />
															<h4 class="text-md col-span-4 font-semibold">Mainboard</h4>
															<span class="text-muted-foreground text-sm">Vendor</span>
															<span class="text-sm">
																{machine.hardwareInformation.mainboard_vendor}
															</span>
															<span class="text-muted-foreground text-sm">Product</span>
															<span class="text-sm">
																{machine.hardwareInformation.mainboard_product}
															</span>
															<span class="text-muted-foreground text-sm">Firmware</span>
															<span class="text-sm">
																{machine.hardwareInformation.mainboard_firmware_vendor}
															</span>
															<span class="text-muted-foreground text-sm">Boot mode</span>
															<span class="text-sm">
																{machine.biosBootMethod}
															</span>
															<span class="text-muted-foreground text-sm">Version</span>
															<span class="col-span-3 text-sm">
																{machine.hardwareInformation.mainboard_firmware_version}
															</span>
															<span class="text-muted-foreground text-sm">Date</span>
															<span class="col-span-3 text-sm">
																{machine.hardwareInformation.mainboard_firmware_date}
															</span>
															<Separator class="col-span-4 my-2" />
															<h4 class="text-md col-span-4 font-semibold">Network</h4>
															{#each machine.networkInterfaces as network}
																<span class="text-muted-foreground text-sm">Name</span>
																<span class="text-sm">
																	{network.name}
																</span>
																<span class="text-muted-foreground text-sm">MAC Address</span>
																<span class="text-sm">
																	{network.macAddress}
																</span>
															{/each}
														</div>
														<!-- <div class="flex justify-between space-x-4">
															<Avatar.Root>
																<Avatar.Image src="https://github.com/sveltejs.png" />
																<Avatar.Fallback>SK</Avatar.Fallback>
															</Avatar.Root>
															<div class="space-y-1">
																<h4 class="text-sm font-semibold">@sveltejs</h4>
																<p class="text-sm">Cybernetically enhanced web apps.</p>
																<div class="flex items-center pt-2">
																	<CalendarDaysIcon class="mr-2 size-4 opacity-70" />
																	<span class="text-muted-foreground text-xs">
																		Joined September 2022
																	</span>
																</div>
															</div>
														</div> -->
													</HoverCard.Content>
												</HoverCard.Root>
												<span>{machine.hostname}</span>
											</Select.Item>
										{/if}
									{/each}
								</Select.Content>
							</Select.Root>
						</div>

						{#if selectedMachine}
							<div class="grid gap-4">
								<div class="grid gap-1">
									<Label for="storage-devices">
										Storage Devices
										<div class="h-2 w-2 rounded-full bg-yellow-500"></div>
									</Label>
									<p class="text-muted-foreground text-sm">
										Configure dedicated storage devices for Ceph.
									</p>
								</div>
								{#each $machinesStore.find((m) => m.id === selectedMachine)?.blockDevices ?? [] as device}
									<span>/dev/{device.name}</span>
								{/each}
							</div>

							<div class="grid grid-cols-2">
								<div class="grid gap-4">
									<div class="grid gap-1">
										<Label for="calico-cidr">Calico CIDR</Label>
										<p class="text-muted-foreground text-sm">
											Network range for Kubernetes pod communication.
										</p>
									</div>
									<IPv4CIDRInput
										class="font-sans text-sm font-normal"
										placeholder="192.168.0.0/16"
									/>
								</div>
								<div class="grid gap-4">
									<div class="grid gap-1">
										<Label for="virtual-ip">Virtual IP</Label>
										<p class="text-muted-foreground text-sm">
											High availability IP for Kubernetes control plane.
										</p>
									</div>
									<IPv4AddressInput
										class="font-sans text-sm font-normal"
										placeholder="192.168.1.1"
									/>
								</div>
							</div>
						{/if}
					</div>

					<div class="flex gap-8">
						<Button
							size="lg"
							variant="outline"
							class="flex-1"
							onclick={() => {
								openSheet = false;
							}}
						>
							Cancel
						</Button>
						<Button type="submit" size="lg" class="flex-1">Confirm</Button>
					</div>
				</form>
			</div>

			<div class="relative lg:absolute lg:inset-y-0 lg:right-0 lg:w-2/5">
				<img src={plan.image} alt={plan.name} class="absolute inset-0 size-full object-cover" />
			</div>
		</Sheet.Header>
	</Sheet.Content>
</Sheet.Root>

<Dialog.Root bind:open onOpenChange={handleClose}>
	<Dialog.Content showCloseButton={false} class="min-w-4xl overflow-hidden border-0 p-0">
		<Carousel.Root setApi={(emblaApi) => (api = emblaApi)}>
			<Carousel.Content>
				{#each plans as plan, index}
					<Carousel.Item>
						<Card.Root class="relative aspect-[21/9] rounded-lg border-0 shadow-none">
							<Card.Content class="flex items-center justify-center rounded-lg">
								<div class="absolute inset-0 rounded-lg transition-transform duration-500 ease-out">
									<img src={plan.image} alt={plan.name} class="object-cover" />
									<div
										class="from-background/90 via-background/50 absolute inset-0 overflow-hidden bg-gradient-to-r to-transparent"
									></div>

									<div class="absolute top-12 left-12 flex min-h-76 flex-col justify-between">
										<div class="flex max-w-2xl flex-col space-y-4">
											<Badge
												variant="secondary"
												class="bg-primary/10 text-primary flex items-center uppercase"
											>
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
										</div>

										<div class="space-y-2 pt-4">
											{#if plan.disabled}
												<Label for="install" class="tracking-tight text-red-500">
													<Icon icon="ph:info-bold" class="size-4" />
													{m.requires_subscription()}
												</Label>
											{/if}
											<Button
												id="install"
												size="lg"
												disabled={plan.disabled}
												onclick={() => {
													open = false;
													openSheet = true;
													selected = index;
												}}
											>
												<Icon icon="ph:download-bold" />
												{m.install()}
											</Button>
										</div>
									</div>
								</div>
							</Card.Content>
						</Card.Root>
					</Carousel.Item>
				{/each}
			</Carousel.Content>

			<div class="absolute top-14 right-12 flex items-center space-x-2">
				{#each plans as _, index}
					<button
						onclick={() => api?.scrollTo(index)}
						aria-label="Go to slide {index + 1}"
						class="size-2 rounded-full transition-all {index + 1 === current
							? 'bg-primary w-6'
							: 'bg-primary/30 hover:bg-primary/50'}"
					></button>
				{/each}
			</div>

			<div class="absolute right-16 bottom-12 flex items-center">
				<Carousel.Previous
					variant={mode.current === 'dark' ? 'default' : 'outline'}
					class="top-1/2 -left-12 rounded-md"
				/>
				<Carousel.Next
					variant={mode.current === 'dark' ? 'default' : 'outline'}
					class="top-1/2 -right-6 rounded-md"
				/>
			</div>
		</Carousel.Root>
	</Dialog.Content>
</Dialog.Root>
