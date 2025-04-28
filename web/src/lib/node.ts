export function nodeIcon(connector: string) {
    switch (connector) {
        case 'kubernetes': return 'logos:kubernetes';
        case 'Ceph': return 'simple-icons:ceph';
        case 'JUJU': return 'logos:juju';
        case 'MAAS': return 'simple-icons:maas';
        default: return 'ph:circle-dashed';
    }
}

export function nodeLabel(connector: string) {
    switch (connector) {
        case 'kubernetes': return 'ph:cloud';
        case 'Ceph': return 'ph:cloud';
        case 'JUJU': return 'ph:cloud';
        case 'MAAS': return 'ph:cloud';
        default: return 'unknown';
    }
}

export function nodeLabelIcon(connector: string) {
    switch (connector) {
        case 'kubernetes': return 'ph:cloud';
        case 'Ceph': return 'ph:cloud';
        case 'JUJU': return 'ph:cloud';
        case 'MAAS': return 'ph:cloud';
        default: return 'ph:circle-dashed';
    }
}