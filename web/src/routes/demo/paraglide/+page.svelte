<script lang="ts">
	import type { AvailableLanguageTag } from '$lib/paraglide/runtime';
	import { i18n } from '$lib/i18n';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import * as m from '$lib/paraglide/messages.js';

	function switchToLanguage(newLanguage: AvailableLanguageTag) {
		const canonicalPath = i18n.route(page.url.pathname);
		const localisedPath = i18n.resolveRoute(canonicalPath, newLanguage);
		goto(localisedPath);
	}
</script>

<h1>{m.hello_world({ name: 'SvelteKit User' })}</h1>
<div>
	<button onclick={() => switchToLanguage('en')}>en</button>
	<button onclick={() => switchToLanguage('zh-hant')}>zh-hant</button>
	<button onclick={() => switchToLanguage('zh-hans')}>zh-hans</button>
	<button onclick={() => switchToLanguage('de')}>de</button>
	<button onclick={() => switchToLanguage('es')}>es</button>
	<button onclick={() => switchToLanguage('fr')}>fr</button>
	<button onclick={() => switchToLanguage('it')}>it</button>
	<button onclick={() => switchToLanguage('pt')}>pt</button>
	<button onclick={() => switchToLanguage('sk')}>sk</button>
	<button onclick={() => switchToLanguage('de')}>de</button>
	<button onclick={() => switchToLanguage('jp')}>jp</button>
</div>
