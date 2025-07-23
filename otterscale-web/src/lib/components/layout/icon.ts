import { applicationsPath, documentationPath, feedbackPath, machinesPath, modelPath, settingsPath, storagePath } from "$lib/path";

export function getIconFromUrl(url: string): string {
    if (url.startsWith(modelPath)) {
        return "ph:robot"
    } else if (url.startsWith(applicationsPath)) {
        return 'ph:compass';
    } else if (url.startsWith(storagePath)) {
        return "ph:hard-drives"
    } else if (url.startsWith(machinesPath)) {
        return "ph:computer-tower"
    } else if (url.startsWith(settingsPath)) {
        return "ph:sliders-horizontal"
    }
    return 'ph:circle-dashed';
}