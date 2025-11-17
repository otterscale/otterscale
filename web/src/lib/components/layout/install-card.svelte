<script lang="ts">
	import Icon from '@iconify/svelte';
	import { mode } from 'mode-watcher';

	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as Carousel from '$lib/components/ui/carousel';
	import type { CarouselAPI } from '$lib/components/ui/carousel/context.js';
	import { Label } from '$lib/components/ui/label';
	import { m } from '$lib/paraglide/messages';

	import { plans } from './plans';

	let { onSelect }: { onSelect: (index: number) => void } = $props();

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
</script>

<Carousel.Root setApi={(emblaApi) => (api = emblaApi)}>
	<Carousel.Content>
		{#each plans as plan, index (plan.name)}
			<Carousel.Item>
				<Card.Root class="relative aspect-21/9 rounded-lg border-0 shadow-none">
					<Card.Content class="flex items-center justify-center rounded-lg">
						<div class="absolute inset-0 rounded-lg transition-transform duration-500 ease-out">
							<img src={plan.image} alt={plan.name} class="object-cover" />
							<div
								class="absolute inset-0 overflow-hidden bg-linear-to-r from-background/90 via-background/50 to-transparent"
							></div>

							<div class="absolute top-12 left-12 flex max-w-2xl flex-col space-y-4">
								<Badge
									variant="secondary"
									class="flex items-center bg-primary/10 text-primary uppercase"
								>
									{#if plan.star}
										<Icon icon="ph:star-fill" class="text-yellow-500" />
									{/if}
									<span>{plan.tier}</span>
								</Badge>

								<h2 class="text-3xl font-semibold tracking-tight">{plan.name}</h2>
								<p class="text-md text-accent-foreground/80">{plan.description}</p>

								<div class="flex flex-wrap gap-2">
									{#each plan.tags as tag (tag)}
										<Badge variant="outline" class="bg-background/50 backdrop-blur-sm">
											{tag}
										</Badge>
									{/each}
								</div>
							</div>

							<div class="absolute bottom-12 left-12 space-y-2">
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
									onclick={() => onSelect(index)}
								>
									<Icon icon="ph:download-bold" />
									{m.install()}
								</Button>
							</div>
						</div>
					</Card.Content>
				</Card.Root>
			</Carousel.Item>
		{/each}
	</Carousel.Content>

	<div class="absolute top-14 right-12 flex items-center space-x-2">
		<!-- eslint-disable-next-line @typescript-eslint/no-unused-vars -->
		{#each plans as _, index (index)}
			<button
				onclick={() => api?.scrollTo(index)}
				aria-label="Go to slide {index + 1}"
				class="size-2 rounded-full transition-all {index + 1 === current
					? 'w-6 bg-primary'
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
