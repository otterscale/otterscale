<script lang="ts">
	import { page } from '$app/state';
	import { authClient } from '$lib/auth-client';
	import * as Layout from '$lib/components/settings/general/layout';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';
	import { type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	const ssoProviderId = 'otterscale-oidc';

	const ssoFormFields = [
		{
			key: 'issuer',
			label: 'Issuer',
			placeholder: 'https://idp.example.com',
			type: 'text',
			span: 1
		},
		{
			key: 'domain',
			label: 'Domain',
			placeholder: 'example.com',
			type: 'text',
			span: 1
		},
		{
			key: 'clientId',
			label: 'Client ID',
			placeholder: 'client-id',
			type: 'text',
			span: 1
		},
		{
			key: 'clientSecret',
			label: 'Client Secret',
			placeholder: 'client-secret',
			type: 'password',
			span: 1
		},
		{
			key: 'authorizationEndpoint',
			label: 'Authorization Endpoint',
			placeholder: 'https://idp.example.com/authorize',
			type: 'text',
			span: 2
		},
		{
			key: 'tokenEndpoint',
			label: 'Token Endpoint',
			placeholder: 'https://idp.example.com/token',
			type: 'text',
			span: 2
		},
		{
			key: 'jwksEndpoint',
			label: 'JWKS Endpoint',
			placeholder: 'https://idp.example.com/jwks',
			type: 'text',
			span: 2
		},
		{
			key: 'discoveryEndpoint',
			label: 'Discovery Endpoint',
			placeholder: 'https://idp.example.com/.well-known/openid-configuration',
			type: 'text',
			span: 2
		}
	];
	let ssoFormData = {
		issuer: '',
		domain: '',
		clientId: '',
		clientSecret: '',
		authorizationEndpoint: '',
		tokenEndpoint: '',
		jwksEndpoint: '',
		discoveryEndpoint: ''
	};
	let isMounted = $state(false);

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.settings(page.params.scope)],
		current: { title: m.sso(), url: '' }
	});

	const transport: Transport = getContext('transport');

	async function handleSSOSubmit(event: Event) {
		event.preventDefault();

		toast.promise(
			authClient.sso.register({
				providerId: ssoProviderId,
				issuer: ssoFormData.issuer,
				domain: ssoFormData.domain,
				oidcConfig: {
					clientId: ssoFormData.clientId,
					clientSecret: ssoFormData.clientSecret,
					authorizationEndpoint: ssoFormData.authorizationEndpoint,
					tokenEndpoint: ssoFormData.tokenEndpoint,
					jwksEndpoint: ssoFormData.jwksEndpoint,
					discoveryEndpoint: ssoFormData.discoveryEndpoint,
					scopes: ['openid', 'email', 'profile'],
					pkce: true
				},
				mapping: {
					id: 'sub',
					email: 'email',
					emailVerified: 'email_verified',
					name: 'name',
					image: 'picture'
				}
			}),
			{
				loading: 'Loading...',
				success: 'OIDC Provider has been updated!',
				error: 'An error occurred'
			}
		);
	}
</script>

<Layout.Title>Single Sign On</Layout.Title>
<Layout.Description>
	Single Sign-On (SSO) allows users to access multiple applications with one set of credentials.
	Configure your SSO provider details here to enable centralized authentication across your
	infrastructure management system.
</Layout.Description>
<Layout.Controller>
	<Card.Root>
		<Card.Header>
			<Card.Title>Configure OIDC Provider</Card.Title>
			<Card.Description>Set up your OpenID Connect provider configuration</Card.Description>
		</Card.Header>
		<Card.Content>
			<form onsubmit={handleSSOSubmit}>
				<div class="flex flex-col gap-6">
					<div class="grid grid-cols-2 gap-4">
						{#each ssoFormFields as field}
							<div class="grid gap-2">
								<Label for={field.key}>{field.label}</Label>
								<Input
									id={field.key}
									type={field.type}
									placeholder={field.placeholder}
									bind:value={ssoFormData[field.key as keyof typeof ssoFormData]}
									required
								/>
							</div>
						{/each}
					</div>
				</div>
				<Card.Footer class="flex-col gap-2 px-0 pt-6">
					<Button type="submit" class="w-full">Submit</Button>
				</Card.Footer>
			</form>
		</Card.Content>
	</Card.Root>
</Layout.Controller>
