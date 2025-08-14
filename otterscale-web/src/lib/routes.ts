import { dynamicPaths, type Path } from '$lib/path';

export interface Route {
    path: Path;
    items: Path[];
};

export const routes = (scope: string): Route[] => [
    {
        path: dynamicPaths.models(scope),
        items: [dynamicPaths.modelsLLM(scope)]
    },
    {
        path: dynamicPaths.databases(scope),
        items: [
            dynamicPaths.databasesRelational(scope),
            dynamicPaths.databasesNoSQL(scope)
        ]
    },
    {
        path: dynamicPaths.applications(scope),
        items: [
            dynamicPaths.applicationsWorkload(scope),
            dynamicPaths.applicationsService(scope),
            dynamicPaths.applicationsStore(scope)
        ]
    },
    {
        path: dynamicPaths.storage(scope),
        items: [
            dynamicPaths.storageOSD(scope),
            dynamicPaths.storagePool(scope),
            dynamicPaths.storageBlockDevice(scope),
            dynamicPaths.storageFileSystem(scope),
            dynamicPaths.storageObjectGateway(scope)
        ]
    },
    {
        path: dynamicPaths.machines(scope),
        items: [
            dynamicPaths.machinesMetal(scope),
            dynamicPaths.machinesVirtualMachine(scope)
        ]
    },
    {
        path: dynamicPaths.settings(scope),
        items: [
            dynamicPaths.settingsNetwork(scope),
            dynamicPaths.settingsSubscription(scope),
            dynamicPaths.settingsBIST(scope),
        ]
    }
];

export const bookmarks = [
    { title: 'FOO 1', url: '#' },
    { title: 'BAR 1', url: '#' },
    { title: 'FOO 2', url: '#' },
    { title: 'BAR 2', url: '#' },
    { title: 'FOO 3', url: '#' },
    { title: 'BAR 3', url: '#' }
];

