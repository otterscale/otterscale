import { resolve } from '$app/paths';
import { page } from '$app/state';
import { m } from '$lib/paraglide/messages';

const items = [
	{
		icon: 'ph:clock',
		title: m.ntp_server(),
		description: m.setting_ntp_server_description(),
		type: m.machine(),
		url: resolve('/(auth)/scope/[scope]/global-settings/ntp-server', { scope: page.params.scope! }),
		default: true,
	},
	{
		icon: 'ph:disc',
		title: m.boot_image(),
		description: m.setting_boot_image_description(),
		type: m.machine(),
		url: resolve('/(auth)/scope/[scope]/global-settings/boot-image', { scope: page.params.scope! }),
	},
	{
		icon: 'ph:tag-simple',
		title: m.machine_tag(),
		description: m.setting_machine_tag_description(),
		type: m.machine(),
		url: resolve('/(auth)/scope/[scope]/global-settings/machine-tag', { scope: page.params.scope! }),
	},
	{
		icon: 'ph:package',
		title: m.package_repository(),
		description: m.setting_package_repository_description(),
		type: m.machine(),
		url: resolve('/(auth)/scope/[scope]/global-settings/package-repository', { scope: page.params.scope! }),
	},
	{
		icon: 'ph:cube',
		title: m.helm_repository(),
		description: m.setting_helm_repository_description(),
		type: m.application(),
		url: resolve('/(auth)/scope/[scope]/global-settings/helm-repository', { scope: page.params.scope! }),
	},
	{
		icon: 'ph:test-tube',
		title: m.built_in_test(),
		description: m.settings(),
		url: resolve('/(auth)/scope/[scope]/global-settings/built-in-test', { scope: page.params.scope! }),
	},
	{
		icon: 'ph:key',
		title: m.single_sign_on(),
		description: m.setting_single_sign_on_description(),
		url: resolve('/(auth)/scope/[scope]/global-settings/single-sign-on', { scope: page.params.scope! }),
	},
	{
		icon: 'ph:wallet',
		title: m.subscription(),
		description: m.settings(),
		url: resolve('/(auth)/scope/[scope]/global-settings/subscription', { scope: page.params.scope! }),
	},
];

export { items };
