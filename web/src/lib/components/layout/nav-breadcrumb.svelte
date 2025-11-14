<script lang="ts">
	import * as Breadcrumb from '$lib/components/ui/breadcrumb';
	import { breadcrumbs } from '$lib/stores';

	const current = $derived($breadcrumbs.at(-1));
</script>

<nav aria-label="Breadcrumb">
	<Breadcrumb.Root>
		<Breadcrumb.List>
			{#each $breadcrumbs as item (item.url)}
				{#if item.url === current?.url}
					<Breadcrumb.Item>
						<Breadcrumb.Page>{current.title}</Breadcrumb.Page>
					</Breadcrumb.Item>
				{:else}
					<Breadcrumb.Item class="hidden md:block">
						<Breadcrumb.Link href={item.url}>
							{item.title}
						</Breadcrumb.Link>
					</Breadcrumb.Item>
					<Breadcrumb.Separator class="hidden md:block" />
				{/if}
			{/each}
		</Breadcrumb.List>
	</Breadcrumb.Root>
</nav>
