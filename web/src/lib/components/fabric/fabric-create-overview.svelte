<script lang="ts">
	import Icon from '@iconify/svelte';

	import * as Carousel from '$lib/components/ui/carousel';
	import * as Card from '$lib/components/ui/card';

	let {
		items = $bindable()
	}: {
		items: {
			name: string;
			icon: string;
			active: boolean;
		}[];
	} = $props();
</script>

<Carousel.Root
	opts={{
		align: 'center'
	}}
	class="w-full max-w-xs md:max-w-md lg:max-w-lg"
>
	<Carousel.Content>
		{#each items as item}
			<Carousel.Item class="md:basis-1/2 lg:basis-1/3">
				<div class="space-y-4 p-2">
					<Card.Root
						class="hover:bg-accent"
						onclick={() => {
							item.active = !item.active;
						}}
					>
						<Card.Content class="flex aspect-square items-center justify-center p-6">
							<span class="text-4xl font-semibold">
								<Icon icon={item.icon} />
							</span>
						</Card.Content>
					</Card.Root>
					<div class="text-foregroundf text-center text-sm">{item.name}</div>
				</div>
			</Carousel.Item>
		{/each}
	</Carousel.Content>
	<Carousel.Previous />
	<Carousel.Next />
</Carousel.Root>
