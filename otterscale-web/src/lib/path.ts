import { m } from '$lib/paraglide/messages.js';

const createScopePath = (scope: string | undefined, subPath = '') =>
	scope === undefined ? '/' : `/scope/${scope}${subPath}`;

export interface Path {
	title: string;
	url: string;
}

// Static
export const staticPaths: Record<string, Path> = {
	// External
	documentation: {
		title: m.documentation(),
		url: 'https://otterscale.github.io'
	},
	github: {
		title: 'GitHub',
		url: 'https://github.com/otterscale/otterscale'
	},
	feedback: {
		title: m.feedback(),
		url: 'https://github.com/otterscale/otterscale/issues/new/choose'
	},
	contributors: {
		title: m.contributors(),
		url: 'https://github.com/otterscale/otterscale/graphs/contributors'
	},

	// Internal
	home: { title: m.home(), url: '/' },
	login: { title: m.login(), url: '/login' },
	setup: { title: m.setup_environment(), url: '/setup' },
	scopes: { title: m.scopes(), url: '/scopes' },
	privacyPolicy: { title: m.privacy_policy(), url: '/privacy-policy' },
	termsOfService: { title: m.terms_of_service(), url: '/terms-of-service' }
};

// Dynamic
export const dynamicPaths = {
	scope: (scope: string | undefined): Path => ({ title: m.scopes(), url: createScopePath(scope) }),
	changelog: (scope: string | undefined): Path => ({
		title: m.changelog(),
		url: createScopePath(scope, '/changelog')
	}),
	account: (scope: string | undefined): Path => ({
		title: m.account(),
		url: createScopePath(scope, '/account')
	}),
	accountSettings: (scope: string | undefined): Path => ({
		title: m.settings(),
		url: createScopePath(scope, '/account/settings')
	}),
	models: (scope: string | undefined): Path => ({
		title: m.models(),
		url: createScopePath(scope, '/models')
	}),
	modelsLLM: (scope: string | undefined): Path => ({
		title: m.llm(),
		url: createScopePath(scope, '/models/llm')
	}),
	databases: (scope: string | undefined): Path => ({
		title: m.databases(),
		url: createScopePath(scope, '/databases')
	}),
	databasesRelational: (scope: string | undefined): Path => ({
		title: m.relational(),
		url: createScopePath(scope, '/databases/relational')
	}),
	databasesNoSQL: (scope: string | undefined): Path => ({
		title: m.no_sql(),
		url: createScopePath(scope, '/databases/no-sql')
	}),
	applications: (scope: string | undefined): Path => ({
		title: m.applications(),
		url: createScopePath(scope, '/applications')
	}),
	applicationsWorkload: (scope: string | undefined): Path => ({
		title: m.workload(),
		url: createScopePath(scope, '/applications/workload')
	}),
	applicationsService: (scope: string | undefined): Path => ({
		title: m.service(),
		url: createScopePath(scope, '/applications/service')
	}),
	applicationsStore: (scope: string | undefined): Path => ({
		title: m.store(),
		url: createScopePath(scope, '/applications/store')
	}),
	storage: (scope: string | undefined): Path => ({
		title: m.storage(),
		url: createScopePath(scope, '/storage')
	}),
	storageOSD: (scope: string | undefined): Path => ({
		title: m.osd(),
		url: createScopePath(scope, '/storage/osd')
	}),
	storagePool: (scope: string | undefined): Path => ({
		title: m.pool(),
		url: createScopePath(scope, '/storage/pool')
	}),
	storageBlockDevice: (scope: string | undefined): Path => ({
		title: m.block_device(),
		url: createScopePath(scope, '/storage/block-device')
	}),
	storageFileSystem: (scope: string | undefined): Path => ({
		title: m.network_file_system(),
		url: createScopePath(scope, '/storage/file-system')
	}),
	storageObjectGateway: (scope: string | undefined): Path => ({
		title: m.object_gateway(),
		url: createScopePath(scope, '/storage/object-gateway')
	}),
	machines: (scope: string | undefined): Path => ({
		title: m.machines(),
		url: createScopePath(scope, '/machines')
	}),
	machinesMetal: (scope: string | undefined): Path => ({
		title: m.metal(),
		url: createScopePath(scope, '/machines/metal')
	}),
	virtualization: (scope: string | undefined): Path => ({
		title: m.virtualization(),
		url: createScopePath(scope, '/virtualization')
	}),
	virtualizationVirtualMachine: (scope: string | undefined): Path => ({
		title: m.virtual_machine(),
		url: createScopePath(scope, '/virtualization/virtual-machine')
	}),
	settings: (scope: string | undefined): Path => ({
		title: m.settings(),
		url: createScopePath(scope, '/settings')
	}),
	settingsNetwork: (scope: string | undefined): Path => ({
		title: m.network(),
		url: createScopePath(scope, '/settings/network')
	}),
	settingsBIST: (scope: string | undefined): Path => ({
		title: m.built_in_test(),
		url: createScopePath(scope, '/settings/built-in-self-test')
	}),
	settingsSubscription: (scope: string | undefined): Path => ({
		title: m.subscription(),
		url: createScopePath(scope, '/settings/subscription')
	}),
	setupScope: (scope: string | undefined): Path => ({
		title: m.setup_scope(),
		url: createScopePath(scope, '/setup')
	}),
	setupScopeCeph: (scope: string | undefined): Path => ({
		title: 'Ceph',
		url: createScopePath(scope, '/setup/ceph')
	}),
	setupScopeKubernetes: (scope: string | undefined): Path => ({
		title: 'Kubernetes',
		url: createScopePath(scope, '/setup/kubernetes')
	})
};

