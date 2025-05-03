<script lang="ts">
	import { i18n } from '$lib/i18n';
	import { page } from '$app/state';
	import { ParaglideJS } from '@inlang/paraglide-sveltekit';

	let { children } = $props();

	let items = ['Profile', 'Billing', 'Appearance', 'Notification', 'Advanced'];
</script>

<main
	class="flex min-h-[calc(100vh_-_theme(spacing.16))] flex-1 flex-col gap-4 bg-muted/40 p-4 md:gap-8 md:p-10"
>
	<div class="mx-auto grid w-full max-w-6xl gap-2">
		<h1 class="text-3xl font-semibold">Settings</h1>
	</div>
	<div
		class="mx-auto grid w-full max-w-6xl items-start gap-6 md:grid-cols-[180px_1fr] lg:grid-cols-[250px_1fr]"
	>
		<nav
			class="grid gap-4 text-sm text-muted-foreground"
			data-x-chunk-container="chunk-container after:right-0"
		>
			{#each items as item}
				{@const url = '/settings/' + item.toLowerCase()}
				{#if page.url.pathname === url}
					<a href={url} class="font-semibold text-primary">{item}</a>
				{:else}
					<a href={url}>{item}</a>
				{/if}
			{/each}
		</nav>
		<div class="grid gap-6">
			<ParaglideJS {i18n}>
				{@render children()}
			</ParaglideJS>
		</div>
	</div>
</main>
