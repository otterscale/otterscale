import { page } from '$app/state';
import { m } from '$lib/paraglide/messages';
import { dynamicPaths } from '$lib/path';

const items = [
	{
		icon: 'ph:cube',
		title: m.extensions(),
		url: dynamicPaths.settingsExtensions(page.params.scope).url,
		default: true,
	},
	{
		icon: 'ph:hard-drives',
		title: m.data_volume(),
		type: m.virtual_machine(),
		url: dynamicPaths.settingsDataVolume(page.params.scope).url,
	},
	{
		icon: 'ph:cpu',
		title: m.instance_type(),
		type: m.virtual_machine(),
		url: dynamicPaths.settingsInstanceType(page.params.scope).url,
	},
];

export { items };
