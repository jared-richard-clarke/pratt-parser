import constants from "./constants.js";

function is_zero(x) {
    return x.coefficient === constants.BIGINT_ZERO;
}
function is_negative(x) {
    return x.coefficient < constants.BIGINT_ZERO;
}
function is_positive(x) {
    return x.coefficient >= constants.BIGINT_ZERO;
}
function is_odd(x) {
    return (x.coefficient % constants.BIGINT_TWO) !== constants.BIGINT_ZERO;
}
function bigint_abs(x) {
    return x < constants.BIGINT_ZERO ? -x : x;
}
function bigint_signum(x) {
    if (x === constants.BIGINT_ZERO) {
        return constants.BIGINT_ZERO;
    }
    if (x < constants.BIGINT_ZERO) {
        return constants.BIGINT_NEG_ONE;
    }
    return constants.BIGINT_ONE;
}

export default Object.freeze({
    is_zero,
    is_negative,
    is_positive,
    is_odd,
    bigint_abs,
    bigint_signum,
});
