<script lang="ts">
	import Icon from '@iconify/svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	import { resolve } from '$app/paths';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { FieldDescription } from '$lib/components/ui/field';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	let { class: className, ...restProps }: HTMLAttributes<HTMLDivElement> = $props();
</script>

<div class={cn('flex flex-col gap-6', className)} {...restProps}>
	<Card.Root class="overflow-hidden p-0">
		<Card.Content class="grid p-0">
			<div class="p-6 md:p-8">
				<div class="flex flex-col gap-6">
					<div class="flex flex-col items-center text-center">
						<h1 class="text-2xl font-bold">{m.login_title()}</h1>
						<p class="text-balance text-muted-foreground">{m.login_description()}</p>
					</div>

					<Button class="w-full" href={resolve('/login/keycloak')}>
						<Icon icon="simple-icons:keycloak" />
						{m.login_with({ provider: 'Keycloak' })}
					</Button>
				</div>
			</div>
		</Card.Content>
	</Card.Root>

	<FieldDescription class="px-6 text-center">
		<!-- eslint-disable-next-line svelte/no-at-html-tags -->
		{@html m.login_footer({
			terms_of_service: `<a href="${resolve('/terms-of-service')}">${m.terms_of_service()}</a>`,
			privacy_policy: `<a href="${resolve('/privacy-policy')}">${m.privacy_policy()}</a>`
		})}
	</FieldDescription>
</div>
