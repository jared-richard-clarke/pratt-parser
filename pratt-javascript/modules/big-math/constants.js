const BIGINT_NEG_ONE = -1n;
const BIGINT_ZERO = 0n;
const BIGINT_ONE = 1n;
const BIGINT_TWO = 2n;
const BIGINT_TEN = 10n;
const BIGINT_TEN_MILLION = 10000000n;

const BIGFLOAT_ZERO = (function () {
    const x = Object.create(null);
    x.coefficient = BIGINT_ZERO;
    x.exponent = 0;
    return Object.freeze(x);
})();

const BIGFLOAT_ONE = (function () {
    const x = Object.create(null);
    x.coefficient = BIGINT_ONE;
    x.exponent = 0;
    return Object.freeze(x);
})();

const BIGFLOAT_TWO = (function () {
    const x = Object.create(null);
    x.coefficient = BIGINT_TWO;
    x.exponent = 0;
    return Object.freeze(x);
})();

const PRECISION = -4;

export default Object.freeze({
    BIGINT_NEG_ONE,
    BIGINT_ZERO,
    BIGINT_ONE,
    BIGINT_TWO,
    BIGINT_TEN,
    BIGINT_TEN_MILLION,
    BIGFLOAT_ZERO,
    BIGFLOAT_ONE,
    BIGFLOAT_TWO,
    PRECISION,
});
