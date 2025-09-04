<script lang="ts">
	import { LanguageSwitcher, LightSwitch, LoginForm } from '$lib/components/login';
	import { m } from '$lib/paraglide/messages.js';
	import { getLocale, setLocale, type Locale } from '$lib/paraglide/runtime.js';
	import { staticPaths } from '$lib/path';

	const { data } = $props();

	const languages = [
		{ code: 'en', label: 'English' },
		{ code: 'zh-hant', label: '繁體中文' },
	];

	let locale = $state(getLocale());

	const handleLanguageChange = (newLocale: Locale) => {
		setLocale(newLocale);
		locale = newLocale;
	};
</script>

<svelte:head>
	<title>{m.welcome_to({ name: 'OtterScale 🦦' })}</title>
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

<div class="bg-muted flex min-h-svh flex-col items-center justify-center p-6 md:p-10">
	<div class="w-full max-w-sm md:max-w-3xl">
		<div class="flex flex-col gap-6">
			<LoginForm {data} />

			<!-- Terms and Privacy -->
			<div
				class="text-muted-foreground *:[a]:hover:text-primary text-center text-xs text-balance *:[a]:underline *:[a]:underline-offset-4"
			>
				<!-- eslint-disable-next-line svelte/no-at-html-tags -->
				{@html m.login_footer({
					terms_of_service: `<a href="${staticPaths.termsOfService.url}">${m.terms_of_service()}</a>`,
					privacy_policy: `<a href="${staticPaths.privacyPolicy.url}">${m.privacy_policy()}</a>`,
				})}
			</div>
		</div>
	</div>
</div>
