<script lang="ts">
	import type { ComponentProps } from 'svelte';
	import Icon from '@iconify/svelte';
	import * as Sidebar from '$lib/components/ui/sidebar';

	interface NavItem {
		title: string;
		url: string;
		icon: string;
	}

	let {
		ref = $bindable(null),
		items,
		...restProps
	}: {
		items: NavItem[];
	} & ComponentProps<typeof Sidebar.Group> = $props();
</script>

<Sidebar.Group bind:ref {...restProps}>
	<Sidebar.GroupContent>
		<Sidebar.Menu>
			{#each items as { title, url, icon } (title)}
				<Sidebar.MenuItem>
					<Sidebar.MenuButton size="sm">
						{#snippet child({ props })}
							<a href={url} {...props}>
								<Icon {icon} />
								<span>{title}</span>
							</a>
						{/snippet}
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
			{/each}
		</Sidebar.Menu>
	</Sidebar.GroupContent>
</Sidebar.Group>
