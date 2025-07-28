import {
    applicationsPath,
    applicationsServicePath,
    applicationsStorePath,
    applicationsWorkloadPath,
    machinesMetalPath,
    machinesPath,
    machinesVirtualMachinePath,
    modelsLLMPath,
    modelsPath,
    settingsPath,
    settingsNetworkPath,
    storageBlockDevicePath,
    storageClusterPath,
    storageFileSystemPath,
    storageObjectGatewayPath,
    storagePath,
    settingsBISTPath,
    databasesPath,
    databasesRelationalPath,
    databasesNoSQLPath,
    getPath,
} from '$lib/path';

export const routes = [
    {
        path: getPath(modelsPath),
        items: [getPath(modelsLLMPath)]
    },
    {
        path: getPath(databasesPath),
        items: [
            getPath(databasesRelationalPath),
            getPath(databasesNoSQLPath)
        ]
    },
    {
        path: getPath(applicationsPath),
        items: [
            getPath(applicationsWorkloadPath),
            getPath(applicationsServicePath),
            getPath(applicationsStorePath)
        ]
    },
    {
        path: getPath(storagePath),
        items: [
            getPath(storageClusterPath),
            getPath(storageBlockDevicePath),
            getPath(storageFileSystemPath),
            getPath(storageObjectGatewayPath)
        ]
    },
    {
        path: getPath(machinesPath),
        items: [
            getPath(machinesMetalPath),
            getPath(machinesVirtualMachinePath)
        ]
    },
    {
        path: getPath(settingsPath),
        items: [
            getPath(settingsNetworkPath),
            getPath(settingsBISTPath)
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