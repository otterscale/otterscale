<script lang="ts">
	import type { ComponentProps } from 'svelte';
	import Icon from '@iconify/svelte';
	import { page } from '$app/state';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { m } from '$lib/paraglide/messages.js';
	import { dynamicPaths, staticPaths } from '$lib/path';
	import changelogRead from '$lib/stores/changelog';
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
						<a href={staticPaths.documentation.url} {...props}>
							<Icon icon="ph:book-open" />
							<span>{m.documentation()}</span>
						</a>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton size="sm">
					{#snippet child({ props })}
						<a href={dynamicPaths.changelog(page.params.scope).url} {...props}>
							<Icon icon="ph:clock-counter-clockwise" />
							<span>{m.changelog()}</span>
							{#if !$changelogRead}
								<span class="relative flex size-2">
									<span
										class="absolute inline-flex h-full w-full animate-ping rounded-full bg-blue-400 opacity-75"
									></span>
									<span class="relative inline-flex size-2 rounded-full bg-blue-500"></span>
								</span>
							{/if}
							<!-- prevent [&>span:last-child]:truncate -->
							<span></span>
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
