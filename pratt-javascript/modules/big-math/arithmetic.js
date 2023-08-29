import constants from "./constants.js";
import encoders from "./encoders.js";
import utils from "./utils.js";

function neg(x) {
    return constants.make_bigfloat(-x.coefficient, x.exponent);
}
function floor(x) {
    const { coefficient, exponent } = x;
    if (exponent === 0) {
        return x;
    }
    if (exponent > 0) {
        return constants.make_bigfloat(
            coefficient * constants.BIGINT_TEN ** BigInt(exponent),
            0,
        );
    }
    return constants.make_bigfloat(
        coefficient / constants.BIGINT_TEN ** BigInt(-exponent),
        0,
    );
}
function match_terms(op) {
    return function (x, y) {
        const differential = x.exponent - y.exponent;
        if (differential === 0) {
            return constants.make_bigfloat(
                op(x.coefficient, y.coefficient),
                x.exponent,
            );
        }
        if (differential > 0) {
            return constants.make_bigfloat(
                op(
                    x.coefficient *
                        constants.BIGINT_TEN ** BigInt(differential),
                    y.coefficient,
                ),
                y.exponent,
            );
        }
        return constants.make_bigfloat(
            op(
                x.coefficient,
                y.coefficient * constants.BIGINT_TEN ** BigInt(-differential),
            ),
            x.exponent,
        );
    };
}

const add = match_terms((x, y) => x + y);
const sub = match_terms((x, y) => x - y);

function eq(x, y) {
    x = encoders.normalize(x);
    y = encoders.normalize(y);
    return (x.coefficient === y.coefficient) && (x.exponent === y.exponent);
}
function ge(x, y) {
    const difference = sub(x, y);
    return difference.coefficient >= constants.BIGINT_ZERO;
}
function gt(x, y) {
    const difference = sub(x, y);
    return difference.coefficient > constants.BIGINT_ZERO;
}
function le(x, y) {
    const difference = sub(x, y);
    return difference.coefficient <= constants.BIGINT_ZERO;
}
function lt(x, y) {
    const difference = sub(x, y);
    return difference.coefficient < constants.BIGINT_ZERO;
}

function mul(x, y) {
    return constants.make_bigfloat(
        x.coefficient * y.coefficient,
        x.exponent + y.exponent,
    );
}

// Division by zero is undefined.
function div(x, y) {
    if (utils.is_zero(x)) {
        return constants.BIGFLOAT_ZERO;
    }
    if (utils.is_zero(y)) {
        return constants.DIVIDE_ZERO;
    }
    const precision = constants.PRECISION;
    let { coefficient, exponent } = x;
    exponent -= y.exponent;
    if (exponent > precision) {
        coefficient = coefficient *
            constants.BIGINT_TEN ** BigInt(exponent - precision);
        exponent = precision;
    }
    coefficient = coefficient / y.coefficient;
    const remainder = coefficient % y.coefficient;
    // round
    if (
        utils.bigint_abs(remainder + remainder) >=
            utils.bigint_abs(y.coefficient)
    ) {
        coefficient += utils.bigint_signum(x.coefficient);
    }
    return constants.make_bigfloat(coefficient, exponent);
}

// Exponentiation by squaring.
// - (-2) ^ 0.5 returns NaN. A negative base and a non-integer exponent
// cannot be defined as a real number.
// -2 ^ (-2 ^ -2) -> NaN | (-2 ^ -2) ^ -2 -> 16
//  ^ negative base and non-integer exponent.
// Some mathematicians consider 0 ^ 0 an undefined value.
// 0 ** -y -> Infinity, 0n ** -y -> Range Error
// pow(a,b) = exp( b * log(a) )

// pow(a,b) = exp( b * log(a) )
// exp(x)   = sum[n = 0->inf]  x^n / n!    // <- Taylor Series
// ln(x)    = sum[n = 1->inf]  (x-1)^n / n
// x^n      = n == 0  ? 1   // unless x == 0
//           (n%2==0) ? x^(n/2) * x^(n/2)
//           otherwise x*x^(n-1)
//           // faster than loop for large n,

// A typical implementation of pow(x,y) (i.e. x^y) involves computing exp(y*log(x))

// A fractional exponent like 1/n means to take the n-th root:

// x1/n  =  n√x

function pow(x, y) {
    // Current implementation does not support non-integer exponents.
    if (encoders.normalize(y).exponent !== 0) {
        return constants.NON_INTEGER_EXPONENT;
    }
    if (utils.is_zero(x)) {
        return constants.BIGFLOAT_ZERO;
    }
    if (utils.is_zero(y)) {
        return constants.BIGFLOAT_ONE;
    }
    if (utils.is_negative(y)) {
        x = div(constants.BIGFLOAT_ONE, x);
        y = neg(y);
    }
    let z = constants.BIGFLOAT_ONE;
    while (gt(y, constants.BIGFLOAT_ONE)) {
        if (utils.is_odd(y)) {
            z = mul(x, z);
        }
        x = mul(x, x);
        y = floor(div(y, constants.BIGFLOAT_TWO));
    }
    return mul(x, z);
}

export default Object.freeze({
    eq,
    ge,
    gt,
    le,
    lt,
    neg,
    add,
    sub,
    mul,
    div,
    pow,
});
