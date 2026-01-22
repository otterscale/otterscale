<script lang="ts">
	import BookOpenIcon from '@lucide/svelte/icons/book-open';
	import InfoIcon from '@lucide/svelte/icons/info';
	import type { ComponentProps } from 'svelte';

	import * as Sidebar from '$lib/components/ui/sidebar';
	import { m } from '$lib/paraglide/messages.js';

	import DialogAbout from './dialog-about.svelte';

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
						<a href="https://otterscale.github.io" target="_blank" {...props}>
							<BookOpenIcon />
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
							<InfoIcon />
							<span>{m.about()}</span>
						</button>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.GroupContent>
</Sidebar.Group>
