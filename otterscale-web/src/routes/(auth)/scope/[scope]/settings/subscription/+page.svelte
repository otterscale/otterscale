<script lang="ts">
	import Icon from '@iconify/svelte';
	import { page } from '$app/state';
	import { PremiumTier } from '$lib/api/premium/v1/premium_pb';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb, premiumTier } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.settings(page.params.scope)],
		current: dynamicPaths.settingsSubscription(page.params.scope)
	});

	// Tier configurations
	const tiers = [
		{
			id: PremiumTier.BASIC,
			title: m.basic_tier(),
			description: m.basic_setup_description(),
			icons: [
				{ name: 'simple-icons:ceph', color: '#f0424d', hasAnimation: false },
				{ name: 'simple-icons:kubernetes', color: '#326de6', hasAnimation: false }
			],
			features: [
				m.single_node_ceph(),
				m.single_node_ceph_description(),
				m.single_node_kubernetes(),
				m.single_node_kubernetes_description(),
				m.vm_support(),
				m.vm_support_description(),
				m.app_store_access(),
				m.app_store_access_description(),
				m.open_source_code(),
				m.open_source_code_description(),
				m.community_support(),
				m.community_support_description(),
				m.usage_limitations(),
				m.usage_limitations_description()
			],
			isRecommended: false
		},
		{
			id: PremiumTier.ADVANCED,
			title: m.advanced_tier(),
			description: m.advanced_setup_description(),
			icons: [{ name: 'simple-icons:ceph', color: '#f0424d', hasAnimation: true }],
			features: [
				m.all_basic_features(),
				m.all_basic_features_description(),
				m.multi_node_ceph(),
				m.multi_node_ceph_description(),
				m.team_access(),
				m.team_access_description(),
				m.priority_support(),
				m.priority_support_description()
			],
			isRecommended: false
		},
		{
			id: PremiumTier.ENTERPRISE,
			title: m.enterprise_tier(),
			description: m.enterprise_setup_description(),
			icons: [
				{ name: 'simple-icons:ceph', color: '#f0424d', hasAnimation: true },
				{ name: 'simple-icons:kubernetes', color: '#326de6', hasAnimation: true }
			],
			features: [
				m.all_advanced_features(),
				m.all_advanced_features_description(),
				m.enterprise_multi_node_ceph(),
				m.enterprise_multi_node_ceph_description(),
				m.multi_node_kubernetes(),
				m.multi_node_kubernetes_description(),
				m.enterprise_capabilities(),
				m.enterprise_capabilities_description(),
				m.dedicated_support(),
				m.dedicated_support_description(),
				m.customization_services(),
				m.customization_services_description(),
				m.training_services(),
				m.training_services_description()
			],
			isRecommended: true
		}
	];
</script>

<!-- just-in-time  -->
<dummy class="bg-[#326de6]"></dummy>
<dummy class="bg-[#f0424d]"></dummy>

<div class="mx-auto max-w-7xl px-4 xl:px-0">
	<div class="grid gap-12 md:grid-cols-2 lg:grid-cols-3 lg:items-start">
		{#each tiers as tier}
			<Card.Root class="flex h-full flex-col {tier.isRecommended ? 'border-primary relative' : ''}">
				{#if tier.isRecommended}
					<div
						class="bg-primary text-primary-foreground absolute top-0 right-0 rounded-tr-lg rounded-bl-lg px-3 py-1 text-xs font-medium uppercase"
					>
						{m.recommended()}
					</div>
				{/if}

				<Card.Header class="justify-center pb-2 text-center">
					<Card.Title class="mb-7">{tier.title}</Card.Title>
					<div class="flex {tier.icons.length > 1 ? 'space-x-2' : ''}">
						{#each tier.icons as iconConfig}
							<div class="flex">
								<Icon icon={iconConfig.name} class="size-12 text-[{iconConfig.color}]" />
								{#if iconConfig.hasAnimation}
									<span class="relative flex size-3">
										<span
											class="absolute inline-flex h-full w-full animate-ping rounded-full bg-[{iconConfig.color}] opacity-75"
										></span>
										<span class="relative inline-flex size-3 rounded-full bg-[{iconConfig.color}]"
										></span>
									</span>
								{/if}
							</div>
						{/each}
					</div>
				</Card.Header>

				<Card.Description class="text-center {tier.isRecommended ? 'mx-auto w-11/12' : ''}">
					{tier.description}
				</Card.Description>

				<Card.Content>
					<ul>
						{#each tier.features as feature, index}
							{#if index % 2 === 0}
								<li class="flex space-x-2 pt-4">
									<Icon icon="ph:check-bold" class="mt-0.5 h-4 w-4 flex-shrink-0 text-blue-500" />
									<span class="text-foreground text-sm">{feature}</span>
								</li>
							{:else}
								<li class="flex space-x-2 py-1 pl-6">
									<span class="text-muted-foreground text-xs">{feature}</span>
								</li>
							{/if}
						{/each}
					</ul>
				</Card.Content>

				<Card.Footer class="mt-auto">
					<Button
						href="mailto:paul_tsai@phison.com"
						class="w-full"
						variant={tier.isRecommended ? 'default' : 'outline'}
						disabled={tier.id <= $premiumTier}
					>
						{#if tier.id > $premiumTier}
							{m.contact_sales()}
						{:else}
							{m.enabled()}
						{/if}
					</Button>
				</Card.Footer>
			</Card.Root>
		{/each}
	</div>
</div>
