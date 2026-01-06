<script lang="ts">
	import { LanguageSwitcher, LightSwitch, LoginForm } from '$lib/components/login';
	import { m } from '$lib/paraglide/messages.js';
	import { getLocale, type Locale, setLocale } from '$lib/paraglide/runtime.js';

	const languages = [
		{ code: 'en', label: 'English' },
		{ code: 'zh-hant', label: '繁體中文' }
	];

	let locale = $state(getLocale());

	const handleLanguageChange = (newLocale: Locale) => {
		setLocale(newLocale);
		locale = newLocale;
	};
</script>

<svelte:head>
	<title>{m.welcome_to({ name: 'OtterScale' })}</title>
</svelte:head>

<div class="absolute top-4 right-4 flex p-4 md:top-12 md:right-12">
	<LightSwitch />
	<LanguageSwitcher
		{languages}
		bind:value={locale}
		onChange={(newLocale: string) => {
			handleLanguageChange(newLocale as Locale);
		}}
	/>
</div>

<div class="flex h-screen w-full items-center justify-center bg-muted px-4">
	<LoginForm />
</div>
