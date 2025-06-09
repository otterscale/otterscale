import type { OptionType, valueSetterType, valueGetterType } from './types';

class OptionManager {
    // selectedOption: OptionType = $state({} as OptionType);
    options: OptionType[];
    valueSetter: valueSetterType;
    valueGetter: valueGetterType;

    constructor(options: OptionType[], valueSetter: valueSetterType, valueGetter: valueGetterType) {
        // this.selectedOption = selectedOption ?? {} as OptionType;
        this.options = options;
        this.valueSetter = valueSetter
        this.valueGetter = valueGetter
    }

    get selectedOption(): OptionType {
        return this.options.find((option) => (option.value === this.valueGetter())) ?? {} as OptionType
    }

    handleSelect(option: OptionType) {
        // this.selectedOption = option;
        this.valueSetter(option.value);
    }
}

export {
    OptionManager
}