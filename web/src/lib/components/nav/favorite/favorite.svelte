<script lang="ts">
	import { onMount } from 'svelte';
	import Icon from '@iconify/svelte';
	import { toast } from 'svelte-sonner';

	import { page } from '$app/state';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import pb, { addFavorite, deleteFavorite, isFavorite } from '$lib/pb';
	import { i18n } from '$lib/i18n';
	import { cn } from '$lib/utils';

	let favorited = $state(false);
	let currentPage = i18n.route(page.url.pathname);

	onMount(async () => {
		favorited = await isFavorite();

		pb.collection('favorites').subscribe('*', (r) => {
			if (r.record.path === i18n.route(page.url.pathname)) {
				favorited = r.action === 'create';
			}
		});
	});

	$effect(() => {
		if (currentPage !== i18n.route(page.url.pathname)) {
			(async () => (favorited = await isFavorite()))();
			currentPage = i18n.route(page.url.pathname);
		}
	});

	async function toggleFavorite() {
		favorited
			? await deleteFavorite().then(() => toast.success('Removed from favorites!'))
			: await addFavorite().then(() => toast.success('Added to favorites!'));
	}
</script>

<Tooltip.Provider>
	<Tooltip.Root>
		<Tooltip.Trigger
			class={cn(buttonVariants({ variant: 'outline', size: 'icon' }), 'bg-header [&_svg]:size-5')}
			onclick={toggleFavorite}
		>
			{#if favorited}
				<Icon icon="ph:heart-fill" />
			{:else}
				<Icon icon="ph:heart" />
			{/if}
		</Tooltip.Trigger>
		<Tooltip.Content>
			<p>Favorite</p>
		</Tooltip.Content>
	</Tooltip.Root>
</Tooltip.Provider>
