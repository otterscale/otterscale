<script lang="ts">
	import Autoplay from 'embla-carousel-autoplay';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Carousel from '$lib/components/ui/carousel/index.js';
</script>

<div class="flex items-center justify-center">
	<Carousel.Root
		plugins={[
			Autoplay({
				delay: 3000
			})
		]}
		opts={{
			align: 'start'
		}}
		class="w-full max-w-sm"
	>
		<Carousel.Content>
			{#each Array(5) as _, i (i)}
				<Carousel.Item class="md:basis-1/2 lg:basis-1/3">
					<div class="p-1">
						<Card.Root>
							<Card.Content class="flex aspect-square items-center justify-center p-6">
								<span class="text-3xl font-semibold">{i + 1}</span>
							</Card.Content>
						</Card.Root>
					</div>
				</Carousel.Item>
			{/each}
		</Carousel.Content>
		<Carousel.Previous />
		<Carousel.Next />
	</Carousel.Root>
</div>
