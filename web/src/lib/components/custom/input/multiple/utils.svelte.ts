import { typeToIcon } from '../single';
import type { AccessorType, InputType } from './types';

class InputManager {
	input: any = $state();
	type: InputType | undefined;
	icon: string | undefined;

	constructor(type: InputType | undefined, icon: string | undefined) {
		this.type = type ?? 'text';
		this.icon = icon ?? typeToIcon[this.type];
	}

	reset() {
		this.input = undefined;
	}
}

class ValuesManager {
	values: any[] = $state([] as any[]);
	accessor: AccessorType;

	constructor(initialValues: any, accessor: AccessorType) {
		this.values = Array.isArray(initialValues)
			? initialValues
			: initialValues
				? [initialValues]
				: [];
		this.accessor = accessor;
	}

	append(value: any) {
		if (value === undefined || value === null || String(value).trim() === '') {
			return;
		}
		if (this.values.includes(value)) return;
		this.values = [...this.values, value];
		this.accessor.values = this.values;
	}

	remove(value: any) {
		this.values = this.values.filter((v) => v !== value);
		this.accessor.values = this.values;
	}

	reset() {
		this.values = [];
		this.accessor.values = this.values;
	}
}

function validate(required: boolean | undefined, valuesManager: ValuesManager) {
	return required && valuesManager.values.length === 0;
}

export { InputManager, typeToIcon, validate, ValuesManager };
