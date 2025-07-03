import type { OptionType, valueGetterType, valueSetterType } from './types';

class OptionManager {
    options = $state([] as OptionType[]);
    valueSetter: valueSetterType;
    valueGetter: valueGetterType;

    constructor(options: OptionType[], valueSetter: valueSetterType, valueGetter: valueGetterType) {
        this.options = options;
        this.valueSetter = valueSetter
        this.valueGetter = valueGetter
    }

    updateOptions(newOptions: OptionType[]) {
        this.options = newOptions
    }

    get selectedOption(): OptionType {
        return this.options.find((option) => (JSON.stringify(option.value) === JSON.stringify(this.valueGetter()))) ?? {} as OptionType
    }

    handleSelect(option: OptionType) {
        this.valueSetter(option.value);
    }
}

export {
    OptionManager
};
