import { m } from "$lib/paraglide/messages";

export const siteConfig = {
    title: m.site_title(),
    description: m.site_description()
};

export type SiteConfig = typeof siteConfig;