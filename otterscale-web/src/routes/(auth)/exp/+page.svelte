<script lang="ts">
	import Icon from '@iconify/svelte';
	import { PremiumTier } from '$lib/api/premium/v1/premium_pb';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as Carousel from '$lib/components/ui/carousel';
	import type { CarouselAPI } from '$lib/components/ui/carousel/context.js';
	import { Label } from '$lib/components/ui/label';
	import AdvancedTierImage from '$lib/static/advanced-tier.jpg';
	import BasicTierImage from '$lib/static/basic-tier.jpg';
	import EnterpriseTierImage from '$lib/static/enterprise-tier.jpg';
	import { premiumTier } from '$lib/stores';
	import { m } from '$lib/paraglide/messages';

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

<div class="mx-auto w-full max-w-5xl p-6">
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

								<div class="absolute top-12 left-12 flex min-h-72 flex-col justify-between">
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
		<div class="absolute right-12 bottom-12 flex items-center">
			<Carousel.Previous class="top-1/2 -left-12 rounded-md" />
			<Carousel.Next class="top-1/2 -right-6 rounded-md" />
		</div>
	</Carousel.Root>
</div>
