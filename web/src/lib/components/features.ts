import * as m from '$lib/paraglide/messages.js';

export interface Feature {
    path: string;
    enable: boolean;
    visible: boolean;
}

export const features: Feature[] = [
    { path: '/tutorial', enable: false, visible: false },
    { path: '/data-fabric', enable: false, visible: true },
    { path: '/explore', enable: false, visible: true },
    { path: '/dashboard', enable: true, visible: true },
    { path: '/applications', enable: false, visible: true },
    { path: '/integrations', enable: false, visible: true },
    { path: '/dev-tools', enable: false, visible: false }
];

export function featureTitle(path: string): string {
    switch (path) {
        case '/recents':
            return m.avatar_recents();
        case '/favorites':
            return m.avatar_favorites();
        case '/tutorial':
            return m.nav_tutorial();
        case '/data-fabric':
            return m.nav_data_fabric();
        case '/explore':
            return m.nav_explore();
        case '/dashboard':
            return m.nav_dashboard();
        case '/applications':
            return m.nav_applications();
        case '/integrations':
            return m.nav_integrations();
        case '/dev-tools':
            return m.nav_dev_tools();
        default:
            return '';
    }
}