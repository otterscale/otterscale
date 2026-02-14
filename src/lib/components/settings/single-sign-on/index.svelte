<script lang="ts" module>
	import * as Layout from '$lib/components/settings/layout';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { m } from '$lib/paraglide/messages';

	const ssoFormFields = [
		{
			key: 'issuer',
			label: m.issuer(),
			placeholder: 'https://idp.example.com',
			type: 'text',
			span: 1
		},
		{
			key: 'domain',
			label: m.domain(),
			placeholder: 'example.com',
			type: 'text',
			span: 1
		},
		{
			key: 'clientId',
			label: m.client_id(),
			placeholder: 'client-id',
			type: 'text',
			span: 1
		},
		{
			key: 'clientSecret',
			label: m.client_secret(),
			placeholder: 'client-secret',
			type: 'password',
			span: 1
		},
		{
			key: 'authorizationEndpoint',
			label: m.authorization_endpoint(),
			placeholder: 'https://idp.example.com/authorize',
			type: 'text',
			span: 2
		},
		{
			key: 'tokenEndpoint',
			label: m.token_endpoint(),
			placeholder: 'https://idp.example.com/token',
			type: 'text',
			span: 2
		},
		{
			key: 'jwksEndpoint',
			label: m.jwks_endpoint(),
			placeholder: 'https://idp.example.com/jwks',
			type: 'text',
			span: 2
		},
		{
			key: 'discoveryEndpoint',
			label: m.discovery_endpoint(),
			placeholder: 'https://idp.example.com/.well-known/openid-configuration',
			type: 'text',
			span: 2
		}
	];
</script>

<script lang="ts">
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

	async function handleSSOSubmit(event: Event) {
		event.preventDefault();

		// TODO: implement logic to save SSO settings
	}
</script>

<Layout.Root>
	<Layout.Title>{m.single_sign_on()}</Layout.Title>
	<Layout.Description>
		{m.setting_single_sign_on_description()}
	</Layout.Description>
	<Layout.Viewer>
		<div class="w-full">
			<Card.Root>
				<Card.Header>
					<Card.Title>{m.single_sign_on_form_title()}</Card.Title>
					<Card.Description>{m.single_sign_on_form_description()}</Card.Description>
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
							<Button type="submit" class="w-full">{m.submit()}</Button>
						</Card.Footer>
					</form>
				</Card.Content>
			</Card.Root>
		</div>
	</Layout.Viewer>
</Layout.Root>
