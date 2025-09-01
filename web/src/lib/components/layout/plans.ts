import { get } from 'svelte/store';
import { PremiumTier } from '$lib/api/premium/v1/premium_pb';
import AdvancedTierImage from '$lib/assets/advanced-tier.jpg';
import BasicTierImage from '$lib/assets/basic-tier.jpg';
import EnterpriseTierImage from '$lib/assets/enterprise-tier.jpg';
import { m } from '$lib/paraglide/messages';
import { premiumTier } from '$lib/stores';

export interface Plan {
	tier: string;
	star: boolean;
	name: string;
	description: string;
	tags: string[];
	image: string;
	disabled: boolean;
}

export const plans: Plan[] = [
	{
		tier: m.basic_tier(),
		star: false,
		name: m.basic_tier_name(),
		description: m.basic_tier_description(),
		tags: ['Ceph', 'Kubernetes', m.single_node()],
		image: BasicTierImage,
		disabled: get(premiumTier) < PremiumTier.BASIC,
	},
	{
		tier: m.advanced_tier(),
		star: true,
		name: m.advanced_tier_name(),
		description: m.advanced_tier_description(),
		tags: ['Ceph', 'Multi-Node', m.multi_node(), m.cluster()],
		image: AdvancedTierImage,
		disabled: get(premiumTier) < PremiumTier.ADVANCED,
	},
	{
		tier: m.enterprise_tier(),
		star: true,
		name: m.enterprise_tier_name(),
		description: m.enterprise_tier_description(),
		tags: ['Ceph', 'Kubernetes', m.multi_node(), m.cluster()],
		image: EnterpriseTierImage,
		disabled: get(premiumTier) < PremiumTier.ENTERPRISE,
	},
];
