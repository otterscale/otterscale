<script lang="ts">
	import type { ComponentProps } from 'svelte';
	import Icon from '@iconify/svelte';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { m } from '$lib/paraglide/messages.js';
	import { documentationPath } from '$lib/path';
	import DialogAbout from './dialog-about.svelte';

	interface Props extends ComponentProps<typeof Sidebar.Group> {}

	let { ref = $bindable(null), ...restProps }: Props = $props();

	let open = $state(false);
</script>

<DialogAbout bind:open />

<Sidebar.Group bind:ref {...restProps}>
	<Sidebar.GroupContent>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton size="sm">
					{#snippet child({ props })}
						<a href={documentationPath} {...props}>
							<Icon icon="ph:book-open" />
							<span>{m.documentation()}</span>
						</a>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton
					size="sm"
					onclick={() => {
						open = true;
					}}
				>
					{#snippet child({ props })}
						<a href="##" {...props}>
							<Icon icon="ph:info" />
							<span>{m.about()}</span>
						</a>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.GroupContent>
</Sidebar.Group>
