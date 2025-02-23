<script lang="ts">
	import { onMount } from 'svelte';
	import Icon from '@iconify/svelte';
	import { toast } from 'svelte-sonner';

	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { addFavorite, deleteFavorite, isFavorite } from '$lib/pb';

	let favorited = false;

	onMount(async () => {
		favorited = await isFavorite();
	});

	async function add() {
		await addFavorite();
		toast.success('Added to favorites!');
	}

	async function del() {
		await deleteFavorite();
		toast.success('Removed from favorites!');
	}

	async function toggleFavorite() {
		favorited ? await del() : await add();
		favorited = await isFavorite();
	}
</script>

<Tooltip.Root>
	<Tooltip.Trigger>
		<Button variant="outline" size="icon" class="bg-header" on:click={toggleFavorite}>
			{#if favorited}
				<Icon icon="ph:heart-fill" class="h-5 w-5" />
			{:else}
				<Icon icon="ph:heart" class="h-5 w-5" />
			{/if}
		</Button>
	</Tooltip.Trigger>
	<Tooltip.Content>
		<p>Favorite</p>
	</Tooltip.Content>
</Tooltip.Root>
