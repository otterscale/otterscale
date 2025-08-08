import { m } from '$lib/paraglide/messages.js';

const createScopePath = (scope: string, subPath = '') => `/scope/${scope}${subPath}`;

export interface Path {
    title: string;
    url: string;
}

// Static
export const staticPaths: Record<string, Path> = {
    // External
    documentation: { title: m.documentation(), url: "https://otterscale.github.io" },
    github: { title: "GitHub", url: "https://github.com/otterscale/otterscale" },
    feedback: { title: m.feedback(), url: "https://github.com/otterscale/otterscale/issues/new/choose" },
    contributors: { title: m.contributors(), url: "https://github.com/otterscale/otterscale/graphs/contributors" },

    // Internal
    home: { title: m.home(), url: "/" },
    login: { title: m.login(), url: "/login" },
    setup: { title: m.setup_environment(), url: "/setup" },
    scopes: { title: m.scopes(), url: "/scopes" },
    privacyPolicy: { title: m.privacy_policy(), url: "/privacy-policy" },
    termsOfService: { title: m.terms_of_service(), url: "/terms-of-service" }
};

// Dynamic
export const dynamicPaths = {
    scope: (scope: string): Path => ({ title: m.scopes(), url: createScopePath(scope) }),
    changelog: (scope: string): Path => ({ title: m.changelog(), url: createScopePath(scope, '/changelog') }),
    account: (scope: string): Path => ({ title: m.account(), url: createScopePath(scope, '/account') }),
    accountSettings: (scope: string): Path => ({ title: m.settings(), url: createScopePath(scope, '/account/settings') }),
    models: (scope: string): Path => ({ title: m.models(), url: createScopePath(scope, '/models') }),
    modelsLLM: (scope: string): Path => ({ title: m.llm(), url: createScopePath(scope, '/models/llm') }),
    databases: (scope: string): Path => ({ title: m.databases(), url: createScopePath(scope, '/databases') }),
    databasesRelational: (scope: string): Path => ({ title: m.relational(), url: createScopePath(scope, '/databases/relational') }),
    databasesNoSQL: (scope: string): Path => ({ title: m.no_sql(), url: createScopePath(scope, '/databases/no-sql') }),
    applications: (scope: string): Path => ({ title: m.applications(), url: createScopePath(scope, '/applications') }),
    applicationsWorkload: (scope: string): Path => ({ title: m.workload(), url: createScopePath(scope, '/applications/workload') }),
    applicationsService: (scope: string): Path => ({ title: m.service(), url: createScopePath(scope, '/applications/service') }),
    applicationsStore: (scope: string): Path => ({ title: m.store(), url: createScopePath(scope, '/applications/store') }),
    storage: (scope: string): Path => ({ title: m.storage(), url: createScopePath(scope, '/storage') }),
    storageCluster: (scope: string): Path => ({ title: m.cluster(), url: createScopePath(scope, '/storage/cluster') }),
    storageBlockDevice: (scope: string): Path => ({ title: m.block_device(), url: createScopePath(scope, '/storage/block-device') }),
    storageFileSystem: (scope: string): Path => ({ title: m.file_system(), url: createScopePath(scope, '/storage/file-system') }),
    storageObjectGateway: (scope: string): Path => ({ title: m.object_gateway(), url: createScopePath(scope, '/storage/object-gateway') }),
    machines: (scope: string): Path => ({ title: m.machines(), url: createScopePath(scope, '/machines') }),
    machinesMetal: (scope: string): Path => ({ title: m.metal(), url: createScopePath(scope, '/machines/metal') }),
    machinesVirtualMachine: (scope: string): Path => ({ title: m.virtual_machine(), url: createScopePath(scope, '/machines/virtual-machine') }),
    settings: (scope: string): Path => ({ title: m.settings(), url: createScopePath(scope, '/settings') }),
    settingsSSO: (scope: string): Path => ({ title: m.sso(), url: createScopePath(scope, '/settings/sso') }),
    settingsNetwork: (scope: string): Path => ({ title: m.network(), url: createScopePath(scope, '/settings/network') }),
    settingsBIST: (scope: string): Path => ({ title: m.built_in_test(), url: createScopePath(scope, '/settings/built-in-self-test') }),
    settingsSubscription: (scope: string): Path => ({ title: m.subscription(), url: createScopePath(scope, '/settings/subscription') }),
    setupScope: (scope: string): Path => ({ title: m.setup_scope(), url: createScopePath(scope, '/setup') }),
    setupScopeCeph: (scope: string): Path => ({ title: "Ceph", url: createScopePath(scope, '/setup/ceph') }),
    setupScopeKubernetes: (scope: string): Path => ({ title: "Kubernetes", url: createScopePath(scope, '/setup/kubernetes') })
};

const ICON_MAP = new Map([
    ['/models', 'ph:robot'],
    ['/databases', 'ph:database'],
    ['/applications', 'ph:compass'],
    ['/storage', 'ph:hard-drives'],
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

const disabledPaths = (scope: string) => ({
    ceph: [
        dynamicPaths.storage(scope),
    ],
    kube: [
        dynamicPaths.models(scope),
        dynamicPaths.databases(scope),
        dynamicPaths.applications(scope),
    ]
});

export const pathDisabled = (cephName: string | undefined, kubeName: string | undefined, scope: string, url: string): boolean => {
    const paths = disabledPaths(scope);
    return (!cephName && paths.ceph.some((path) => path.url === url)) ||
        (!kubeName && paths.kube.some((path) => path.url === url));
};

export const findDynamicPath = (pathname: string, scope: string): keyof typeof dynamicPaths | null => {
    for (const [key, pathFn] of Object.entries(dynamicPaths)) {
        const path = pathFn(scope);
        if (path.url === pathname) {
            return key as keyof typeof dynamicPaths;
        }
    }
    return null;
};
