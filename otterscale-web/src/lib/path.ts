import { m } from '$lib/paraglide/messages.js';

// External
export const documentationPath = "https://otterscale.github.io";
export const githubPath = "https://github.com/otterscale/otterscale";
export const feedbackPath = "https://github.com/otterscale/otterscale/issues/new/choose";
export const contributorsPath = "https://github.com/otterscale/otterscale/graphs/contributors";

// Misc
export const homePath = "/";
export const privacyPolicyPath = "/privacy-policy";
export const termsOfServicePath = "/terms-of-service";

// Account
export const loginPath = "/login";
export const accountPath = "/account"
export const accountSettingsPath = "/account/settings"

// Model
export const modelsPath = "/models"
export const modelsLLMPath = "/models/llm"

// Database
export const databasesPath = "/databases"
export const databasesRelationalPath = "/databases/relational"
export const databasesNoSQLPath = "/databases/no-sql"

// Application
export const applicationsPath = "/applications"
export const applicationsWorkloadPath = "/applications/workload"
export const applicationsServicePath = "/applications/service"
export const applicationsStorePath = "/applications/store"

// Storage
export const storagePath = "/storage"
export const storageClusterPath = "/storage/cluster"
export const storageBlockDevicePath = "/storage/block-device"
export const storageFileSystemPath = "/storage/file-system"
export const storageObjectGatewayPath = "/storage/object-gateway"

// Machine
export const machinesPath = "/machines"
export const machinesMetalPath = "/machines/metal"
export const machinesVirtualMachinePath = "/machines/virtual-machine"

// Setting
export const settingsPath = "/settings"
export const settingsNetworkPath = "/settings/network"
export const settingsBISTPath = "/settings/built-in-self-test"

// Functions
export function getIconFromUrl(url: string): string {
    const iconMap = new Map<string, string>([
        [modelsPath, "ph:robot"],
        [databasesPath, "ph:database"],
        [applicationsPath, "ph:compass"],
        [storagePath, "ph:hard-drives"],
        [machinesPath, "ph:computer-tower"],
        [settingsPath, "ph:sliders-horizontal"]
    ]);

    for (const [path, icon] of iconMap) {
        if (url.startsWith(path)) {
            return icon;
        }
    }
    return 'ph:circle-dashed';
}

export interface Path {
    title: string;
    url: string;
}

const routesMap = new Map<string, Path>([
    [homePath, { title: m.home(), url: homePath }],
    [accountPath, { title: m.account(), url: accountPath }],
    [modelsPath, { title: m.models(), url: modelsPath }],
    [modelsLLMPath, { title: m.llm(), url: modelsLLMPath }],
    [databasesPath, { title: m.databases(), url: databasesPath }],
    [databasesRelationalPath, { title: m.relational(), url: databasesRelationalPath }],
    [databasesNoSQLPath, { title: m.no_sql(), url: databasesNoSQLPath }],
    [applicationsPath, { title: m.applications(), url: applicationsPath }],
    [applicationsWorkloadPath, { title: m.workload(), url: applicationsWorkloadPath }],
    [applicationsServicePath, { title: m.service(), url: applicationsServicePath }],
    [applicationsStorePath, { title: m.store(), url: applicationsStorePath }],
    [storagePath, { title: m.storage(), url: storagePath }],
    [storageClusterPath, { title: m.cluster(), url: storageClusterPath }],
    [storageBlockDevicePath, { title: m.block_device(), url: storageBlockDevicePath }],
    [storageFileSystemPath, { title: m.file_system(), url: storageFileSystemPath }],
    [storageObjectGatewayPath, { title: m.object_gateway(), url: storageObjectGatewayPath }],
    [machinesPath, { title: m.machines(), url: machinesPath }],
    [machinesMetalPath, { title: m.metal(), url: machinesMetalPath }],
    [machinesVirtualMachinePath, { title: m.virtual_machine(), url: machinesVirtualMachinePath }],
    [settingsPath, { title: m.settings(), url: settingsPath }],
    [settingsNetworkPath, { title: m.network(), url: settingsNetworkPath }],
    [settingsBISTPath, { title: m.built_in_test(), url: settingsBISTPath }],
]);

export function getPath(url: string): Path {
    return routesMap.get(url) ?? { title: m.home(), url: homePath };
}
