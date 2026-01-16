<script lang="ts">
	import { BasicForm, createForm, getValueSnapshot } from '@sjsf/form';
	import { setThemeContext } from '@sjsf/shadcn4-theme';
	import * as components from '@sjsf/shadcn4-theme/new-york';

	import { type CreateUser, schema, uiSchema, withFile } from './data';
	import * as defaults from './defaults';

	let initialValue = $state({
		name: 'Sarah Johnson',
		email: 'invalid@email',
		age: 28,
		country: 'CA',
		skills: ['HTML', 'CSS', 'JS/TS', 'Svelte'],
		experience: 'intermediate',
		startDate: new Date().toLocaleDateString('en-CA'),
		bio: 'Bio'
	}) as CreateUser;

	const form = createForm<CreateUser>({
		...defaults,
		// required due to several forms on the page
		idPrefix: 'shadcn4',
		initialValue,
		schema,
		uiSchema,
		onSubmit: ({ name }) => window.alert(`Hello, ${name}`)
	});

	//@ts-expect-error shadcn-svelte-extras
	setThemeContext({ components });

	function addAge() {
		initialValue.age = (initialValue.age ?? 0) + 2;
		initialValue['aaa'] = 'add';
	}
</script>

<BasicForm {form} />
<button type="button" onclick={addAge}>Add Age (+2)</button>

<pre>{JSON.stringify(getValueSnapshot(form), withFile, 2)}</pre>
