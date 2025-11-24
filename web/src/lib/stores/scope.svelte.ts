let currentScope = $state<string>('');

export const scopeStore = {
	get value() {
		return currentScope;
	},
	set(value: string) {
		currentScope = value;
	}
};
