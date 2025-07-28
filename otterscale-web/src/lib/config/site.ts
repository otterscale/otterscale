import { m } from "$lib/paraglide/messages";

export const siteConfig = {
    title: m.site_title(),
    description: m.site_description(),
    version: "v1.0.0-beta.3",
};

export type SiteConfig = typeof siteConfig;