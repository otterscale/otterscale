<script lang="ts">
	import { UserCreateForm } from '$lib/components/login';
	import Button from '$lib/components/ui/button/button.svelte';
	import { appendCallback } from '$lib/utils';
</script>

<div
	class="container relative h-screen flex-col items-center justify-center md:grid lg:max-w-none lg:grid-cols-2 lg:px-0"
>
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
				<h1 class="text-2xl font-semibold tracking-tight">Sign Up</h1>
				<p class="text-sm text-muted-foreground">Enter your information to create an account</p>
			</div>
			<UserCreateForm />
			<div class="text-center text-sm text-muted-foreground">
				Already have an account?
				<a href={appendCallback('/login')} class="underline"> Sign in </a>
			</div>
		</div>
	</div>
</div>
