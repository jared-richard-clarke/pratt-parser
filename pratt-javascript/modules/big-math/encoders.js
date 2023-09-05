import constants from "./constants.js";
import utils from "./utils.js";

// Adjust the coefficient and exponent so that the absolute value
// of the big-float object is at least 1 but less than 10.
function normalize(x) {
    let { coefficient, exponent } = x;
    if (exponent > 0) {
        coefficient = coefficient * constants.BIGINT_TEN ** BigInt(exponent);
        exponent = 0;
        return constants.make_bigfloat(coefficient, exponent);
    }
    if (exponent < 0) {
        let quotient;
        let remainder;
        while (exponent <= -7) {
            quotient = coefficient / constants.BIGINT_TEN_MILLION;
            remainder = coefficient % constants.BIGINT_TEN_MILLION;
            if (remainder !== constants.BIGINT_ZERO) {
                break;
            }
            coefficient = quotient;
            exponent += 7;
        }
        while (exponent < 0) {
            quotient = coefficient / constants.BIGINT_TEN;
            remainder = coefficient % constants.BIGINT_TEN;
            if (remainder !== constants.BIGINT_ZERO) {
                break;
            }
            coefficient = quotient;
            exponent += 1;
        }
        return constants.make_bigfloat(coefficient, exponent);
    }
    return constants.make_bigfloat(coefficient, exponent);
}

// Transform string representation of a decimal number into a big-float object.
function decode(x, y) {
    const base = y || 0;
    const pattern = /^(-?\d+)(?:\.(\d*))?(?:[eE]([+-]?\d+))?$/;
    //                ^-----^^----------^^-----------------^
    //                | int || fraction ||    exponent     |
    if (typeof x === "bigint") {
        return constants.make_bigfloat(x, base);
    }
    const match = x.match(pattern);
    if (match) {
        const integer = match[1];
        const fraction = match[2] || "";
        const exponent = match[3];
        return decode(
            BigInt(integer + fraction),
            (Number(exponent) || base) - fraction.length,
        );
    }
    return constants.BIGFLOAT_ZERO;
}

// Transform a big-float object into its decimal string representation.
function encode(x) {
    if (x.coefficient === constants.BIGINT_ZERO) {
        return "0";
    }
    x = normalize(x);
    let text = String(
        x.coefficient < constants.BIGINT_ZERO ? -x.coefficient : x.coefficient,
    );
    if (x.exponent < 0) {
        let point = text.length + x.exponent;
        if (point <= 0) {
            text = "0".repeat(1 - point) + text;
            point = 1;
        }
        text = text.slice(0, point) + "." + text.slice(point);
    }
    if (x.exponent > 0) {
        text += "0".repeat(x.exponent);
    }
    if (x.coefficient < constants.BIGINT_ZERO) {
        text = "-" + text;
    }
    return text;
}

// Transform a big-float object into its scientific notation string representation.
function encode_scientific(x) {
    if (utils.is_zero(x)) {
        return "0";
    }
    x = normalize(x);
    let text = String(
        utils.is_negative(x) ? -x.coefficient : x.coefficient,
    );
    const e = x.exponent + text.length - 1;
    if (text.length > 1) {
        text = text.slice(0, 1) + "." + text.slice(1);
        text = utils.trim_zeroes(text);
    }
    if (e !== 0) {
        text += "e" + e;
    }
    if (utils.is_negative(x)) {
        text = "-" + text;
    }
    return text;
}

export default Object.freeze({
    normalize,
    decode,
    encode,
    encode_scientific,
});
