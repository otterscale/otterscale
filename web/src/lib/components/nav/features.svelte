<script lang="ts">
	import { page } from '$app/state';

	import { cn } from '$lib/utils.js';
	import * as m from '$lib/paraglide/messages.js';

	const features = [
		{
			path: '/tutorial',
			name: m.nav_tutorial(),
			enable: false
		},
		{
			path: '/data-fabric',
			name: m.nav_data_fabric(),
			enable: true
		},
		{
			path: '/explore',
			name: m.nav_explore(),
			enable: false
		},
		{
			path: '/dashboard',
			name: m.nav_dashboard(),
			enable: true
		},
		{
			path: '/applications',
			name: m.nav_applications(),
			enable: true
		},
		{
			path: '/integrations',
			name: m.nav_integrations(),
			enable: false
		},
		{
			path: '/dev-tools',
			name: m.nav_dev_tools(),
			enable: false
		}
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
		<span>{feature.name}</span>
	</a>
{/each}
