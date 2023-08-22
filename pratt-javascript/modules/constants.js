function pad(x) {
    return x + " ";
}

export default Object.freeze({
    // spaces
    WHITE_SPACE: " ",
    TAB: "\t",
    LINEFEED: "\n",
    CARRIAGE_RETURN: "\r",
    VERTICAL_TAB: "\v",
    FORM_FEED: "\f",
    // general symbols
    DECIMAL_POINT: ".",
    OPEN_PAREN: "(",
    CLOSE_PAREN: ")",
    UPPER_E: "E",
    LOWER_E: "e",
    NAN: "NaN",
    // end-of-input flag
    EOF: "eof",
    // digits
    ZERO: "0",
    ONE: "1",
    TWO: "2",
    THREE: "3",
    FOUR: "4",
    FIVE: "5",
    SIX: "6",
    SEVEN: "7",
    EIGHT: "8",
    NINE: "9",
    // operators
    ADD: "+",
    SUBTRACT: "-",
    MULTIPLY: "×",
    MULTIPLY_ALT: "*",
    IMPLIED_MULTIPLY: "imp-x",
    DIVIDE: "÷",
    DIVIDE_ALT: "/",
    EXPONENT: "^",
    // token labels
    NUMBER: "number",
    ERROR: "error",
    // errors
    UNKNOWN: pad("Unknown character."),
    LEADING_ZERO: pad("Leading zero."),
    DIVIDE_ZERO: pad("Divide by zero."),
    NON_INTEGER_EXPONENT: pad("Non-integer exponent."),
    MISPLACED_DECIMAL: pad("Misplaced decimal."),
    NOT_NUMBER: pad("Not a number."),
    NO_PREFIX: pad("Undefined prefix operation."),
    NO_INFIX: pad("Undefined infix operation."),
    EMPTY_EXPRESSION: pad("Empty expression."),
    INCOMPLETE_EXPRESSION: pad("Incomplete expression."),
    MISMATCHED_PAREN: pad("Mismatched parenthesis."),
    EMPTY_PARENS: pad("Empty parentheses."),
});
