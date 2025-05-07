import { iconExists } from "@iconify/svelte";

export function fuzzLogosIcon(name: string, defaultName: string): string {
    if (!name) {
        return defaultName
    }
    let icon = `logos:${name}-icon`
    if (iconExists(icon)) {
        return icon
    }
    icon = `logos:${name}`
    if (iconExists(icon)) {
        return icon
    }
    return defaultName;
}