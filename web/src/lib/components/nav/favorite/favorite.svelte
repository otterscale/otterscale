<script lang="ts">
	import Icon from '@iconify/svelte';
	import { toast } from 'svelte-sonner';

	import { Button } from '$lib/components/ui/button';

	let favorited = false;
	function toggleFavorite() {
		if (!favorited) {
			favorited = true;
			toast.success('Added to favorites!');
			return;
		}
		favorited = false;
		toast.error('Removed from favorites!');
	}
</script>

<Button variant="outline" size="icon" class="bg-header" on:click={toggleFavorite}>
	{#if favorited}
		<Icon icon="ph:heart-fill" class="h-5 w-5" />
	{:else}
		<Icon icon="ph:heart" class="h-5 w-5" />
	{/if}
</Button>
