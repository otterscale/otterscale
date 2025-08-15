<script lang="ts" module>
	import { authClient } from '$lib/auth-client';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { toast } from 'svelte-sonner';

	const providerId = 'otterscale-oidc';

	const formFields = [
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
</script>

<script lang="ts">
	let formData = {
		issuer: '',
		domain: '',
		clientId: '',
		clientSecret: '',
		authorizationEndpoint: '',
		tokenEndpoint: '',
		jwksEndpoint: '',
		discoveryEndpoint: ''
	};

	async function handleSubmit(event: Event) {
		event.preventDefault();

		toast.promise(
			authClient.sso.register({
				providerId,
				issuer: formData.issuer,
				domain: formData.domain,
				oidcConfig: {
					clientId: formData.clientId,
					clientSecret: formData.clientSecret,
					authorizationEndpoint: formData.authorizationEndpoint,
					tokenEndpoint: formData.tokenEndpoint,
					jwksEndpoint: formData.jwksEndpoint,
					discoveryEndpoint: formData.discoveryEndpoint,
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

<!-- <span>TODO: LAYOUT</span>
<span>TODO: TITLE</span>
<span>TODO: I18N</span>
<span>TODO: SAML</span> -->

<Card.Root>
	<Card.Header>
		<Card.Title>Configure OIDC Provider</Card.Title>
		<Card.Description>Set up your OpenID Connect provider configuration</Card.Description>
	</Card.Header>
	<Card.Content>
		<form on:submit={handleSubmit}>
			<div class="flex flex-col gap-6">
				<div class="grid grid-cols-2 gap-4">
					{#each formFields as field}
						<div class="grid gap-2">
							<Label for={field.key}>{field.label}</Label>
							<Input
								id={field.key}
								type={field.type}
								placeholder={field.placeholder}
								bind:value={formData[field.key as keyof typeof formData]}
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
