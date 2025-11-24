import { get } from 'svelte/store';

import { PremiumTier_Level } from '$lib/api/environment/v1/environment_pb';
import AdvancedTierImage from '$lib/assets/advanced-tier.jpg';
import BasicTierImage from '$lib/assets/basic-tier.jpg';
import EnterpriseTierImage from '$lib/assets/enterprise-tier.jpg';
import { m } from '$lib/paraglide/messages';
import { premiumTier } from '$lib/stores';

interface Plan {
	tier: string;
	star: boolean;
	name: string;
	description: string;
	tags: string[];
	image: string;
	disabled: boolean;
}

const plans: Plan[] = [
	{
		tier: m.standard_tier(),
		star: false,
		name: m.standard_tier_name(),
		description: m.standard_tier_description(),
		tags: [m.ceph(), m.kubernetes(), m.single_node()],
		image: BasicTierImage,
		disabled: get(premiumTier).level < PremiumTier_Level.STANDARD
	},
	{
		tier: m.premium_tier(),
		star: true,
		name: m.premium_tier_name(),
		description: m.premium_tier_description(),
		tags: [m.ceph(), m.kubernetes(), m.multi_node(), m.cluster()],
		image: AdvancedTierImage,
		disabled: get(premiumTier).level < PremiumTier_Level.PREMIUM
	},
	{
		tier: m.enterprise_tier(),
		star: true,
		name: m.enterprise_tier_name(),
		description: m.enterprise_tier_description(),
		tags: [m.ceph(), m.kubernetes(), m.multi_node(), m.cluster()],
		image: EnterpriseTierImage,
		disabled: get(premiumTier).level < PremiumTier_Level.ENTERPRISE
	}
];

export type { Plan };
export { plans };
