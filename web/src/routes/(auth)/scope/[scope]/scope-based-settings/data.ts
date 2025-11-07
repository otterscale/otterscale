import { resolve } from '$app/paths';
import { page } from '$app/state';
import { m } from '$lib/paraglide/messages';

const items = [
	{
		icon: 'ph:cube',
		title: m.extensions(),
		url: resolve('/(auth)/scope/[scope]/scope-based-settings/extensions', { scope: page.params.scope! }),
		default: true,
	},
	{
		icon: 'ph:hard-drives',
		title: m.data_volume(),
		type: m.virtual_machine(),
		url: resolve('/(auth)/scope/[scope]/scope-based-settings/data-volume', { scope: page.params.scope! }),
	},
	{
		icon: 'ph:cpu',
		title: m.instance_type(),
		type: m.virtual_machine(),
		url: resolve('/(auth)/scope/[scope]/scope-based-settings/instance-type', { scope: page.params.scope! }),
	},
];

export { items };
