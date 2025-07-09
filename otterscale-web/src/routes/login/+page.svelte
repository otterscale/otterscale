<script lang="ts">
	import { authClient } from '$lib/auth-client';
	import { Button } from '$lib/components/ui/button';

	const { data } = $props();

	const providers = [
		{ id: 'apple', name: 'Apple', enabled: data.apple },
		{ id: 'github', name: 'GitHub', enabled: data.github },
		{ id: 'google', name: 'Google', enabled: data.google }
	] as const;

	function handleSignIn(provider: string) {
		authClient.signIn.social({
			provider: provider as any,
			callbackURL: data.nextPath
		});
	}
</script>

<div class="flex min-h-screen items-center justify-center">
	<div class="bg-card flex flex-col items-center gap-8 rounded-xl border p-10">
		Logo here
		<div class="flex flex-col gap-2">
			{#each providers as provider}
				<Button
					type="button"
					disabled={!provider.enabled}
					onclick={() => handleSignIn(provider.id)}
					variant="outline"
					size="lg"
				>
					Sign in with {provider.name}
				</Button>
			{/each}
		</div>
	</div>
</div>
