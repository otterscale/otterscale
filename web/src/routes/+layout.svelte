<script lang="ts">
	import { ParaglideJS } from '@inlang/paraglide-sveltekit';
	import { ModeWatcher } from 'mode-watcher';
	import { i18n } from '$lib/i18n';
	import { Toaster } from '$lib/components/ui/sonner';
	import { Metadata } from '$lib/components';
	import { setContext } from 'svelte';
	import { createConnectTransport } from '@connectrpc/connect-web';
	import { addCollection } from '@iconify/svelte';
	import phIcons from '@iconify-json/ph/icons.json';
	import logosIcons from '@iconify-json/logos/icons.json';
	import { getIsEnterprise, setIsEnterprise } from '$lib/components/custom/enterprise';

	import '../app.css';
	import 'inter-ui/inter-variable.css';
	import '@fontsource-variable/noto-sans-tc';

	addCollection(phIcons);
	addCollection(logosIcons);

	let { children } = $props();

	const transport = createConnectTransport({
		baseUrl: import.meta.env.PUBLIC_API_URL
	});

	setContext('transport', transport);

	setContext('getIsEnterprise', getIsEnterprise);
	setContext('setIsEnterprise', setIsEnterprise);
</script>

<Metadata />

<ModeWatcher />
<Toaster closeButton richColors />

<ParaglideJS {i18n}>
	{@render children()}
</ParaglideJS>
