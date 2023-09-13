export default Object.freeze({
    // === spaces ===
    WHITE_SPACE: " ",
    TAB: "\t",
    LINEFEED: "\n",
    CARRIAGE_RETURN: "\r",
    VERTICAL_TAB: "\v",
    FORM_FEED: "\f",
    // === general symbols ===
    DECIMAL_POINT: ".",
    OPEN_PAREN: "(",
    CLOSE_PAREN: ")",
    UPPER_E: "E",
    LOWER_E: "e",
    NAN: "NaN",
    UNDEFINED: "undefined",
    INFINITY: "Infinity",
    // === end-of-input flag ===
    EOF: "eof",
    // === digits ===
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
    // === operators ===
    ADD: "+",
    // U+2212: unicode minus
    SUBTRACT: "−",
    // U+002d: hyphen-minus
    SUBTRACT_ALT: "-",
    // U+00d7: multiplication
    MULTIPLY: "×",
    MULTIPLY_ALT: "*",
    IMPLIED_MULTIPLY: "imp-x",
    DIVIDE: "÷",
    DIVIDE_ALT: "/",
    EXPONENT: "^",
    // === token labels ===
    NUMBER: "number",
    ERROR: "error",
    // === errors ===
    // Right padding added for formatting.
    UNKNOWN: "Unknown identifier. ",
    LEADING_ZERO: "Leading zero. ",
    MISPLACED_NUMBER: "Misplaced number. ",
    MISPLACED_DECIMAL: "Misplaced decimal. ",
    MISPLACED_EXPONENT: "Misplaced exponent suffix. ",
    NOT_NUMBER: "Not a number. ",
    NOT_PREFIX: "Not prefix operation. ",
    NOT_INFIX: "Not infix operation. ",
    EMPTY_EXPRESSION: "Empty expression. ",
    INCOMPLETE_EXPRESSION: "Incomplete expression. ",
    MISMATCHED_PAREN: "Mismatched parenthesis. ",
    EMPTY_PARENS: "Empty parentheses. ",
});
