<script lang="ts">
	import * as Avatar from '$lib/components/ui/avatar';
	import { Button } from '$lib/components/ui/button';
	import { SignupForm } from '$lib/components/login';
	import { appendCallback } from '$lib/callback';
</script>

<div class="flex min-h-svh flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10">
	<div class="flex w-full max-w-sm flex-col items-center gap-6">
		<Avatar.Root class="size-24 shadow-xl">
			<Avatar.Image src="/images/otterscale.jpeg" alt="@shadcn" />
			<Avatar.Fallback>OS</Avatar.Fallback>
		</Avatar.Root>
		<SignupForm />
	</div>
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

<div class="absolute right-4 top-4 flex md:right-8 md:top-8">
	<Button href={appendCallback('/login')} variant="ghost">Login</Button>
</div>