const ICON_MAP = new Map([
	['/models', 'ph:robot'],
	['/databases', 'ph:database'],
	['/applications', 'ph:compass'],
	['/storage', 'ph:hard-drives'],
	['/virtualization', 'ph:cpu'],
	['/machines', 'ph:computer-tower'],
	['/settings', 'ph:sliders-horizontal']
]);

export function urlIcon(url: string): string {
	for (const [section, icon] of ICON_MAP) {
		if (url.endsWith(section)) {
			return icon;
		}
	}
	return 'ph:circle-dashed';
}

const disabledPaths = (scope: string | undefined) => ({
	ceph: [dynamicPaths.virtualization(scope), dynamicPaths.storage(scope)],
	kube: [
		dynamicPaths.models(scope),
		dynamicPaths.databases(scope),
		dynamicPaths.applications(scope),
		dynamicPaths.virtualization(scope)
	]
});

export const pathDisabled = (
	cephName: string | undefined,
	kubeName: string | undefined,
	scope: string | undefined,
	url: string
): boolean => {
	const paths = disabledPaths(scope);
	return (
		(!cephName && paths.ceph.some((path) => path.url === url)) ||
		(!kubeName && paths.kube.some((path) => path.url === url))
	);
};

const findDynamicPath = (
	pathname: string,
	scope: string | undefined
): keyof typeof dynamicPaths | null => {
	for (const [key, pathFn] of Object.entries(dynamicPaths)) {
		const path = pathFn(scope);
		if (path.url === pathname) {
			return key as keyof typeof dynamicPaths;
		}
	}
	return null;
};

const pathBypass = (pathname: string, scope: string | undefined): boolean => {
	const bypass = ['/machines', '/settings'];
	for (const url of bypass) {
		if (pathname.includes(url)) {
			return true;
		}
	}
	return false;
};

export const getValidURL = (
	pathname: string,
	scope: string | undefined,
	cephName: string | undefined,
	kubeName: string | undefined
): string => {
	if (pathBypass(pathname, scope)) {
		return pathname.replace(/\/scope\/[^\/]+/, `/scope/${scope}`);
	}

	const homeURL = dynamicPaths.scope(scope).url;
	const currentPathKey = findDynamicPath(pathname, scope);
	if (!currentPathKey) {
		return homeURL;
	}

	const path = dynamicPaths[currentPathKey](scope);
	if (pathDisabled(cephName, kubeName, scope, path.url)) {
		return homeURL;
	}

	return path.url;
};
