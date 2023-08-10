const BIGINT_ZERO = 0n;
const BIGINT_ONE = 1n;
const BIGINT_TEN = 10n;
const BIGINT_TEN_MILLION = 10000000n;

const ZERO = (function () {
    const x = Object.create(null);
    x.coefficient = BIGINT_ZERO;
    x.exponent = 0;
    return Object.freeze(x);
})();

const PRECISION = -4;

export default Object.freeze({
    BIGINT_ZERO,
    BIGINT_ONE,
    BIGINT_TEN,
    BIGINT_TEN_MILLION,
    ZERO,
    PRECISION,
});
