<script lang="ts">
	import { page } from '$app/state';
	import { cn, getFeatureTitle } from '$lib/utils.js';
	import { Button } from '$lib/components/ui/button';
	import { goto } from '$app/navigation';
	import { i18n } from '$lib/i18n';

	export const features = [
		// { path: '/tutorial', enable: false },
		{ path: '/data-fabric', enable: true },
		{ path: '/explore', enable: false },
		{ path: '/dashboard', enable: true },
		{ path: '/applications', enable: true },
		{ path: '/integrations', enable: false },
		{ path: '/dev-tools', enable: false }
	];
</script>

{#each features as feature}
	<Button
		variant="ghost"
		class={cn(
			'px-2 py-0 transition-colors',
			!feature.enable ? 'disabled:pointer-events-auto disabled:cursor-not-allowed' : '',
			i18n.route(page.url.pathname).startsWith(feature.path)
				? 'text-foreground'
				: 'text-muted-foreground'
		)}
		disabled={!feature.enable}
		onclick={() => goto(i18n.resolveRoute(feature.path))}
	>
		<span>{getFeatureTitle(feature.path)}</span>
	</Button>
{/each}
