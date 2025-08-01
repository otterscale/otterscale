<script lang="ts">
	import Icon from '@iconify/svelte';
	import { PremiumTier } from '$lib/api/premium/v1/premium_pb';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as Carousel from '$lib/components/ui/carousel';
	import type { CarouselAPI } from '$lib/components/ui/carousel/context.js';
	import { Label } from '$lib/components/ui/label';
	import { m } from '$lib/paraglide/messages';
	import {
		setupPath,
		setupScopeCephPath,
		setupScopeKubernetesPath,
		setupScopePath
	} from '$lib/path';
	import AdvancedTierImage from '$lib/static/advanced-tier.jpg';
	import BasicTierImage from '$lib/static/basic-tier.jpg';
	import EnterpriseTierImage from '$lib/static/enterprise-tier.jpg';
	import { breadcrumb, currentCeph, currentKubernetes, premiumTier } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [setupPath], current: setupScopePath });

	let api = $state<CarouselAPI>();

	let current = $state(0);
	$effect(() => {
		if (api) {
			current = api.selectedScrollSnap() + 1;
			api.on('select', () => {
				current = api!.selectedScrollSnap() + 1;
			});
		}
	});

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

<!-- just-in-time  -->
<dummy class="bg-[#326de6]"></dummy>
<dummy class="bg-[#f0424d]"></dummy>

<h2 class="text-center text-3xl font-bold tracking-tight sm:text-4xl">{m.setup_scope()}</h2>

{#if $currentKubernetes || $currentCeph}
	<p class="text-muted-foreground mt-4 text-center text-lg">
		{m.setup_scope_configured_description()}
	</p>
	<div class="mx-auto max-w-5xl px-4 py-10 xl:px-0">
		<div class="bg-card rounded-xl border shadow-sm">
			<div class="rounded-xl p-4 lg:p-8">
				<div class="grid min-w-2xl grid-cols-1 items-center gap-x-12 gap-y-20 sm:grid-cols-2">
					<div
						class="before:bg-border relative text-center before:absolute before:start-1/2 before:-top-full before:mt-3.5 before:h-20 before:w-px before:-translate-x-1/2 before:rotate-[60deg] before:transform before:shadow-sm first:before:hidden sm:before:-start-6 sm:before:top-1/2 sm:before:mt-0 sm:before:-translate-x-0 sm:before:-translate-y-1/2 sm:before:rotate-12"
					>
						<div class="space-y-2">
							<Icon icon="simple-icons:ceph" class="mx-auto size-14 shrink-0 text-[#f0424d]" />
							<h3 class="text-lg font-semibold sm:text-2xl">Ceph</h3>
						</div>
						<Button
							href={setupScopeCephPath}
							variant="ghost"
							class="text-muted-foreground text-sm sm:text-base"
						>
							<Icon icon="ph:wrench" class="size-5" />
							{$currentCeph ? $currentCeph.name : '-'}
						</Button>
					</div>
					<div
						class="before:bg-border relative text-center before:absolute before:start-1/2 before:-top-full before:mt-3.5 before:h-20 before:w-px before:-translate-x-1/2 before:rotate-[60deg] before:transform before:shadow-sm first:before:hidden sm:before:-start-6 sm:before:top-1/2 sm:before:mt-0 sm:before:-translate-x-0 sm:before:-translate-y-1/2 sm:before:rotate-12"
					>
						<div class="space-y-2">
							<Icon
								icon="simple-icons:kubernetes"
								class="mx-auto size-14 shrink-0 text-[#326de6]"
							/>
							<h3 class="text-lg font-semibold sm:text-2xl">Kubernetes</h3>
						</div>
						<Button
							href={setupScopeKubernetesPath}
							variant="ghost"
							class="text-muted-foreground text-sm sm:text-base"
						>
							<Icon icon="ph:wrench" class="size-5" />
							{$currentKubernetes ? $currentKubernetes.name : '-'}
						</Button>
					</div>
				</div>
			</div>
		</div>
	</div>
{:else}
	<p class="text-muted-foreground mt-4 text-center text-lg">
		{m.setup_scope_not_configured_description()}
	</p>
	<div class="mx-auto w-full max-w-5xl px-4 py-10 xl:px-0">
		<Carousel.Root setApi={(emblaApi) => (api = emblaApi)} class="w-full">
			<Carousel.Content>
				{#each plans as plan}
					<Carousel.Item>
						<Card.Root class="relative aspect-[21/9] overflow-hidden rounded-xl shadow-none">
							<Card.Content class="flex items-center justify-center">
								<div class="absolute inset-0 transition-transform duration-500 ease-out">
									<div class="absolute inset-0">
										<img src={plan.image} alt={plan.name} class="object-cover" />
										<div
											class="from-background/90 via-background/50 absolute inset-0 bg-gradient-to-r to-transparent"
										></div>
									</div>

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
												<Label for="install" class="text-red-500">
													<Icon icon="ph:info" />
													{m.requires_subscription()}
												</Label>
											{/if}
											<Button id="install" size="lg" disabled={plan.disabled}>
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

			<!-- Navigation dots -->
			<div class="absolute bottom-10 left-12 flex items-center space-x-2">
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

			<!-- Previous/Next buttons -->
			<div class="absolute right-16 bottom-12 flex items-center">
				<Carousel.Previous class="top-1/2 -left-12 rounded-md" />
				<Carousel.Next class="top-1/2 -right-6 rounded-md" />
			</div>
		</Carousel.Root>
	</div>
{/if}
