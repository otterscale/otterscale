import { page } from '$app/state';
import { m } from '$lib/paraglide/messages';
import { dynamicPaths } from '$lib/path';

const items = [
	{
		icon: 'ph:clock',
		title: m.ntp_server(),
		description: m.setting_ntp_server_description(),
		type: m.machine(),
		url: dynamicPaths.settingsNTPServer(page.params.scope).url,
		default: true,
	},
	{
		icon: 'ph:disc',
		title: m.boot_image(),
		description: m.setting_boot_image_description(),
		type: m.machine(),
		url: dynamicPaths.settingsBootImage(page.params.scope).url,
	},
	{
		icon: 'ph:tag-simple',
		title: m.machine_tag(),
		description: m.setting_machine_tag_description(),
		type: m.machine(),
		url: dynamicPaths.settingsMachineTag(page.params.scope).url,
	},
	{
		icon: 'ph:package',
		title: m.package_repository(),
		description: m.setting_package_repository_description(),
		type: m.machine(),
		url: dynamicPaths.settingsPackageRepository(page.params.scope).url,
	},
	{
		icon: 'ph:cube',
		title: m.helm_repository(),
		description: m.setting_helm_repository_description(),
		type: m.application(),
		url: dynamicPaths.settingsHelmRepository(page.params.scope).url,
	},
	{
		icon: 'ph:test-tube',
		title: m.built_in_test(),
		description: m.settings(),
		url: dynamicPaths.settingsBuiltInTest(page.params.scope).url,
	},
	{
		icon: 'ph:key',
		title: m.sso(),
		description: m.setting_single_sign_on_description(),
		url: dynamicPaths.settingsSSO(page.params.scope).url,
	},
	{
		icon: 'ph:wallet',
		title: m.subscription(),
		description: m.settings(),
		url: dynamicPaths.settingsSubscription(page.params.scope).url,
	},
];

export { items };
