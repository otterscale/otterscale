<script lang="ts" module>
	import * as Collapsible from '$lib/components/ui/collapsible';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { siteConfig } from '$lib/config/site';
	import * as m from '$lib/paraglide/messages.js';
	import Icon from '@iconify/svelte';
	import type { ComponentProps } from 'svelte';
	import { featureTitle } from '../features';
	import Footer from './footer.svelte';
	import Namespace from './namespace.svelte';
	import Section from './section.svelte';

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
				title: m.dashboard(),
				url: '/dashboard?intervals=30',
				icon: 'ph:gauge'
			},
			{
				title: m.orchestration(),
				url: '/orchestration',
				icon: 'ph:tree-structure'
			},
			{
				title: m.management(),
				url: '#',
				icon: 'ph:command',
				isActive: false,
				items: [
					{
						title: m.model(),
						url: '/management/llm'
					},
					{
						title: m.application(),
						url: '/management/application'
					},
					{
						title: m.facility(),
						url: '/management/facility'
					},

					{
						title: m.machine(),
						url: '/management/machine'
					},
					{
						title: m.network(),
						url: '/management/network'
					},
					{
						title: m.configuration(),
						url: '/management/general'
					}
				]
			},
			{
				title: m.store(),
				url: '/market',
				icon: 'ph:magnifying-glass'
			}
		],
		general: [
			{
				title: m.home(),
				url: '/',
				icon: 'ph:house'
			},
			{
				title: m.settings(),
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
				title: m.support(),
				url: 'https://openhdc.github.io',
				icon: 'ph:lifebuoy'
			},
			{
				title: siteConfig.version,
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
</script>

<script lang="ts">
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
		<Section label={m.general()} items={data.general} />
		<Section label={m.shortcut()} items={data.shortcuts} />

		<Sidebar.Group class="group-data-[collapsible=icon]:hidden">
			<Sidebar.GroupLabel>Platform</Sidebar.GroupLabel>

			<Sidebar.Menu>
				<Sidebar.MenuItem>
					<Sidebar.MenuButton>
						{#snippet child({ props })}
							<a href="/management/llm" {...props}>
								<Icon icon="ph:robot" />
								{m.model()}
							</a>
						{/snippet}
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>

				<Sidebar.MenuItem>
					<Sidebar.MenuButton>
						{#snippet child({ props })}
							<a href="/market" {...props}>
								<Icon icon="ph:magnifying-glass" />
								{m.store()}
							</a>
						{/snippet}
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>

				<Collapsible.Root class="group/main-collapsible">
					{#snippet child({ props })}
						<Sidebar.MenuItem {...props}>
							<Collapsible.Trigger>
								{#snippet child({ props })}
									<Sidebar.MenuButton {...props}>
										<Icon icon="ph:hard-drives" />
										Storage
										<Icon
											icon="ph:caret-right"
											class="ml-auto transition-transform duration-200 group-data-[state=open]/main-collapsible:rotate-90"
										/>
									</Sidebar.MenuButton>
								{/snippet}
							</Collapsible.Trigger>
							<Collapsible.Content>
								<Sidebar.MenuSub class="mr-0 pr-0">
									<Collapsible.Root class="group/sub-collapsible">
										<Collapsible.Trigger
											>{#snippet child({ props })}
												<Sidebar.MenuButton {...props}>
													<Icon icon="ph:cube" />
													Cluster
													<Icon
														icon="ph:caret-right"
														class="ml-auto transition-transform duration-200 group-data-[state=open]/sub-collapsible:rotate-90"
													/>
												</Sidebar.MenuButton>
											{/snippet}
										</Collapsible.Trigger>
										<Collapsible.Content>
											<Sidebar.MenuSub>
												<Sidebar.MenuSubItem>
													<Sidebar.MenuSubButton>
														{#snippet child({ props })}
															<a href="/management/storage/cluster_monitor" {...props}>Monitor</a>
														{/snippet}
													</Sidebar.MenuSubButton>
												</Sidebar.MenuSubItem>
												<Sidebar.MenuSubItem>
													<Sidebar.MenuSubButton>
														{#snippet child({ props })}
															<a href="/management/storage/cluster_pool" {...props}>Pool</a>
														{/snippet}
													</Sidebar.MenuSubButton>
												</Sidebar.MenuSubItem>
												<Sidebar.MenuSubItem>
													<Sidebar.MenuSubButton>
														{#snippet child({ props })}
															<a
																href="/management/storage/cluster_object_storage_daemon"
																{...props}
															>
																OSD
															</a>
														{/snippet}
													</Sidebar.MenuSubButton>
												</Sidebar.MenuSubItem>
											</Sidebar.MenuSub>
										</Collapsible.Content>
									</Collapsible.Root>

									<Collapsible.Root class="group/sub-collapsible">
										<Collapsible.Trigger>
											{#snippet child({ props })}
												<Sidebar.MenuButton {...props}>
													<Icon icon="ph:cube" />
													Block
													<Icon
														icon="ph:caret-right"
														class="ml-auto transition-transform duration-200 group-data-[state=open]/sub-collapsible:rotate-90"
													/>
												</Sidebar.MenuButton>
											{/snippet}
										</Collapsible.Trigger>
										<Collapsible.Content>
											<Sidebar.MenuSub>
												<Sidebar.MenuSubItem>
													<Sidebar.MenuSubButton>
														{#snippet child({ props })}
															<a href="/management/storage/block" {...props}>Image</a>
														{/snippet}
													</Sidebar.MenuSubButton>
												</Sidebar.MenuSubItem>
											</Sidebar.MenuSub>
										</Collapsible.Content>
									</Collapsible.Root>

									<Collapsible.Root class="group/sub-collapsible">
										<Collapsible.Trigger
											>{#snippet child({ props })}
												<Sidebar.MenuButton {...props}>
													<Icon icon="ph:cube" />
													File
													<Icon
														icon="ph:caret-right"
														class="ml-auto transition-transform duration-200 group-data-[state=open]/sub-collapsible:rotate-90"
													/>
												</Sidebar.MenuButton>
											{/snippet}
										</Collapsible.Trigger>
										<Collapsible.Content>
											<Sidebar.MenuSub>
												<Sidebar.MenuSubItem>
													<Sidebar.MenuSubButton>
														{#snippet child({ props })}
															<a href="/management/storage/file_system" {...props}>File System</a>
														{/snippet}
													</Sidebar.MenuSubButton>
												</Sidebar.MenuSubItem>
											</Sidebar.MenuSub>
										</Collapsible.Content>
									</Collapsible.Root>

									<Collapsible.Root class="group/sub-collapsible">
										<Collapsible.Trigger>
											{#snippet child({ props })}
												<Sidebar.MenuButton {...props}>
													<Icon icon="ph:cube" />
													Object
													<Icon
														icon="ph:caret-right"
														class="ml-auto transition-transform duration-200 group-data-[state=open]/sub-collapsible:rotate-90"
													/>
												</Sidebar.MenuButton>
											{/snippet}
										</Collapsible.Trigger>
										<Collapsible.Content>
											<Sidebar.MenuSub>
												<Sidebar.MenuSubItem>
													<Sidebar.MenuSubButton>
														{#snippet child({ props })}
															<a href="/management/storage/object_bucket" {...props}>Bucket</a>
														{/snippet}
													</Sidebar.MenuSubButton>
												</Sidebar.MenuSubItem>
												<Sidebar.MenuSubItem>
													<Sidebar.MenuSubButton>
														{#snippet child({ props })}
															<a href="/management/storage/object_user" {...props}>User</a>
														{/snippet}
													</Sidebar.MenuSubButton>
												</Sidebar.MenuSubItem>
											</Sidebar.MenuSub>
										</Collapsible.Content>
									</Collapsible.Root>
								</Sidebar.MenuSub>
							</Collapsible.Content>
						</Sidebar.MenuItem>
					{/snippet}
				</Collapsible.Root>

				<Collapsible.Root class="group/main-collapsible">
					{#snippet child({ props })}
						<Sidebar.MenuItem {...props}>
							<Collapsible.Trigger>
								{#snippet child({ props })}
									<Sidebar.MenuButton {...props}>
										<Icon icon="ph:speedometer" />
										Dashboard
										<Icon
											icon="ph:caret-right"
											class="ml-auto transition-transform duration-200 group-data-[state=open]/main-collapsible:rotate-90"
										/>
									</Sidebar.MenuButton>
								{/snippet}
							</Collapsible.Trigger>
							<Collapsible.Content>
								<Sidebar.MenuSub>
									<Sidebar.MenuSubItem>
										<Sidebar.MenuSubItem>
											<Sidebar.MenuSubButton>
												{#snippet child({ props })}
													<a href="/dashboard/overall" {...props}>Overall</a>
												{/snippet}
											</Sidebar.MenuSubButton>
										</Sidebar.MenuSubItem>

										<Sidebar.MenuSubButton>
											{#snippet child({ props })}
												<a href="/dashboard/application" {...props}>Application</a>
											{/snippet}
										</Sidebar.MenuSubButton>
									</Sidebar.MenuSubItem>

									<Sidebar.MenuSubItem>
										<Sidebar.MenuSubButton>
											{#snippet child({ props })}
												<a href="/dashboard/storage" {...props}>Storage</a>
											{/snippet}
										</Sidebar.MenuSubButton>
									</Sidebar.MenuSubItem>

									<Sidebar.MenuSubItem>
										<Sidebar.MenuSubButton>
											{#snippet child({ props })}
												<a href="/dashboard/hardware" {...props}>Hardware</a>
											{/snippet}
										</Sidebar.MenuSubButton>
									</Sidebar.MenuSubItem>
								</Sidebar.MenuSub>
							</Collapsible.Content>
						</Sidebar.MenuItem>
					{/snippet}
				</Collapsible.Root>

				<Collapsible.Root class="group/main-collapsible">
					{#snippet child({ props })}
						<Sidebar.MenuItem {...props}>
							<Collapsible.Trigger>
								{#snippet child({ props })}
									<Sidebar.MenuButton {...props}>
										<Icon icon="ph:gear-six" />
										Management
										<Icon
											icon="ph:caret-right"
											class="ml-auto transition-transform duration-200 group-data-[state=open]/main-collapsible:rotate-90"
										/>
									</Sidebar.MenuButton>
								{/snippet}
							</Collapsible.Trigger>
							<Collapsible.Content>
								<Sidebar.MenuSub>
									<Sidebar.MenuSubItem>
										<Sidebar.MenuSubButton>
											{#snippet child({ props })}
												<a href="/management/application" {...props}>
													{m.application()}
												</a>
											{/snippet}
										</Sidebar.MenuSubButton>
									</Sidebar.MenuSubItem>

									<Sidebar.MenuSubItem>
										<Sidebar.MenuSubButton>
											{#snippet child({ props })}
												<a href="/management/facility" {...props}>
													{m.facility()}
												</a>
											{/snippet}
										</Sidebar.MenuSubButton>
									</Sidebar.MenuSubItem>

									<Sidebar.MenuSubItem>
										<Sidebar.MenuSubButton>
											{#snippet child({ props })}
												<a href="/management/machine" {...props}>
													{m.machine()}
												</a>
											{/snippet}
										</Sidebar.MenuSubButton>
									</Sidebar.MenuSubItem>

									<Sidebar.MenuSubItem>
										<Sidebar.MenuSubButton>
											{#snippet child({ props })}
												<a href="/management/network" {...props}>
													{m.network()}
												</a>
											{/snippet}
										</Sidebar.MenuSubButton>
									</Sidebar.MenuSubItem>

									<Sidebar.MenuSubItem>
										<Sidebar.MenuSubButton>
											{#snippet child({ props })}
												<a href="/management/general" {...props}>
													{m.configuration()}
												</a>
											{/snippet}
										</Sidebar.MenuSubButton>
									</Sidebar.MenuSubItem>
								</Sidebar.MenuSub>
							</Collapsible.Content>
						</Sidebar.MenuItem>
					{/snippet}
				</Collapsible.Root>
			</Sidebar.Menu>
		</Sidebar.Group>
	</Sidebar.Content>
	<Sidebar.Separator />
	<Sidebar.Footer>
		<Footer items={data.footers} />
	</Sidebar.Footer>
</Sidebar.Root>
