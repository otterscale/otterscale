<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { goto } from '$app/navigation';
	import { ScopeService, type Scope } from '$lib/api/scope/v1/scope_pb';
	import { authClient } from '$lib/auth-client';
	import { Button } from '$lib/components/ui/button';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths, staticPaths } from '$lib/path';

	const transport: Transport = getContext('transport');
	const scopeClient = createClient(ScopeService, transport);
	const scopes = writable<Scope[]>([]);

	async function fetchScopes() {
		try {
			const response = await scopeClient.listScopes({});
			scopes.set(response.scopes);
		} catch (error) {
			console.error('Failed to fetch scopes:', error);
		}
	}

	onMount(fetchScopes);
</script>

<div class="flex space-x-2">
	{#each $scopes as scope}
		<Button href={dynamicPaths.scope(scope.name).url}>{scope.name}</Button>
	{/each}
</div>

<Button
	onclick={() => {
		authClient.signOut({
			fetchOptions: {
				onSuccess: () => {
					toast.success(m.sign_out_success());
					goto(staticPaths.login.url);
				}
			}
		});
	}}>LOGOUT</Button
>
