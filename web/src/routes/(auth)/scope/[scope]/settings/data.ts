import { resolve } from '$app/paths';
import { m } from '$lib/paraglide/messages';

export const getItems = (scope: string) => [
	{
		icon: 'ph:cube',
		title: m.extensions(),
		url: resolve('/(auth)/scope/[scope]/settings/extensions', { scope }),
		default: true
	},
	{
		icon: 'ph:test-tube',
		title: m.built_in_test(),
		description: m.configuration(),
		url: resolve('/(auth)/scope/[scope]/settings/built-in-test', { scope })
	},
	{
		icon: 'ph:hard-drives',
		title: m.data_volume(),
		type: m.virtual_machine(),
		url: resolve('/(auth)/scope/[scope]/settings/data-volume', { scope })
	},
	{
		icon: 'ph:cpu',
		title: m.instance_type(),
		type: m.virtual_machine(),
		url: resolve('/(auth)/scope/[scope]/settings/instance-type', { scope })
	},
	{
		icon: 'ph:robot',
		title: m.model_artifact(),
		type: m.model(),
		url: resolve('/(auth)/scope/[scope]/settings/model-artifact', { scope })
	}
];
