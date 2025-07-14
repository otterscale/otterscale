<script lang="ts">
	import * as Carousel from '$lib/components/ui/carousel/index.js';
	import Autoplay from 'embla-carousel-autoplay';
	import Standalone from './standalone.svelte';
	import Cluster from './cluster.svelte';
	import Model from './model.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { setIsEnterprise } from '$lib/components/custom/enterprise';
</script>

<div class="flex w-full justify-end">
	<Button
		onclick={() => {
			setIsEnterprise(true);
		}}
	>
		Pay
	</Button>
</div>

<div class="flex items-center justify-center">
	<Carousel.Root
		plugins={[
			Autoplay({
				delay: 3000
			})
		]}
		opts={{
			align: 'start',
			loop: true
		}}
	>
		<Carousel.Content class="m-2 flex items-center justify-between gap-2 p-2">
			<Carousel.Item class="p-0 md:basis-1/2 lg:basis-1/3">
				<Standalone />
			</Carousel.Item>
			<Carousel.Item class="p-0 md:basis-1/2 lg:basis-1/3">
				<Cluster />
			</Carousel.Item>
			<Carousel.Item class="p-0 md:basis-1/2 lg:basis-1/3">
				<Model />
			</Carousel.Item>
			<Carousel.Item class="p-0 md:basis-1/2 lg:basis-1/3">
				<Model />
			</Carousel.Item>
			<Carousel.Item class="p-0 md:basis-1/2 lg:basis-1/3">
				<Model />
			</Carousel.Item>
		</Carousel.Content>
		<Carousel.Previous class="absolute bottom-0 left-0" />
		<Carousel.Next class="absolute bottom-0 right-0" />
	</Carousel.Root>
</div>
