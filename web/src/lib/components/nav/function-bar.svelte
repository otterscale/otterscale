<script lang="ts">
	import Icon from '@iconify/svelte';
	import { mode, toggleMode } from 'mode-watcher';
	import { toast } from 'svelte-sonner';

	import { page } from '$app/state';
	import { goto } from '$app/navigation';

	import type { AvailableLanguageTag } from '$lib/paraglide/runtime';
	import { Button } from '$lib/components/ui/button';
	import { i18n } from '$lib/i18n';

	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import pb from '$lib/pb';

	let isValid = pb.authStore.isValid;
	pb.authStore.onChange(() => {
		isValid = pb.authStore.isValid;
	});

	function switchToLanguage(newLanguage: AvailableLanguageTag) {
		const canonicalPath = i18n.route(page.url.pathname);
		const localisedPath = i18n.resolveRoute(canonicalPath, newLanguage);
		goto(localisedPath);
	}

	let favorited = false;
	function toggleFavorite() {
		if (!favorited) {
			favorited = true;
			toast.success('Added to favorites!');
			return;
		}
		favorited = false;
		toast.error('Removed from favorites!');
	}

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
</script>

<div class="flex justify-end space-x-2">
	<DropdownMenu.Root>
		<DropdownMenu.Trigger asChild let:builder>
			<Button builders={[builder]} class="w-14" variant="outline" size="icon">
				<Icon icon="material-symbols:add-2-rounded" class="h-5 w-5" />
				<Icon icon="material-symbols:arrow-drop-down-rounded" class="h-5 w-5" />
			</Button>
		</DropdownMenu.Trigger>
		<DropdownMenu.Content align="end">
			<DropdownMenu.Item class="space-x-2" on:click={toggleMode}>
				<span>{$mode === 'light' ? 'Light Mode' : 'Dark Mode'}</span>
			</DropdownMenu.Item>
			<DropdownMenu.Group>
				<DropdownMenu.Sub>
					<DropdownMenu.SubTrigger>{languages.get(language)}</DropdownMenu.SubTrigger>
					<DropdownMenu.SubContent>
						<DropdownMenu.RadioGroup bind:value={language}>
							{#each languages as language}
								<DropdownMenu.RadioItem
									value={language[0]}
									on:click={() => switchToLanguage(language[0])}
								>
									{language[1]}
								</DropdownMenu.RadioItem>
							{/each}
						</DropdownMenu.RadioGroup>
					</DropdownMenu.SubContent>
				</DropdownMenu.Sub>
			</DropdownMenu.Group>
			<DropdownMenu.Item
				class="space-x-2"
				on:click={() => window.open('https://openhdc.github.io', '_blank')}
			>
				<Icon icon="material-symbols:open-in-new" class="h-5 w-5" />
				<span>Documentation</span>
			</DropdownMenu.Item>
		</DropdownMenu.Content>
	</DropdownMenu.Root>
	{#if isValid}
		<Button variant="outline" size="icon" class="bg-header" on:click={toggleFavorite}>
			{#if favorited}
				<Icon icon="material-symbols:favorite" class="h-5 w-5" />
			{:else}
				<Icon icon="material-symbols:favorite-outline" class="h-5 w-5" />
			{/if}
		</Button>
		<Button variant="outline" size="icon" class="bg-header">
			<Icon icon="material-symbols:inbox-outline-rounded" class="h-5 w-5" />
		</Button>
	{:else}
		<Button
			variant="outline"
			size="icon"
			class="h-9 w-20 border-none"
			on:click={() => goto('/account/login')}
		>
			<span>Sign in</span>
		</Button>
		<Button variant="outline" size="icon" class="h-9 w-20" on:click={() => goto('/account/signup')}>
			<span>Sign up</span>
		</Button>
	{/if}
</div>
