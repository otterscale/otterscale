import type { OptionType, valuesGetterType, valuesSetterType } from './types';

class OptionManager {
    visibility = 1;
    options = $state([] as OptionType[]);

    valuesSetter: valuesSetterType;
    valuesGetter: valuesGetterType;

    constructor(
        options: OptionType[],
        valuesSetter: valuesSetterType,
        valuesGetter: valuesGetterType
    ) {
        this.options = options;
        this.valuesSetter = valuesSetter;
        this.valuesGetter = valuesGetter
    }

    updateOptions(newOptions: OptionType[]) {
        this.options = newOptions
    }

    get selectedOptions(): OptionType[] {
        return this.options.filter((option) => (this.valuesGetter().includes(option.value)))
    }

    get isSomeOptionsSelected(): boolean {
        return this.valuesGetter().length > 0;
    }

    isOptionSelected(option: OptionType): boolean {
        return this.valuesGetter().includes(option.value);
    }

    handleSelect(option: OptionType) {
        if (this.isOptionSelected(option)) {
            this.valuesSetter(this.valuesGetter().filter((o) => o !== option.value));
        } else {
            this.valuesSetter([...this.valuesGetter(), option.value]);
        }
    }

    all() {
        this.valuesSetter(
            this.options.map((option) => (option.value))
        );
    }

    clear() {
        this.valuesSetter([]);
    }
}

export {
    OptionManager
};
