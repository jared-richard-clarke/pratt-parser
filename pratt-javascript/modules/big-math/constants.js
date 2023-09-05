// To preserve precision for very large and very small numbers,
// number strings are converted into big-float objects.
//
// Big Float = { coefficient, exponent }
//    where coefficient = BigInt
//          exponent    = integer
//
// Coefficient tracks significant digits. Exponent tracks the decimal point.
// Big floats are null-inheriting, immutable objects.

// Commonly-used BigInt coefficients.
const BIGINT_NEG_ONE = -1n;
const BIGINT_ZERO = 0n;
const BIGINT_ONE = 1n;
const BIGINT_TWO = 2n;
const BIGINT_TEN = 10n;
const BIGINT_TEN_MILLION = 10000000n;

// make_bigfloat(BigInt, integer) -> { coefficient, exponent }
//
// "make_bigfloat" is not a constant, but its inclusion in this module
// eliminates circular dependencies, that is, preventing "make_bigfloat"
// from depending on constants that depend on it.
function make_bigfloat(coefficient, exponent) {
    const x = Object.create(null);
    x.coefficient = coefficient;
    x.exponent = exponent;
    return Object.freeze(x);
}

// big-float errors
const DIVIDE_ZERO = "Divide by zero. ";
const NON_INTEGER_EXPONENT = "Non-integer exponent. ";

// Commonly-used big-float objects.
const BIGFLOAT_ZERO = make_bigfloat(BIGINT_ZERO, 0);
const BIGFLOAT_ONE = make_bigfloat(BIGINT_ONE, 0);
const BIGFLOAT_TWO = make_bigfloat(BIGINT_TWO, 0);

// Desired decimal precision for big-float division.
const PRECISION = -4;

export default Object.freeze({
    make_bigfloat,
    BIGINT_NEG_ONE,
    BIGINT_ZERO,
    BIGINT_ONE,
    BIGINT_TWO,
    BIGINT_TEN,
    BIGINT_TEN_MILLION,
    DIVIDE_ZERO,
    NON_INTEGER_EXPONENT,
    BIGFLOAT_ZERO,
    BIGFLOAT_ONE,
    BIGFLOAT_TWO,
    PRECISION,
});
