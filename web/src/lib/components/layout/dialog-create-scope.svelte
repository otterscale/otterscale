<script lang="ts">
	import { mode } from 'mode-watcher';
	import { writable, type Writable } from 'svelte/store';
	import Icon from '@iconify/svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Card from '$lib/components/ui/card';
	import * as Carousel from '$lib/components/ui/carousel';
	import type { CarouselAPI } from '$lib/components/ui/carousel/context.js';
	import { Label } from '$lib/components/ui/label';
	import { m } from '$lib/paraglide/messages';
	import { plans } from './plans';
	import SheetCreateScope from './sheet-create-scope.svelte';

	let {
		open = $bindable(false),
		trigger = $bindable(writable(false))
	}: { open: boolean; trigger: Writable<boolean> } = $props();

	let api = $state<CarouselAPI>();
	let current = $state(0);
	let openSheet = $state(false);
	let selected = $state(-1);

	$effect(() => {
		if (api) {
			current = api.selectedScrollSnap() + 1;
			api.on('select', () => {
				current = api!.selectedScrollSnap() + 1;
			});
		}
	});

	function handlePlanSelect(index: number) {
		open = false;
		openSheet = true;
		selected = index;
	}
</script>

<SheetCreateScope bind:open={openSheet} plan={plans[selected]} />

<Dialog.Root bind:open>
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
												onclick={() => handlePlanSelect(index)}
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
