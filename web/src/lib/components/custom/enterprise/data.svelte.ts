let isEnterprise = $state(false);

function getIsEnterprise() {
    return isEnterprise;
}
function setIsEnterprise(newValue: boolean) {
    isEnterprise = newValue;
}

export {
    getIsEnterprise, setIsEnterprise
}