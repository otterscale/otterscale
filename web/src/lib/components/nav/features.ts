import * as m from '$lib/paraglide/messages.js';

export const features = [
    {
        path: '/tutorial',
        name: m.nav_tutorial(),
        enable: false
    },
    {
        path: '/data-fabric',
        name: m.nav_data_fabric(),
        enable: true
    },
    {
        path: '/explore',
        name: m.nav_explore(),
        enable: false
    },
    {
        path: '/dashboard',
        name: m.nav_dashboard(),
        enable: true
    },
    {
        path: '/applications',
        name: m.nav_applications(),
        enable: true
    },
    {
        path: '/integrations',
        name: m.nav_integrations(),
        enable: false
    },
    {
        path: '/dev-tools',
        name: m.nav_dev_tools(),
        enable: false
    }
];
