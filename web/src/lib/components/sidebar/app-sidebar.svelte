<script lang="ts">
	import Section from './section.svelte';
	import Namespace from './namespace.svelte';
	import Footer from './footer.svelte';
	import { featureTitle } from '../features';
	import type { Title } from '../ui/sheet';

	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { ComponentProps } from 'svelte';

	const data = {
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
		analysis: [
			{
				title: featureTitle('/data-fabric'),
				url: '#',
				icon: 'ph:tree-structure',
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
			}
		],
		platforms: [
			{
				title: 'Monitor',
				url: '/management',
				icon: 'ph:gauge'
			},
			{
				title: 'Orchestration',
				url: '/orchestration',
				icon: 'ph:tree-structure'
			},
			{
				title: 'Management',
				url: '#',
				icon: 'ph:command',
				isActive: true,
				items: [
					{
						title: 'Model',
						url: '#'
					},
					{
						title: 'Application',
						url: '/management/application'
					},
					{
						title: 'Facility',
						url: '/management/scope'
					},
					{
						title: 'Machine',
						url: '/management/machine?intervals=15'
					},

					{
						title: 'General',
						url: '/management/general'
					}
				]
			},
			{
				title: 'Market',
				url: '/market',
				icon: 'ph:magnifying-glass'
			}
		],
		general: [
			{
				title: 'Settings',
				url: '#',
				icon: 'ph:gear',
				items: [
					{
						title: 'Profile',
						url: '/settings/profile'
					},
					{
						title: 'Billing',
						url: '/settings/billing'
					},
					{
						title: 'Appearance',
						url: '/settings/appearance'
					},
					{
						title: 'Notification',
						url: '/settings/notification'
					},
					{
						title: 'Advanced',
						url: '/settings/advanced'
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
				title: featureTitle('/favorites'),
				url: '/favorites',
				icon: 'ph:clover'
			},
			{
				title: featureTitle('/recents'),
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
		<Section label="Shortcut" items={data.shortcuts} />
		<!-- <Section label="Alanysis" items={data.analysis} /> -->
		<Section label="Platform" items={data.platforms} />
		<Section label="General" items={data.general} />
	</Sidebar.Content>
	<Sidebar.Footer>
		<Sidebar.Separator />
		<Footer items={data.footers} />
	</Sidebar.Footer>
	<Sidebar.Rail />
</Sidebar.Root>
