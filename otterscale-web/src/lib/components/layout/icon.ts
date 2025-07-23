import { applicationsPath, changelogPath, documentationPath, feedbackPath, machinesPath, modelPath, settingsPath, storagePath } from "$lib/path";

export function getIconFromUrl(url: string): string {
    switch (url) {
        case changelogPath:
            return 'ph:clock-counter-clockwise'
        case documentationPath:
            return 'ph:book-open'
        case feedbackPath:
            return 'ph:paper-plane-tilt'
    }

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