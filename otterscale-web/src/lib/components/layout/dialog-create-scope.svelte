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

	let {
		open = $bindable(false),
		trigger = $bindable(writable(false))
	}: { open: boolean; trigger: Writable<boolean> } = $props();

	let api = $state<CarouselAPI>();
	let current = $state(0);
	let openSheet = $state(false);

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
		if (createScopeRequest.name.trim()) {
			scopeClient
				.createScope(createScopeRequest)
				.then((r) => {
					toast.success(m.create_scope_success({ name: r.name }));
					trigger.set(true);
				})
				.catch((e) => {
					toast.error(m.create_scope_error({ name: createScopeRequest.name, error: e.toString() }));
				});

			open = false;
			createScopeRequest = DEFAULT_REQUEST;
		}
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
		{@const plan = plans[current]}
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
							<div class="flex space-x-2">
								{#each $machinesStore as machine}
									{#if machine.status == 'Ready'}
										<!-- required -->
										<!-- <span>{machine.fqdn}</span> -->
										<button class="h-8 w-8 rounded-full bg-slate-200"></button>
									{/if}
								{/each}
							</div>
						</div>

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
							<Select.Root type="single" required>
								<!-- bind:value={area} -->
								<Select.Trigger id="storage-devices" class="w-full">123</Select.Trigger>
								<Select.Content>
									<!-- {#each areas as area (area.value)} -->
									<Select.Item value="ffff">"fff"</Select.Item>
									<Select.Item value="aaa">"aaa"</Select.Item>
									<!-- {/each} -->
								</Select.Content>
							</Select.Root>
						</div>

						<div class="grid grid-cols-2">
							<div class="grid gap-4">
								<div class="grid gap-1">
									<Label for="calico-cidr">Calico CIDR</Label>
									<p class="text-muted-foreground text-sm">
										Network range for Kubernetes pod communication.
									</p>
								</div>
								<IPv4CIDRInput class="font-sans text-sm font-normal" placeholder="192.168.0.0/16" />
							</div>
							<div class="grid gap-4">
								<div class="grid gap-1">
									<Label for="virtual-ip">Virtual IP</Label>
									<p class="text-muted-foreground text-sm">
										High availability IP for Kubernetes control plane.
									</p>
								</div>
								<IPv4AddressInput class="font-sans text-sm font-normal" placeholder="192.168.1.1" />
							</div>
						</div>
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
				{#each plans as plan}
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
