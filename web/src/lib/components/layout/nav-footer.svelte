<script lang="ts">
	import Icon from '@iconify/svelte';
	import type { ComponentProps } from 'svelte';

	import DialogAbout from './dialog-about.svelte';

	import * as Sidebar from '$lib/components/ui/sidebar';
	import { m } from '$lib/paraglide/messages.js';
	import { staticPaths } from '$lib/path';

	type Props = ComponentProps<typeof Sidebar.Group>;

	let { ref = $bindable(null), ...restProps }: Props = $props();

	let open = $state(false);
</script>

<DialogAbout bind:open />

<Sidebar.Group bind:ref {...restProps}>
	<Sidebar.GroupContent>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton size="sm" tooltipContent={m.documentation()}>
					{#snippet child({ props })}
						<a href={staticPaths.documentation.url} target="_blank" {...props}>
							<Icon icon="ph:book-open" />
							<span>{m.documentation()}</span>
						</a>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton
					size="sm"
					tooltipContent={m.about()}
					onclick={() => {
						open = true;
					}}
				>
					{#snippet child({ props })}
						<button type="button" {...props}>
							<Icon icon="ph:info" />
							<span>{m.about()}</span>
						</button>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.GroupContent>
</Sidebar.Group>
