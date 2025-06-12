import { typeToIcon } from "../single";
import type { InputType, valueSetterType } from './types';

class InputManager {
    input: any = $state();
    type: InputType = 'text';

    constructor(type: InputType) {
        this.type = type
    }

    reset() {
        this.input = '';
    }
}

class ValuesManager {
    values: any[] = $state([] as any[]);
    valuesSetter: valueSetterType;

    constructor(initialValues: any, valuesSetter: valueSetterType) {
        this.values = Array.isArray(initialValues) ? initialValues : initialValues ? [initialValues] : []
        this.valuesSetter = valuesSetter
    }

    append(value: any) {
        if (value === undefined || value === null || String(value).trim() === '') {
            return;
        }
        if (this.values.includes(value)) return;
        this.values = [...this.values, value];
        this.valuesSetter(this.values)
    }

    remove(value: any) {
        this.values = this.values.filter((v) => v !== value);
        this.valuesSetter(this.values)
    }

    reset() {
        this.values = [];
        this.valuesSetter(this.values)
    }
}

export {
    InputManager, typeToIcon, ValuesManager
};
