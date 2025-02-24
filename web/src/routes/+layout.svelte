<script lang="ts">
	import { ParaglideJS } from '@inlang/paraglide-sveltekit';
	import { ModeWatcher } from 'mode-watcher';

	import { i18n } from '$lib/i18n';
	import { Toaster } from '$lib/components/ui/sonner';

	import { Metadata } from '$lib/components';

	import '../app.css';

	let { children } = $props();
</script>

<Metadata />

<ModeWatcher />
<Toaster />

<ParaglideJS {i18n}>
	{@render children()}
</ParaglideJS>
