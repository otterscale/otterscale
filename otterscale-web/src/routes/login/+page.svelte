<script lang="ts">
	import { LoginForm } from '$lib/components/login';
	import { m } from '$lib/paraglide/messages.js';
	import { privacyPolicyPath, termsOfServicePath } from '$lib/path';

	const { data } = $props();
</script>

<svelte:head>
	<title>{m.welcome_to({ name: 'OtterScale 🦦' })}</title>
</svelte:head>

<div class="bg-muted flex min-h-svh flex-col items-center justify-center p-6 md:p-10">
	<div class="w-full max-w-sm md:max-w-3xl">
		<div class="flex flex-col gap-6">
			<LoginForm {data} />

			<!-- Terms and Privacy -->
			<div
				class="text-muted-foreground *:[a]:hover:text-primary text-center text-xs text-balance *:[a]:underline *:[a]:underline-offset-4"
			>
				{@html m.login_footer({
					terms_of_service: `<a href="${termsOfServicePath}">${m.terms_of_service()}</a>`,
					privacy_policy: `<a href="${privacyPolicyPath}">${m.privacy_policy()}</a>`
				})}
			</div>
		</div>
	</div>
</div>
