import { typeToIcon } from "../single";
import type { InputType, valueSetterType } from './types'

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

    constructor(values: any[], valuesSetter: valueSetterType) {
        this.values = values
        this.valuesSetter = valuesSetter
    }

    append(value: any) {
        console.log(typeof value);
        if (value === undefined || value === null || (typeof value === 'string' && value.trim() === '')) {
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
    typeToIcon,
    //
    InputManager,
    ValuesManager,
}