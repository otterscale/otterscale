import AdvancedTierImage from '$lib/assets/advanced-tier.jpg';
import BasicTierImage from '$lib/assets/basic-tier.jpg';
import EnterpriseTierImage from '$lib/assets/enterprise-tier.jpg';
import { m } from '$lib/paraglide/messages';

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
		tier: m.community_tier(),
		star: false,
		name: m.community_tier_name(),
		description: m.community_tier_description(),
		tags: [m.ceph(), m.kubernetes(), m.single_node()],
		image: BasicTierImage,
		disabled: false
	},
	// TODO: Standard tier disabled until feature is available
	{
		tier: m.premium_tier(),
		star: true,
		name: m.premium_tier_name(),
		description: m.premium_tier_description(),
		tags: [m.ceph(), m.kubernetes(), m.multi_node(), m.cluster()],
		image: AdvancedTierImage,
		disabled: true
	},
	{
		tier: m.enterprise_tier(),
		star: true,
		name: m.enterprise_tier_name(),
		description: m.enterprise_tier_description(),
		tags: [m.ceph(), m.kubernetes(), m.multi_node(), m.cluster()],
		image: EnterpriseTierImage,
		disabled: true
	}
];

export type { Plan };
export { plans };
