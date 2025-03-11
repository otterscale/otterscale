<script lang="ts">
	import Platform from './platform.svelte';
	import Shortcut from './shortcut.svelte';
	import Namespace from './namespace.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { ComponentProps } from 'svelte';
	import Footer from './footer.svelte';
	import pb, { avatarFallback, avatarURL } from '$lib/pb';
	import { featureTitle } from '../features';

	const data = {
		user: {
			name: pb.authStore.record?.name,
			email: pb.authStore.record?.email,
			avatar: avatarURL(pb.authStore.record?.avatar),
			fallback: avatarFallback(pb.authStore.record?.name)
		},
		namespaces: [
			{
				name: 'Default',
				plan: 'Enterprise',
				icon: 'ph:graph',
				color: '#4682B4'
			},
			{
				name: 'Team A',
				plan: 'Free',
				icon: 'ph:airplane-tilt',
				color: '#FF6347'
			},
			{
				name: 'Team B',
				plan: 'Enterprise',
				icon: 'ph:flower',
				color: '#FFD700'
			}
		],
		platforms: [
			{
				title: featureTitle('/data-fabric'),
				url: '#',
				icon: 'ph:tree-structure',
				isActive: true,
				items: [
					{
						title: 'Browse',
						url: '/data-fabric'
					},
					{
						title: 'Configurations',
						url: '/data-fabric/configurations'
					},
					{
						title: 'Uptime',
						url: '/data-fabric/uptime'
					}
				]
			},
			{
				title: featureTitle('/dashboard'),
				url: '#',
				icon: 'ph:chart-line-up',
				isActive: true,
				items: [
					{
						title: 'My Data',
						url: '/dashboard/my'
					},
					{
						title: 'Shared with Me',
						url: '/dashboard/shared'
					}
				]
			},
			{
				title: 'Models',
				url: '#',
				icon: 'ph:robot',
				items: [
					{
						title: 'Genesis',
						url: '#'
					},
					{
						title: 'Explorer',
						url: '#'
					},
					{
						title: 'Quantum',
						url: '#'
					}
				]
			},
			{
				title: 'Settings',
				url: '#',
				icon: 'ph:gear',
				items: [
					{
						title: 'General',
						url: '#'
					},
					{
						title: 'Team',
						url: '#'
					},
					{
						title: 'Billing',
						url: '#'
					},
					{
						title: 'Limits',
						url: '#'
					}
				]
			}
		],
		footers: [
			{
				title: 'Support',
				url: 'https://openhdc.github.io',
				icon: 'ph:lifebuoy'
			},
			{
				title: 'About',
				url: '/about',
				icon: 'ph:info'
			}
		],
		shortcuts: [
			{
				name: featureTitle('/favorites'),
				url: '/favorites',
				icon: 'ph:clover'
			},
			{
				name: featureTitle('/recents'),
				url: '/recents',
				icon: 'ph:clock'
			}
		]
	};

	let {
		ref = $bindable(null),
		collapsible = 'icon',
		...restProps
	}: ComponentProps<typeof Sidebar.Root> = $props();
</script>

<Sidebar.Root bind:ref {collapsible} {...restProps}>
	<Sidebar.Header>
		<Namespace namespaces={data.namespaces} />
	</Sidebar.Header>
	<Sidebar.Content>
		<Shortcut items={data.shortcuts} />
		<Platform items={data.platforms} />
	</Sidebar.Content>
	<Sidebar.Footer>
		<Sidebar.Separator />
		<Footer items={data.footers} />
	</Sidebar.Footer>
	<Sidebar.Rail />
</Sidebar.Root>
