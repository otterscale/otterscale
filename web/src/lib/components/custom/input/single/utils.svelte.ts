const BORDER_INPUT_CLASSNAME = 'flex items-center rounded-md border shadow';
const UNFOCUS_INPUT_CLASSNAME = 'border-none shadow-none focus-visible:ring-0 bg-transparent';

const typeToIcon: Record<string, string> = {
    color: 'ph:palette',
    'datetime-local': 'ph:clock',
    date: 'ph:calendar',
    time: 'ph:clock',
    url: 'ph:link',
    email: 'ph:mailbox',
    tel: 'ph:phone',
    boolean: 'ph:check',
    text: 'ph:textbox',
    number: 'ph:list-numbers',
    search: 'ph:magnifying-glass',
    password: 'ph:password'
};

class PasswordManager {
    isVisible = $state<boolean>(false);

    enable() {
        this.isVisible = true;
    }

    disable() {
        this.isVisible = false;
    }
}

export {
    BORDER_INPUT_CLASSNAME,
    UNFOCUS_INPUT_CLASSNAME,
    typeToIcon,
    //
    PasswordManager,
}