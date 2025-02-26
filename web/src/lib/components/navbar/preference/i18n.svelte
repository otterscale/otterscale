<script lang="ts">
	import Icon from '@iconify/svelte';

	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { goto } from '$app/navigation';
	import { i18n } from '$lib/i18n';
	import type { AvailableLanguageTag } from '$lib/paraglide/runtime';
	import { page } from '$app/state';

	let language = i18n.getLanguageFromUrl(page.url);
	const languages = new Map([
		['de' as AvailableLanguageTag, 'Deutsch'],
		['en' as AvailableLanguageTag, 'English'],
		['es' as AvailableLanguageTag, 'Español'],
		['fr' as AvailableLanguageTag, 'Français'],
		['it' as AvailableLanguageTag, 'Italiano'],
		['jp' as AvailableLanguageTag, '日本語'],
		['pt' as AvailableLanguageTag, 'Português'],
		['zh-hans' as AvailableLanguageTag, '简体中文'],
		['zh-hant' as AvailableLanguageTag, '繁體中文']
	]);

	function switchToLanguage(newLanguage: AvailableLanguageTag) {
		const canonicalPath = i18n.route(page.url.pathname);
		const localisedPath = i18n.resolveRoute(canonicalPath, newLanguage);
		goto(localisedPath);
	}
</script>

<DropdownMenu.Group class="[&_svg]:size-4">
	<DropdownMenu.Sub>
		<DropdownMenu.SubTrigger>
			<Icon icon="ph:translate" />
			<span class="pl-2">{languages.get(language)}</span>
		</DropdownMenu.SubTrigger>
		<DropdownMenu.SubContent>
			<DropdownMenu.RadioGroup bind:value={language}>
				{#each languages as language}
					<DropdownMenu.RadioItem value={language[0]} onclick={() => switchToLanguage(language[0])}>
						{language[1]}
					</DropdownMenu.RadioItem>
				{/each}
			</DropdownMenu.RadioGroup>
		</DropdownMenu.SubContent>
	</DropdownMenu.Sub>
</DropdownMenu.Group>
