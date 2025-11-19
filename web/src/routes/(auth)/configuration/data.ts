import { resolve } from '$app/paths';
import { m } from '$lib/paraglide/messages';

const items = [
	{
		icon: 'ph:clock',
		title: m.ntp_server(),
		description: m.setting_ntp_server_description(),
		type: m.machine(),
		url: resolve('/(auth)/configuration/ntp-server'),
		default: true
	},
	{
		icon: 'ph:disc',
		title: m.boot_image(),
		description: m.setting_boot_image_description(),
		type: m.machine(),
		url: resolve('/(auth)/configuration/boot-image')
	},
	{
		icon: 'ph:tag-simple',
		title: m.machine_tag(),
		description: m.setting_machine_tag_description(),
		type: m.machine(),
		url: resolve('/(auth)/configuration/machine-tag')
	},
	{
		icon: 'ph:package',
		title: m.package_repository(),
		description: m.setting_package_repository_description(),
		type: m.machine(),
		url: resolve('/(auth)/configuration/package-repository')
	},
	{
		icon: 'ph:cube',
		title: m.helm_repository(),
		description: m.setting_helm_repository_description(),
		type: m.application(),
		url: resolve('/(auth)/configuration/helm-repository')
	},
	{
		icon: 'ph:wallet',
		title: m.subscription(),
		description: m.settings(),
		url: resolve('/(auth)/configuration/subscription')
	}
];

export { items };
