export function connectorIcon(connector: string) {
    switch (connector) {
        case 'postgresql': return 'logos:postgresql';
        case 'csv': return 'ph:file-csv';
        default: return 'ph:circle-dashed';
    }
}

export function connectorLabel(connector: string) {
    switch (connector) {
        case 'postgresql': return 'database';
        case 'csv': return 'file';
        default: return 'unknown';
    }
}

export function connectorLabelIcon(connector: string) {
    switch (connector) {
        case 'postgresql': return 'ph:database';
        case 'csv': return 'ph:folder';
        default: return 'ph:circle-dashed';
    }
}