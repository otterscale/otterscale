<script lang="ts">
	import { login, logout, register } from '$lib/auth';
	import { Button } from '$lib/components/ui/button';
	import { Spinner } from '$lib/components/ui/spinner';
	import { m } from '$lib/paraglide/messages';
</script>

<Button onclick={login}>Login with OIDC</Button>
<Button onclick={logout}>Logout</Button>
<Button onclick={register}>Register</Button>

<div class="h-svh">
	<div class="m-auto flex h-full w-1/4 flex-col items-center justify-center gap-2">
		<h1 class="text-9xl leading-tight font-bold">302</h1>
		<p class="text-muted-foreground line-clamp-2 text-center">
			{m.redirect_message_302()}
		</p>
		<div class="mt-6 flex gap-4">
			<Spinner class="size-8" />
		</div>
	</div>
</div>
