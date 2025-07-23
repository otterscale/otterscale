<script lang="ts">
	import type { ComponentProps } from 'svelte';
	import Icon from '@iconify/svelte';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { getIconFromUrl } from './icon';

	interface NavItem {
		title: string;
		url: string;
	}

	interface Props extends ComponentProps<typeof Sidebar.Group> {
		items: NavItem[];
	}

	let { ref = $bindable(null), items, ...restProps }: Props = $props();
</script>

<Sidebar.Group bind:ref {...restProps}>
	<Sidebar.GroupContent>
		<Sidebar.Menu>
			{#each items as item (item.title)}
				<Sidebar.MenuItem>
					<Sidebar.MenuButton size="sm">
						{#snippet child({ props })}
							<a href={item.url} {...props}>
								<Icon icon={getIconFromUrl(item.url)} />
								<span>{item.title}</span>
							</a>
						{/snippet}
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
			{/each}
		</Sidebar.Menu>
	</Sidebar.GroupContent>
</Sidebar.Group>
