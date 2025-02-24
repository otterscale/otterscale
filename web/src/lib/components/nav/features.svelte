<script lang="ts">
	import { page } from '$app/state';
	import { cn, getFeatureTitle } from '$lib/utils.js';

	export const features = [
		{ path: '/tutorial', enable: false },
		{ path: '/data-fabric', enable: true },
		{ path: '/explore', enable: false },
		{ path: '/dashboard', enable: true },
		{ path: '/applications', enable: true },
		{ path: '/integrations', enable: false },
		{ path: '/dev-tools', enable: false }
	];
</script>

{#each features as feature}
	<a
		href={feature.path}
		class={cn(
			'transition-colors hover:text-foreground',
			page.url.pathname.startsWith(feature.path) ? 'text-foreground' : 'text-muted-foreground'
		)}
	>
		<span>{getFeatureTitle(feature.path)}</span>
	</a>
{/each}
