import type { OptionType, valueSetterType, valueGetterType } from './types';

class OptionManager {
    options: OptionType[];
    valueSetter: valueSetterType;
    valueGetter: valueGetterType;

    constructor(options: OptionType[], valueSetter: valueSetterType, valueGetter: valueGetterType) {
        this.options = options;
        this.valueSetter = valueSetter
        this.valueGetter = valueGetter
    }

    get selectedOption(): OptionType {
        return this.options.find((option) => (option.value === this.valueGetter())) ?? {} as OptionType
    }

    handleSelect(option: OptionType) {
        this.valueSetter(option.value);
    }
}

export {
    OptionManager
}