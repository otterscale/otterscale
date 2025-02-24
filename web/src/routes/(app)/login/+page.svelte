<script lang="ts">
	import { onMount } from 'svelte';

	import { goto } from '$app/navigation';
	import { UserAuthForm } from '$lib/components/login';
	import { Button } from '$lib/components/ui/button';
	import { i18n } from '$lib/i18n';
	import { appendCallback } from '$lib/utils';
	import pb from '$lib/pb';

	onMount(() => {
		if (pb.authStore.isValid) {
			goto(i18n.resolveRoute('/'));
		}
	});
</script>

<div
	class="container relative h-screen flex-col items-center justify-center md:grid lg:max-w-none lg:grid-cols-2 lg:px-0"
>
	<div class="absolute right-4 top-4 flex md:right-8 md:top-8">
		<Button href={appendCallback('/signup')} variant="ghost">Sign up</Button>
	</div>
	<div class="relative hidden h-full flex-col bg-muted p-10 dark:border-r lg:flex">
		<img
			src="/images/placeholder.svg"
			alt="placeholder"
			class="absolute inset-0 h-full w-full object-cover dark:brightness-[0.6] dark:grayscale"
		/>
		<div class="relative z-20 flex items-center text-lg font-medium">
			<img src="/images/phison.svg" alt="phison" class="w-28" />
		</div>
	</div>
	<div class="p-8">
		<div class="mx-auto flex w-full flex-col items-center justify-center space-y-6">
			<div class="flex flex-col space-y-2 text-center">
				<h1 class="text-2xl font-semibold tracking-tight">Sign In</h1>
				<p class="text-sm text-muted-foreground">Enter your email to access your account</p>
			</div>
			<UserAuthForm />
			<div class="space-y-1 text-center text-xs text-muted-foreground">
				<p>By continuing, you agree to our</p>
				<p>
					<a href="/terms" class="underline underline-offset-4 hover:text-primary">
						Terms of Service
					</a>
					and
					<a href="/privacy" class="underline underline-offset-4 hover:text-primary">
						Privacy Policy
					</a>
					.
				</p>
			</div>
		</div>
	</div>
</div>
